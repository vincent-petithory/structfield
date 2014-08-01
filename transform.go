package structfield

import (
	"reflect"
	"strings"
)

// Transformer is the interface which provides the Transform method.
//
// Transform takes a field, whose name is field and value is value,
// and returns a new field name and new value.
type Transformer interface {
	Transform(field string, value interface{}) (string, interface{})
}

// Transform takes a struct or ptr to a struct and converts it to a map,
// applying a set of transformers on it.
//
// transformers is a map of field names to Transformer.
// Before the transformers are run, Transform honors the json struct tag of the field.
//
// When the Transform method of a Transformer returns an empty string for the new field name, the field is discarded.
//
// Note that Transform stays at depth 1: only the fields of the struct are processed, not nested data structures.
// e.g: T.Foo is processed, but T.Foo.Bar isn't.
func Transform(v interface{}, transformers map[string]Transformer) map[string]interface{} {
	rv := reflect.ValueOf(v)
KindTest:
	for {
		switch rv.Kind() {
		case reflect.Ptr:
			rv = rv.Elem()
		case reflect.Struct:
			break KindTest
		default:
			panic("structfield.Transform: not a ptr or a struct")
		}
	}
	rt := rv.Type()
	nFields := rv.NumField()
	m := make(map[string]interface{}, nFields)
	for i := 0; i < nFields; i++ {
		var fieldName string
		tagParts := strings.Split(rt.Field(i).Tag.Get("json"), ",")
		if len(tagParts) > 0 {
			fieldName = strings.TrimSpace(tagParts[0])
		}
		if fieldName == "-" {
			continue
		}
		if fieldName == "" {
			fieldName = rt.Field(i).Name
		}

		var omitempty bool
		for _, tagPart := range tagParts[1:] {
			if tagPart == "omitempty" {
				omitempty = true
				break
			}
		}
		if omitempty && isEmptyValue(rv.Field(i)) {
			continue
		}
		transformer, ok := transformers[fieldName]
		if !ok {
			m[fieldName] = rv.Field(i).Interface()
			continue
		}
		newField, newValue := transformer.Transform(fieldName, rv.Field(i).Interface())
		if newField != "" {
			m[newField] = newValue
		}
	}
	return m
}

// Borrowed from encoding/json.
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

// TransformerFunc is to Transformer what http.HandlerFunc is to http.Handler.
type TransformerFunc func(field string, value interface{}) (string, interface{})

// Transform calls tf(field, value).
func (tf TransformerFunc) Transform(field string, value interface{}) (string, interface{}) {
	return tf(field, value)
}
