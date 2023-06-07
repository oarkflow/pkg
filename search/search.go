package search

import (
	"fmt"
	"hash/fnv"
	"math"
	"reflect"
	"sort"
	"sync"
	"time"

	json "github.com/bytedance/sonic"
	"github.com/oarkflow/xid"

	"github.com/oarkflow/pkg/search/lib"
	"github.com/oarkflow/pkg/search/radix"
	"github.com/oarkflow/pkg/search/tokenizer"
)

const (
	AND Mode = "AND"
	OR  Mode = "OR"
)

const WILDCARD = "*"

type Mode string

type SchemaProps any

type Record[Schema SchemaProps] struct {
	Data Schema
	Id   int64
}

type UpdateParams[Schema SchemaProps] struct {
	Document Schema
	Language tokenizer.Language
	Id       int64
}

type DeleteParams[Schema SchemaProps] struct {
	Language tokenizer.Language
	Id       int64
}

type indexParams struct {
	document map[string]string
	language tokenizer.Language
	id       int64
}

type findParams struct {
	boolMode   Mode
	language   tokenizer.Language
	query      string
	properties []string
	exact      bool
	tolerance  int
}

type Params struct {
	Extra      map[string]any     `json:"extra"`
	BoolMode   Mode               `json:"bool_mode"`
	Language   tokenizer.Language `json:"language"`
	Query      string             `json:"query"`
	Properties []string           `json:"properties"`
	Exact      bool               `json:"exact"`
	Paginate   bool               `json:"paginate"`
	Tolerance  int                `json:"tolerance"`
	Offset     int                `json:"offset"`
	Limit      int                `json:"limit"`
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

type Result[Schema SchemaProps] struct {
	Hits  Hits[Schema]
	Count int
}

type Hit[Schema SchemaProps] struct {
	Data  Schema
	Id    int64
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
	Rules           map[string]bool
	SliceField      string
}

type Store[Schema SchemaProps] struct {
	tokenizerConfig *tokenizer.Config
	rules           map[string]bool
	occurrences     map[string]map[string]int
	cache           map[uint64]map[int64]float64
	documents       map[int64]Schema
	indexes         map[string]*radix.Trie
	defaultLanguage tokenizer.Language
	key             string
	sliceField      string
	indexKeys       []string
	mutex           sync.RWMutex
}

func New[Schema SchemaProps](c *Config) *Store[Schema] {
	if c.DefaultLanguage == "" {
		c.DefaultLanguage = tokenizer.ENGLISH
	}
	db := &Store[Schema]{
		key:             c.Key,
		documents:       make(map[int64]Schema),
		indexes:         make(map[string]*radix.Trie),
		indexKeys:       make([]string, 0),
		occurrences:     make(map[string]map[string]int),
		defaultLanguage: c.DefaultLanguage,
		tokenizerConfig: c.TokenizerConfig,
		rules:           c.Rules,
		sliceField:      c.SliceField,
	}
	db.buildIndexes()
	return db
}

func (db *Store[Schema]) Insert(doc Schema, lang ...tokenizer.Language) (Record[Schema], error) {
	language := tokenizer.ENGLISH
	if len(lang) > 0 {
		language = lang[0]
	}
	if len(db.indexKeys) == 0 {
		for key := range db.flattenSchema(doc) {
			db.mutex.Lock()
			db.indexes[key] = radix.New()
			db.indexKeys = append(db.indexKeys, key)
			db.occurrences[key] = make(map[string]int)
			db.mutex.Unlock()
		}
	}
	idxParams := indexParams{
		id:       xid.New().Int64(),
		document: db.flattenSchema(doc),
		language: language,
	}
	if idxParams.language == "" {
		idxParams.language = db.defaultLanguage

	} else if !tokenizer.IsSupportedLanguage(idxParams.language) {
		return Record[Schema]{}, fmt.Errorf("not supported language")
	}

	db.mutex.Lock()
	defer db.mutex.Unlock()

	if _, ok := db.documents[idxParams.id]; ok {
		return Record[Schema]{}, fmt.Errorf("document id already exists")
	}

	db.documents[idxParams.id] = doc
	db.indexDocument(&idxParams)

	return Record[Schema]{Id: idxParams.id, Data: doc}, nil
}

func (db *Store[Schema]) InsertBatch(docs []Schema, batchSize int, lang ...tokenizer.Language) []error {
	language := tokenizer.ENGLISH
	if len(lang) > 0 {
		language = lang[0]
	}
	batchCount := int(math.Ceil(float64(len(docs)) / float64(batchSize)))
	docsChan := make(chan Schema)
	errsChan := make(chan error)

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

func (db *Store[Schema]) Update(params *UpdateParams[Schema]) (Record[Schema], error) {
	idxParams := indexParams{
		id:       params.Id,
		language: params.Language,
		document: db.flattenSchema(params.Document),
	}

	if idxParams.language == "" {
		idxParams.language = db.defaultLanguage

	} else if !tokenizer.IsSupportedLanguage(idxParams.language) {
		return Record[Schema]{}, fmt.Errorf("not supported language")
	}

	db.mutex.Lock()
	defer db.mutex.Unlock()

	oldDocument, ok := db.documents[idxParams.id]
	if !ok {
		return Record[Schema]{}, fmt.Errorf("document not found")
	}

	db.indexDocument(&idxParams)
	idxParams.document = db.flattenSchema(oldDocument)
	db.deindexDocument(&idxParams)

	db.documents[idxParams.id] = params.Document

	return Record[Schema]{Id: idxParams.id, Data: params.Document}, nil
}

func (db *Store[Schema]) DocumentLen() int {
	return len(db.documents)
}

func (db *Store[Schema]) Delete(params *DeleteParams[Schema]) error {
	idxParams := indexParams{
		id:       params.Id,
		language: params.Language,
	}

	if idxParams.language == "" {
		idxParams.language = db.defaultLanguage

	} else if !tokenizer.IsSupportedLanguage(idxParams.language) {
		return fmt.Errorf("not supported language")
	}

	db.mutex.Lock()
	defer db.mutex.Unlock()

	document, ok := db.documents[idxParams.id]
	if !ok {
		return fmt.Errorf("document not found")
	}

	idxParams.document = db.flattenSchema(document)
	db.deindexDocument(&idxParams)

	delete(db.documents, idxParams.id)

	return nil
}

func (db *Store[Schema]) Search(params *Params) (Result[Schema], error) {
	if db.cache == nil {
		db.cache = make(map[uint64]map[int64]float64)
	}
	cachedKey := params.ToInt64()
	if cachedKey != 0 {
		if score, ok := db.cache[cachedKey]; ok {
			return db.prepareResult(score, params)
		}
	}
	idxParams, err := db.prepareFindParams(params)
	if err != nil {
		return Result[Schema]{}, err
	}
	queryScores := db.findDocumentIds(idxParams)
	if len(params.Extra) > 0 {
		idScores := make(map[int64]float64)
		for key, val := range params.Extra {
			param := &Params{
				Query:      fmt.Sprintf("%v", val),
				Properties: []string{key},
				BoolMode:   params.BoolMode,
				Exact:      true,
				Tolerance:  params.Tolerance,
				Language:   params.Language,
			}
			extraParams, err := db.prepareFindParams(param)
			if err != nil {
				return Result[Schema]{}, err
			}
			scores := db.findDocumentIds(extraParams)
			for id, _ := range scores {
				if v, k := queryScores[id]; k {
					idScores[id] = v
				}
			}
		}
		if cachedKey != 0 {
			db.cache[cachedKey] = idScores
		}
		return db.prepareResult(idScores, params)
	}
	if cachedKey != 0 {
		db.cache[cachedKey] = queryScores
	}
	return db.prepareResult(queryScores, params)
}

func (db *Store[Schema]) ClearCache() {
	db.cache = nil
}

func (db *Store[Schema]) prepareFindParams(params *Params) (*findParams, error) {
	if params.Language == "" {
		params.Language = tokenizer.ENGLISH
	}
	if len(params.Properties) == 0 {
		params.Properties = []string{WILDCARD}
	}

	idxParams := &findParams{
		query:      params.Query,
		properties: params.Properties,
		boolMode:   params.BoolMode,
		exact:      params.Exact,
		tolerance:  params.Tolerance,
		language:   params.Language,
	}

	if idxParams.language == "" {
		idxParams.language = db.defaultLanguage

	} else if !tokenizer.IsSupportedLanguage(idxParams.language) {
		return nil, fmt.Errorf("not supported language")
	}

	if len(idxParams.properties) == 1 && idxParams.properties[0] == WILDCARD {
		idxParams.properties = db.indexKeys
	}
	return idxParams, nil
}

func (db *Store[Schema]) prepareResult(idScores map[int64]float64, params *Params) (Result[Schema], error) {
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

func (db *Store[Schema]) buildIndexes() {
	var s Schema
	for key := range db.flattenSchema(s) {
		db.mutex.Lock()
		db.indexes[key] = radix.New()
		db.indexKeys = append(db.indexKeys, key)
		db.occurrences[key] = make(map[string]int)
		db.mutex.Unlock()
	}
}

func (db *Store[Schema]) findDocumentIds(params *findParams) map[int64]float64 {
	tokenParams := tokenizer.TokenizeParams{
		Text:            params.query,
		Language:        params.language,
		AllowDuplicates: false,
	}
	tokens, _ := tokenizer.Tokenize(&tokenParams, db.tokenizerConfig)

	idScores := make(map[int64]float64)
	for _, prop := range params.properties {
		if index, ok := db.indexes[prop]; ok {
			idTokensCount := make(map[int64]int)

			for _, token := range tokens {
				infos := index.Find(&radix.FindParams{
					Term:      token,
					Tolerance: params.tolerance,
					Exact:     params.exact,
				})
				for _, info := range infos {
					idScores[info.Id] += lib.TfIdf(info.TermFrequency, db.occurrences[prop][token], len(db.documents))
					idTokensCount[info.Id]++
				}
			}

			for id, tokensCount := range idTokensCount {
				if params.boolMode == AND && tokensCount != len(tokens) {
					delete(idScores, id)
				}
			}
		}
	}

	return idScores
}

func (db *Store[Schema]) indexDocument(params *indexParams) {
	tokenParams := tokenizer.TokenizeParams{
		Language:        params.language,
		AllowDuplicates: true,
	}

	for propName, index := range db.indexes {
		tokenParams.Text = params.document[propName]
		tokens, _ := tokenizer.Tokenize(&tokenParams, db.tokenizerConfig)
		tokensCount := lib.Count(tokens)

		for token, count := range tokensCount {
			tokenFrequency := float64(count) / float64(len(tokens))
			index.Insert(&radix.InsertParams{
				Id:            params.id,
				Word:          token,
				TermFrequency: tokenFrequency,
			})

			db.occurrences[propName][token]++
		}
	}
}

func (db *Store[Schema]) deindexDocument(params *indexParams) {
	tokenParams := tokenizer.TokenizeParams{
		Language:        params.language,
		AllowDuplicates: false,
	}

	for propName, index := range db.indexes {
		tokenParams.Text = params.document[propName]
		tokens, _ := tokenizer.Tokenize(&tokenParams, db.tokenizerConfig)

		for _, token := range tokens {
			index.Delete(&radix.DeleteParams{
				Id:   params.id,
				Word: token,
			})

			db.occurrences[propName][token]--
			if db.occurrences[propName][token] == 0 {
				delete(db.occurrences[propName], token)
			}
		}
	}
}

func (db *Store[Schema]) flattenSchema(obj any, prefix ...string) map[string]string {
	if obj == nil {
		return nil
	}
	fields := make(map[string]string)
	switch obj := obj.(type) {
	case string, bool, time.Time, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		fields[db.sliceField] = fmt.Sprintf("%v", obj)
	case any:
		switch obj := obj.(type) {
		case map[string]any:
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
		}
	case map[string]any:
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

	default:
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
	}

	return fields
}
