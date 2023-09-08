package search

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"math"
	"reflect"
	"sort"
	"sync"
	"time"

	"github.com/oarkflow/xid"

	"github.com/oarkflow/pkg/search/lib"
	"github.com/oarkflow/pkg/search/tokenizer"
	"github.com/oarkflow/pkg/str"
	"github.com/oarkflow/pkg/utils"
)

const (
	AND Mode = "AND"
	OR  Mode = "OR"
)

const WILDCARD = "*"

type Mode string

type SchemaProps any

type Record[Schema SchemaProps] struct {
	Id   int64
	Data Schema
}

type InsertParams[Schema SchemaProps] struct {
	Document Schema
	Language tokenizer.Language
}

type InsertBatchParams[Schema SchemaProps] struct {
	Documents []Schema
	BatchSize int
	Language  tokenizer.Language
}

type UpdateParams[Schema SchemaProps] struct {
	Id       int64
	Document Schema
	Language tokenizer.Language
}

type DeleteParams[Schema SchemaProps] struct {
	Id       int64
	Language tokenizer.Language
}

type Params struct {
	Extra      map[string]any     `json:"extra"`
	Query      string             `json:"query"`
	Properties []string           `json:"properties"`
	BoolMode   Mode               `json:"boolMode"`
	Exact      bool               `json:"exact"`
	Tolerance  int                `json:"tolerance"`
	Relevance  BM25Params         `json:"relevance"`
	Paginate   bool               `json:"paginate"`
	Offset     int                `json:"offset"`
	Limit      int                `json:"limit"`
	Language   tokenizer.Language `json:"lang"`
}

func (p *Params) ToInt64() uint64 {
	bt, err := json.Marshal(p)
	if err != nil {
		return 0
	}
	f := fnv.New64()
	f.Write(bt)
	return f.Sum64()
}

type BM25Params struct {
	K float64 `json:"k"`
	B float64 `json:"b"`
	D float64 `json:"d"`
}

type Result[Schema SchemaProps] struct {
	Hits  Hits[Schema]
	Count int
}

type Hit[Schema SchemaProps] struct {
	Id    int64
	Data  Schema
	Score float64
}

type Hits[Schema SchemaProps] []Hit[Schema]

func (r Hits[Schema]) Len() int { return len(r) }

func (r Hits[Schema]) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

func (r Hits[Schema]) Less(i, j int) bool { return r[i].Score > r[j].Score }

type Config struct {
	Key             string
	DefaultLanguage tokenizer.Language
	TokenizerConfig *tokenizer.Config
	IndexKeys       []string
	Rules           map[string]bool
	SliceField      string
}

type Search[Schema SchemaProps] struct {
	mutex           sync.RWMutex
	documents       map[int64]Schema
	indexes         map[string]*Index
	indexKeys       []string
	defaultLanguage tokenizer.Language
	tokenizerConfig *tokenizer.Config
	rules           map[string]bool
	cache           map[uint64]map[int64]float64
	key             string
	sliceField      string
}

func New[Schema SchemaProps](c *Config) *Search[Schema] {
	if c.TokenizerConfig == nil {
		c.TokenizerConfig = &tokenizer.Config{
			EnableStemming:  true,
			EnableStopWords: true,
		}
	}
	if c.DefaultLanguage == "" {
		c.DefaultLanguage = tokenizer.ENGLISH
	}
	db := &Search[Schema]{
		key:             c.Key,
		documents:       make(map[int64]Schema),
		indexes:         make(map[string]*Index),
		indexKeys:       make([]string, 0),
		defaultLanguage: c.DefaultLanguage,
		tokenizerConfig: c.TokenizerConfig,
		rules:           c.Rules,
		sliceField:      c.SliceField,
	}
	db.buildIndexes()
	if len(db.indexKeys) == 0 {
		for _, key := range c.IndexKeys {
			db.indexes[key] = NewIndex()
			db.indexKeys = append(db.indexKeys, key)
		}
	}
	return db
}

func (db *Search[Schema]) buildIndexes() {
	var s Schema
	for key := range db.flattenSchema(s) {
		db.indexes[key] = NewIndex()
		db.indexKeys = append(db.indexKeys, key)
	}
}

func (db *Search[Schema]) DocumentLen() int {
	return len(db.documents)
}

func (db *Search[Schema]) Insert(doc Schema, lang ...tokenizer.Language) (Record[Schema], error) {
	language := tokenizer.ENGLISH
	if len(lang) > 0 {
		language = lang[0]
	}
	id := xid.New().Int64()
	document := db.flattenSchema(doc)

	if language == "" {
		language = db.defaultLanguage

	} else if !tokenizer.IsSupportedLanguage(language) {
		return Record[Schema]{}, fmt.Errorf("not supported language")
	}

	db.mutex.Lock()
	defer db.mutex.Unlock()

	if _, ok := db.documents[id]; ok {
		return Record[Schema]{}, fmt.Errorf("document id already exists")
	}

	db.documents[id] = doc
	db.indexDocument(id, document, language)
	return Record[Schema]{Id: id, Data: doc}, nil
}

func (db *Search[Schema]) InsertBatch(docs []Schema, batchSize int, lang ...tokenizer.Language) []error {
	docLen := len(docs)
	if docLen == 0 {
		return nil
	}
	if len(db.indexKeys) == 0 {
		keys := DocFields(docs[0])
		for _, key := range keys {
			db.indexes[key] = NewIndex()
			db.indexKeys = append(db.indexKeys, key)
		}
	}
	batchCount := int(math.Ceil(float64(len(docs)) / float64(batchSize)))
	docsChan := make(chan Schema)
	errsChan := make(chan error)
	language := tokenizer.ENGLISH
	if len(lang) > 0 {
		language = lang[0]
	}
	var wg sync.WaitGroup
	wg.Add(batchCount)

	go func() {
		for _, doc := range docs {
			docsChan <- doc
		}
		close(docsChan)
	}()

	for i := 0; i < batchCount; i++ {
		go func() {
			defer wg.Done()
			for doc := range docsChan {
				if _, err := db.Insert(doc, language); err != nil {
					errsChan <- err
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errsChan)
	}()

	errs := make([]error, 0)
	for err := range errsChan {
		errs = append(errs, err)
	}
	return errs
}

func (db *Search[Schema]) Update(params *UpdateParams[Schema]) (Record[Schema], error) {
	document := db.flattenSchema(params.Document)

	language := params.Language
	if language == "" {
		language = db.defaultLanguage

	} else if !tokenizer.IsSupportedLanguage(language) {
		return Record[Schema]{}, fmt.Errorf("not supported language")
	}

	db.mutex.Lock()
	defer db.mutex.Unlock()

	oldDocument, ok := db.documents[params.Id]
	if !ok {
		return Record[Schema]{}, fmt.Errorf("document not found")
	}

	db.indexDocument(params.Id, document, language)
	document = db.flattenSchema(oldDocument)
	db.deindexDocument(params.Id, document, language)

	db.documents[params.Id] = params.Document

	return Record[Schema]{Id: params.Id, Data: params.Document}, nil
}

func (db *Search[Schema]) Delete(params *DeleteParams[Schema]) error {
	language := params.Language
	if language == "" {
		language = db.defaultLanguage

	} else if !tokenizer.IsSupportedLanguage(language) {
		return fmt.Errorf("not supported language")
	}

	db.mutex.Lock()
	defer db.mutex.Unlock()

	document, ok := db.documents[params.Id]
	if !ok {
		return fmt.Errorf("document not found")
	}

	doc := db.flattenSchema(document)
	db.deindexDocument(params.Id, doc, language)

	delete(db.documents, params.Id)

	return nil
}

func (db *Search[Schema]) prepareResult(idScores map[int64]float64, params *Params) (Result[Schema], error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	results := make(Hits[Schema], 0)

	for id, score := range idScores {
		if doc, ok := db.documents[id]; ok {
			results = append(results, Hit[Schema]{Id: id, Data: doc, Score: score})
		}
	}

	sort.Sort(results)

	if params.Paginate {
		if params.Limit == 0 {
			params.Limit = 20
		}
		start, stop := lib.Paginate(params.Offset, params.Limit, len(results))
		return Result[Schema]{Hits: results[start:stop], Count: len(results)}, nil
	}
	return Result[Schema]{Hits: results, Count: len(results)}, nil
}

func (db *Search[Schema]) ClearCache() {
	db.cache = nil
}

func (db *Search[Schema]) prepareParams(params *Params) (map[int64]float64, error) {
	allIdScores := make(map[int64]float64)

	properties := params.Properties
	if len(params.Properties) == 0 {
		properties = db.indexKeys
	}
	language := params.Language
	if language == "" {
		language = db.defaultLanguage
	} else if !tokenizer.IsSupportedLanguage(language) {
		return nil, fmt.Errorf("not supported language")
	}
	if language == "" {
		language = tokenizer.ENGLISH
	}
	tokens, _ := tokenizer.Tokenize(&tokenizer.TokenizeParams{
		Text:            params.Query,
		Language:        language,
		AllowDuplicates: false,
	}, db.tokenizerConfig)

	db.mutex.RLock()
	defer db.mutex.RUnlock()
	for _, prop := range properties {
		if index, ok := db.indexes[prop]; ok {
			idScores := index.Find(&FindParams{
				Tokens:    tokens,
				BoolMode:  params.BoolMode,
				Exact:     params.Exact,
				Tolerance: params.Tolerance,
				Relevance: params.Relevance,
				DocsCount: len(db.documents),
			})
			for id, score := range idScores {
				allIdScores[id] += score
			}
		}
	}
	return allIdScores, nil
}

func (db *Search[Schema]) Search(params *Params) (Result[Schema], error) {
	if db.cache == nil {
		db.cache = make(map[uint64]map[int64]float64)
	}
	cachedKey := params.ToInt64()
	if cachedKey != 0 {
		if score, ok := db.cache[cachedKey]; ok && len(score) > 0 {
			return db.prepareResult(score, params)
		}
	}
	if params.Query == "" && len(params.Extra) > 0 {
		for key, val := range params.Extra {
			params.Query = fmt.Sprintf("%v", val)
			params.Properties = append(params.Properties, key)
			delete(params.Extra, key)
			break
		}
	}
	allIdScores, err := db.prepareParams(params)
	if err != nil {
		return Result[Schema]{}, err
	}
	if len(params.Extra) > 0 {
		idScores := make(map[int64]float64)
		commonKeys := make(map[string][]int64)
		for key, val := range params.Extra {
			param := &Params{
				Query:      fmt.Sprintf("%v", val),
				Properties: []string{key},
				BoolMode:   params.BoolMode,
				Exact:      true,
				Tolerance:  params.Tolerance,
				Relevance:  params.Relevance,
				Language:   params.Language,
			}
			scores, err := db.prepareParams(param)
			if err != nil {
				return Result[Schema]{}, err
			}
			for id, _ := range scores {
				if v, k := allIdScores[id]; k {
					idScores[id] = v
					commonKeys[key] = append(commonKeys[key], id)
				}
			}
			var keys [][]int64
			for _, k := range commonKeys {
				keys = append(keys, k)
			}
			if len(keys) > 0 {
				d := utils.Intersection(keys...)
				for id, _ := range idScores {
					if !str.Contains(d, id) {
						delete(idScores, id)
					}
				}
			}
		}
		commonKeys = nil
		if cachedKey != 0 {
			db.cache[cachedKey] = idScores
		}
		return db.prepareResult(idScores, params)
	}
	return db.prepareResult(allIdScores, params)
}

func (db *Search[Schema]) indexDocument(id int64, document map[string]string, language tokenizer.Language) {
	for propName, index := range db.indexes {
		tokens, _ := tokenizer.Tokenize(&tokenizer.TokenizeParams{
			Text:            document[propName],
			Language:        language,
			AllowDuplicates: true,
		}, db.tokenizerConfig)

		index.Insert(&IndexParams{
			Id:        id,
			Tokens:    tokens,
			DocsCount: len(db.documents),
		})
	}
}

func (db *Search[Schema]) deindexDocument(id int64, document map[string]string, language tokenizer.Language) {
	for propName, index := range db.indexes {
		tokens, _ := tokenizer.Tokenize(&tokenizer.TokenizeParams{
			Text:            document[propName],
			Language:        language,
			AllowDuplicates: false,
		}, db.tokenizerConfig)

		index.Delete(&IndexParams{
			Id:        id,
			Tokens:    tokens,
			DocsCount: len(db.documents),
		})
	}
}

func (db *Search[Schema]) getFieldsFromMap(obj map[string]any, prefix ...string) map[string]string {
	fields := make(map[string]string)
	rules := make(map[string]bool)
	if db.rules != nil {
		rules = db.rules
	}
	for field, val := range obj {
		if reflect.TypeOf(field).Kind() == reflect.Map {
			for key, value := range db.flattenSchema(val, field) {
				fields[key] = value
			}
		} else {
			if len(rules) > 0 {
				if canIndex, ok := rules[field]; ok && canIndex {
					fields[field] = fmt.Sprintf("%v", val)
				}
			} else {
				fields[field] = fmt.Sprintf("%v", val)
			}
		}
	}
	return fields
}

func (db *Search[Schema]) getFieldsFromStruct(obj any, prefix ...string) map[string]string {
	fields := make(map[string]string)
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	visibleFields := reflect.VisibleFields(t)
	hasIndexField := false
	for i, field := range visibleFields {
		if propName, ok := field.Tag.Lookup("index"); ok {
			hasIndexField = true
			if len(prefix) == 1 {
				propName = fmt.Sprintf("%s.%s", prefix[0], propName)
			}
			if field.Type.Kind() == reflect.Struct {
				for key, value := range db.flattenSchema(v.Field(i).Interface(), propName) {
					fields[key] = value
				}
			} else {
				fields[propName] = v.Field(i).String()
			}
		}
	}
	if !hasIndexField {
		for i, field := range visibleFields {
			propName := field.Name
			if len(prefix) == 1 {
				propName = fmt.Sprintf("%s.%s", prefix[0], propName)
			}

			if field.Type.Kind() == reflect.Struct {
				for key, value := range db.flattenSchema(v.Field(i).Interface(), propName) {
					fields[key] = value
				}
			} else {
				fields[propName] = v.Field(i).String()
			}
		}
	}
	return fields
}

func (db *Search[Schema]) flattenSchema(obj any, prefix ...string) map[string]string {
	if obj == nil {
		return nil
	}
	fields := make(map[string]string)
	if reflect.TypeOf(obj).Kind() == reflect.Struct {
		return db.getFieldsFromStruct(obj, prefix...)
	} else {
		switch obj := obj.(type) {
		case string, bool, time.Time, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			fields[db.sliceField] = fmt.Sprintf("%v", obj)
			return fields
		case map[string]any:
			return db.getFieldsFromMap(obj, prefix...)
		default:
			switch obj := obj.(type) {
			case map[string]any:
				return db.getFieldsFromMap(obj, prefix...)
			default:
				return db.getFieldsFromStruct(obj, prefix...)
			}
		}
	}
}

func getFieldsFromMap(obj map[string]any) []string {
	var fields []string
	rules := make(map[string]bool)
	for field, val := range obj {
		if reflect.TypeOf(field).Kind() == reflect.Map {
			for _, key := range DocFields(val, field) {
				fields = append(fields, key)
			}
		} else {
			if len(rules) > 0 {
				if canIndex, ok := rules[field]; ok && canIndex {
					fields = append(fields, field)
				}
			} else {
				fields = append(fields, field)
			}
		}
	}
	return fields
}
func getFieldsFromStruct(obj any, prefix ...string) []string {
	var fields []string
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	visibleFields := reflect.VisibleFields(t)
	hasIndexField := false
	for i, field := range visibleFields {
		if propName, ok := field.Tag.Lookup("index"); ok {
			hasIndexField = true
			if len(prefix) == 1 {
				propName = fmt.Sprintf("%s.%s", prefix[0], propName)
			}

			if field.Type.Kind() == reflect.Struct {
				for _, key := range DocFields(v.Field(i).Interface(), propName) {
					fields = append(fields, key)
				}
			} else {
				fields = append(fields, propName)
			}
		}
	}

	if !hasIndexField {
		for i, field := range visibleFields {
			propName := field.Name
			if len(prefix) == 1 {
				propName = fmt.Sprintf("%s.%s", prefix[0], propName)
			}

			if field.Type.Kind() == reflect.Struct {
				for _, key := range DocFields(v.Field(i).Interface(), propName) {
					fields = append(fields, key)
				}
			} else {
				fields = append(fields, propName)
			}
		}
	}
	return fields
}

func DocFields(obj any, prefix ...string) []string {
	if obj == nil {
		return nil
	}

	switch obj := obj.(type) {
	case map[string]any:
		return getFieldsFromMap(obj)
	case map[string]string:
		data := make(map[string]any)
		for k, v := range obj {
			data[k] = v
		}
		return getFieldsFromMap(data)
	default:
		switch obj := obj.(type) {
		case map[string]any:
			return getFieldsFromMap(obj)
		case map[string]string:
			data := make(map[string]any)
			for k, v := range obj {
				data[k] = v
			}
			return getFieldsFromMap(data)
		default:
			return getFieldsFromStruct(obj, prefix...)
		}
	}
}
