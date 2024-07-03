package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"-"`
}

type FieldMapping struct {
	Source string
	Target string
}

type MapperConfig struct {
	Mappings            []FieldMapping
	AllowUnmappedFields bool
}

// Map maps data from src to dst using the provided mapperConfig.
// If mapperConfig is nil or has AllowUnmappedFields set to true, it will try to map common fields.
func Map(src interface{}, dst interface{}, cfg ...*MapperConfig) error {
	var mapperConfig *MapperConfig
	if len(cfg) > 0 {
		mapperConfig = cfg[0]
	}
	srcVal := reflect.ValueOf(src)
	srcType := reflect.TypeOf(src)
	dstVal := reflect.ValueOf(dst).Elem()
	dstType := reflect.TypeOf(dst).Elem()

	if srcType.Kind() == reflect.Slice && dstType.Kind() == reflect.Slice {
		// Handle mapping between slices
		return mapSlice(srcVal, dstVal, mapperConfig)
	} else if srcType.Kind() == reflect.Map && dstType.Kind() == reflect.Map {
		// Handle mapping between maps
		return mapMap(srcVal, dstVal, mapperConfig)
	} else if srcType.Kind() == reflect.Struct && dstType.Kind() == reflect.Struct {
		// Handle mapping between structs
		return mapStruct(srcVal, dstVal, mapperConfig)
	} else {
		return fmt.Errorf("unsupported types: %v -> %v", srcType.Kind(), dstType.Kind())
	}
}

// mapSlice maps data from a slice of src to a slice of dst using the provided mapperConfig.
func mapSlice(src reflect.Value, dst reflect.Value, mapperConfig *MapperConfig) error {
	srcType := src.Type().Elem()
	dstType := dst.Type().Elem()

	if srcType.Kind() != reflect.Struct || dstType.Kind() != reflect.Struct {
		return fmt.Errorf("slice elements must be structs")
	}

	for i := 0; i < src.Len(); i++ {
		srcElem := src.Index(i)
		dstElem := reflect.New(dstType).Elem()

		if err := mapStruct(srcElem, dstElem, mapperConfig); err != nil {
			return err
		}

		dst.Set(reflect.Append(dst, dstElem))
	}

	return nil
}

// mapMap maps data from a map of src to a map of dst using the provided mapperConfig.
func mapMap(src reflect.Value, dst reflect.Value, mapperConfig *MapperConfig) error {
	srcKeyType := src.Type().Key()
	srcValueType := src.Type().Elem()
	dstKeyType := dst.Type().Key()
	dstValueType := dst.Type().Elem()

	if srcKeyType != dstKeyType || srcValueType.Kind() != reflect.Interface || dstValueType.Kind() != reflect.Interface {
		return fmt.Errorf("map key or value types are not compatible")
	}

	dst.Set(reflect.MakeMap(dst.Type()))

	for _, key := range src.MapKeys() {
		srcVal := src.MapIndex(key)
		dstVal := reflect.New(dstValueType).Elem()

		if err := mapValue(srcVal, dstVal, mapperConfig); err != nil {
			return err
		}

		dst.SetMapIndex(key, dstVal)
	}

	return nil
}

// mapStruct maps data from src to dst using the provided mapperConfig.
func mapStruct(src reflect.Value, dst reflect.Value, mapperConfig *MapperConfig) error {
	srcType := src.Type()
	dstType := dst.Type()

	if srcType.Kind() != reflect.Struct || dstType.Kind() != reflect.Struct {
		return fmt.Errorf("source and destination must be structs")
	}

	// Create a map for quick lookup of mappings
	mappingMap := make(map[string]string)
	if mapperConfig != nil {
		for _, mapping := range mapperConfig.Mappings {
			mappingMap[mapping.Source] = mapping.Target
		}
	}

	for i := 0; i < src.NumField(); i++ {
		srcField := srcType.Field(i)
		dstFieldName := srcField.Name
		if mappedName, mapped := mappingMap[srcField.Name]; mapped {
			dstFieldName = mappedName
		} else if mapperConfig != nil && !mapperConfig.AllowUnmappedFields {
			continue
		}

		dstField := dst.FieldByName(dstFieldName)
		if !dstField.IsValid() {
			return fmt.Errorf("field %s does not exist in destination struct", dstFieldName)
		}

		if srcField.Type != dstField.Type() {
			return fmt.Errorf("type mismatch for field %s", srcField.Name)
		}

		dstField.Set(src.Field(i))
	}

	return nil
}

// mapValue maps a single value from src to dst using the provided mapperConfig.
func mapValue(src reflect.Value, dst reflect.Value, mapperConfig *MapperConfig) error {
	srcType := src.Type()
	dstType := dst.Type()

	if srcType.Kind() != dstType.Kind() {
		return fmt.Errorf("type mismatch for map value: %v -> %v", srcType.Kind(), dstType.Kind())
	}

	switch srcType.Kind() {
	case reflect.Struct:
		return mapStruct(src, dst, mapperConfig)
	case reflect.Slice:
		return mapSlice(src, dst, mapperConfig)
	default:
		dst.Set(src)
	}

	return nil
}

func main() {
	data := []byte(`[
		{"username":"test_user1", "password":"test_pass1"},
		{"username":"test_user2", "password":"test_pass2"}
	]`)
	var loginSlice []Login
	if err := json.Unmarshal(data, &loginSlice); err != nil {
		log.Fatalf("Error unmarshalling data: %v", err)
	}

	var userSlice []User

	// Example with custom mapper config
	mapperConfig := &MapperConfig{
		Mappings: []FieldMapping{
			{Source: "Username", Target: "Username"},
			// Password mapping is not added to demonstrate handling of unmapped fields
		},
		AllowUnmappedFields: true,
	}

	if err := Map(loginSlice, &userSlice, mapperConfig); err != nil {
		log.Fatalf("Error mapping data: %v", err)
	}

	fmt.Println("Mapped userSlice with custom mapper config: %+v\n", userSlice)

	// Example with no mapper config provided
	if err := Map(loginSlice, &userSlice); err != nil {
		log.Fatalf("Error mapping data: %v", err)
	}

	fmt.Println("Mapped userSlice without mapper config: %+v\n", userSlice)
}
