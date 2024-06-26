package dipper

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// setOption is a type used for special assignments in a set operation.
type setOption int

const (
	// Zero is used as the new value in Set() to set the attribute to its zero
	// value (e.g. "" for string, nil for any, etc.).
	Zero setOption = 0
	// Delete is used as the new value in Set() to delete a map key. If the
	// field is not a map value, the value will be zeroed (see Zero).
	Delete setOption = 1
)

// Options defines the configuration of a Dipper instance.
type Options struct {
	Separator string
	Slice     string
	Wildcard  string
}

// Dipper allows to access deeply-nested object attributes to get or set their
// values. Attributes are specified by a string with its fields separated by
// some delimiter (e.g. “Books.3.Author" or "Books->3->Author", with "." and
// "->" as delimiters, respectively).
type Dipper struct {
	separator string
	slice     string
	wildcard  string
}

// New returns a new Dipper instance.
func New(opts Options) *Dipper {
	return &Dipper{separator: opts.Separator, slice: opts.Separator + opts.Slice + opts.Separator, wildcard: opts.Separator + opts.Wildcard + opts.Separator}
}

// Get returns the value of the given obj attribute. The attribute uses some
// delimiter-notation to allow accessing nested fields, slice elements or map
// keys. Field names and key maps are case-sensitive.
// All the struct fields accessed must be exported.
// If an error occurs, it will be returned as the attribute value, so it should
// be handled. All the returned errors are fieldError.
//
// Example:
//
//	 // Using "." as the Dipper separator
//		v := my_dipper.Get(myObj, "SomeStructField.1.some_key_map")
//		if err := Error(v); err != nil {
//		    return err
//		}
func (d *Dipper) Get(obj any, attribute string) any {
	if !strings.Contains(attribute, d.slice) {
		value, _, err := getReflectValue(reflect.ValueOf(obj), attribute, d.separator, false)
		if err != nil {
			return err
		}
		return value.Interface()
	}
	fields := strings.Split(attribute, d.slice)
	left := fields[0]
	right := fields[1]
	value, _, err := getReflectValue(reflect.ValueOf(obj), left, d.separator, false)
	if err != nil {
		return err
	}
	var values []any
	switch v := value.Interface().(type) {
	case []map[string]any:
		for _, vt := range v {
			val := d.Get(vt, right)
			values = append(values, val)
		}
		return values
	case []any:
		for _, vt := range v {
			val := d.Get(vt, right)
			values = append(values, val)
		}
		return values
	default:
		value, _, err := getReflectValue(reflect.ValueOf(obj), attribute, d.separator, false)
		if err != nil {
			return err
		}
		return value.Interface()
	}
}

func (d *Dipper) FilterSlice(obj any, attribute string, search []any) any {
	if !strings.Contains(attribute, d.slice) {
		value, _, err := getReflectValue(reflect.ValueOf(obj), attribute, d.separator, false)
		if err != nil {
			return err
		}
		switch i := value.Interface().(type) {
		case []any:
			var values []any
			for _, it := range i {
				for _, s := range search {
					if s == it {
						values = append(values, it)
					}
				}
			}
			return values
		case []string:
			var values []any
			for _, it := range i {
				for _, s := range search {
					if s == it {
						values = append(values, it)
					}
				}
			}
			return values
		case []float32:
			var values []any
			for _, it := range i {
				for _, s := range search {
					if s == it {
						values = append(values, it)
					}
				}
			}
			return values
		case []float64:
			var values []any
			for _, it := range i {
				for _, s := range search {
					switch s := s.(type) {
					case int:
						if float64(s) == it {
							values = append(values, it)
						}
					case float64:
						if s == it {
							values = append(values, it)
						}
					case float32:
						if float64(s) == it {
							values = append(values, it)
						}
					}

				}
			}
			return values
		}
		return value.Interface()
	}
	fields := strings.Split(attribute, d.slice)
	left := fields[0]
	right := fields[1]
	value, _, err := getReflectValue(reflect.ValueOf(obj), left, d.separator, false)
	if err != nil {
		return err
	}
	var values []any
	switch v := value.Interface().(type) {
	case []map[string]any:
		for _, vt := range v {
			val := d.Get(vt, right)
			for _, t := range search {
				if fmt.Sprintf("%v", val) == fmt.Sprintf("%v", t) {
					values = append(values, vt)
				}
			}
		}
		return values
	case []any:
		for _, vt := range v {
			val := d.Get(vt, right)
			for _, t := range search {
				if fmt.Sprintf("%v", val) == fmt.Sprintf("%v", t) {
					values = append(values, vt)
				}
			}
		}
		return values
	default:
		value, _, err := getReflectValue(reflect.ValueOf(obj), attribute, d.separator, false)
		if err != nil {
			return err
		}
		return value.Interface()
	}
}

// GetMany returns a map with the values of the given obj attributes.
// It works as Dipper.Get(), but it takes a slice of attributes to return their
// corresponding values. The returned map will have the same length as the
// attributes slice, with the attributes as keys.
//
// Example:
//
//	 // Using "." as the Dipper separator
//		v := my_dipper.GetMany(myObj, []string{"Name", "Age", "Skills.skydiving})
//		if err := v.FirstError(); err != nil {
//		    return err
//		}
func (d *Dipper) GetMany(obj any, attributes []string) Fields {
	m := make(Fields, len(attributes))

	for _, attr := range attributes {
		if _, ok := m[attr]; !ok {
			m[attr] = d.Get(obj, attr)
		}
	}

	return m
}

// Set sets the value of the given obj attribute to the new provided value.
// The attribute uses some delimiter-notation to allow accessing nested fields,
// slice elements or map keys. Field names and key maps are case-sensitive.
// All the struct fields accessed must be exported.
// ErrUnaddressable will be returned if obj is not addressable.
// It returns nil if the value was successfully set, otherwise it will return
// a fieldError.
//
// Example:
//
//	 // Using "." as the Dipper separator
//		v := my_dipper.Set(&myObj, "SomeStructField.1.some_key_map", 123)
//		if err != nil {
//		    return err
//		}
func (d *Dipper) Set(obj any, attribute string, new any) error {
	var err error

	value := reflect.ValueOf(obj)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	var lastField string
	value, lastField, err = getReflectValue(value, attribute, d.separator, true)
	if err != nil {
		return err
	}

	var optZero, optDelete bool

	var newValue reflect.Value
	switch new {
	case Zero:
		optZero = true
	case Delete:
		optDelete = true
	default:
		newValue = reflect.ValueOf(new)
		if newValue.Kind() == reflect.Ptr {
			newValue = newValue.Elem()
		}
	}

	if value.Kind() == reflect.Map {
		if !optZero && !optDelete {
			mapValueType := value.Type().Elem()
			if mapValueType.Kind() != reflect.Interface && mapValueType != newValue.Type() {
				return ErrTypesDoNotMatch
			}
		}

		// Initialize map if needed
		if value.IsNil() {
			keyType := value.Type().Key()
			valueType := value.Type().Elem()
			mapType := reflect.MapOf(keyType, valueType)
			value.Set(reflect.MakeMapWithSize(mapType, 0))
		}

		value.SetMapIndex(reflect.ValueOf(lastField), newValue)
	} else {
		if !optZero && !optDelete {
			if !value.CanAddr() {
				return ErrUnaddressable
			}
			if value.Kind() != reflect.Interface && value.Type() != newValue.Type() {
				return ErrTypesDoNotMatch
			}
		} else {
			newValue = reflect.Zero(value.Type())
		}
		value.Set(newValue)
	}
	return nil
}

// getReflectValue gets the reflect.Value of the given value attribute.
// It splits the attribute into the field names, map keys and slice indexes
// and uses reflection to get the final value.
// toSet indicates that the function must return a value that will be set to
// another value, which is used in the special case of maps (maps elements are
// not addressable).
// It also returns the name of the accessed field.
func getReflectValue(value reflect.Value, attribute string, sep string, toSet bool) (_ reflect.Value, fieldName string, _ error) {
	if attribute == "" {
		return value, "", nil
	}

	if len(sep) == 0 {
		sep = "."
	}

	var i, maxSetDepth int
	if toSet {
		maxSetDepth = strings.Count(attribute, sep)
	}

	splitter := newAttributeSplitter(attribute, sep)
	for splitter.HasMore() {
		fieldName, i = splitter.Next()

		if value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
			value = value.Elem()
		}

		switch value.Kind() {
		case reflect.Map:
			// Check that the map accept string keys
			keyKind := value.Type().Key().Kind()
			if keyKind != reflect.String && keyKind != reflect.Interface {
				return value, "", ErrMapKeyNotString
			}

			// If a map key has to be set, skip the last attribute and return the map
			if toSet && i == maxSetDepth {
				return value, fieldName, nil
			}

			mapValue := value.MapIndex(reflect.ValueOf(fieldName))
			if !mapValue.IsValid() {
				return value, "", ErrNotFound
			}

			value = mapValue

		case reflect.Struct:
			field, ok := value.Type().FieldByName(fieldName)
			if !ok {
				return value, "", ErrNotFound
			}
			// Check if field is unexported (method IsExported() was introduced in Go 1.17)
			if field.PkgPath != "" {
				return value, "", ErrUnexported
			}

			value = value.FieldByName(fieldName)

		case reflect.Slice, reflect.Array:
			sliceIndex, err := strconv.Atoi(fieldName)
			if err != nil {
				return value, "", ErrInvalidIndex
			}
			if sliceIndex < 0 || sliceIndex >= value.Len() {
				return value, "", ErrIndexOutOfRange
			}
			field := value.Index(sliceIndex)
			value = field

		default:
			return value, "", ErrNotFound
		}
	}

	return value, fieldName, nil
}
