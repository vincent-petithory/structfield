package structfield

import (
	"reflect"
	"testing"
)

type testTransformer string

func (t testTransformer) Transform(field string, value interface{}) (string, interface{}) {
	var nv interface{}
	switch x := value.(type) {
	case string:
		nv = x + "=changed"
	default:
		nv = value
	}
	return field + string(t), nv
}

func TestTransformFields(t *testing.T) {
	tests := []struct {
		value        interface{}
		transformers map[string]Transformer
		m            map[string]interface{}
	}{
		{
			struct {
				A string
				B string
				C string
			}{
				"foo", "bar", "meow",
			},
			map[string]Transformer{
				"B": testTransformer("_url"),
			},
			map[string]interface{}{
				"A":     "foo",
				"B_url": "bar=changed",
				"C":     "meow",
			},
		},
		{
			&struct {
				A string
				B string
				C string
			}{
				"foo", "bar", "meow",
			},
			map[string]Transformer{
				"B": testTransformer("_url"),
			},
			map[string]interface{}{
				"A":     "foo",
				"B_url": "bar=changed",
				"C":     "meow",
			},
		},
		{
			struct {
				A string `json:",omitempty"`
				B string
				C string
			}{
				"foo", "bar", "meow",
			},
			map[string]Transformer{
				"A": testTransformer("_url"),
			},
			map[string]interface{}{
				"A_url": "foo=changed",
				"B":     "bar",
				"C":     "meow",
			},
		},
		{
			struct {
				A string
				B string
				C string `json:"field_c,omitempty"`
			}{
				"foo", "bar", "meow",
			},
			map[string]Transformer{
				"field_c": testTransformer("_url"),
			},
			map[string]interface{}{
				"A":           "foo",
				"B":           "bar",
				"field_c_url": "meow=changed",
			},
		},
		{
			struct {
				A string `json:"field_a"`
				B string
				C string
			}{
				"foo", "bar", "meow",
			},
			map[string]Transformer{
				"field_a": testTransformer("_url"),
			},
			map[string]interface{}{
				"field_a_url": "foo=changed",
				"B":           "bar",
				"C":           "meow",
			},
		},
		{
			struct {
				A string `json:"field_a,omitempty"`
				B string
				C string
			}{
				"", "bar", "meow",
			},
			map[string]Transformer{
				"field_a": testTransformer("_url"),
			},
			map[string]interface{}{
				"B": "bar",
				"C": "meow",
			},
		},
		{
			struct {
				A string `json:"-"`
				B string
				C string
			}{
				"foo", "bar", "meow",
			},
			map[string]Transformer{
				"field_a": testTransformer("_url"),
			},
			map[string]interface{}{
				"B": "bar",
				"C": "meow",
			},
		},
		{
			struct {
				A string `json:"-"`
				B string
				C string `json:"c"`
			}{
				"foo", "bar", "meow",
			},
			nil,
			map[string]interface{}{
				"B": "bar",
				"c": "meow",
			},
		},
	}
	for _, test := range tests {
		m := Transform(test.value, test.transformers)
		if !reflect.DeepEqual(m, test.m) {
			t.Errorf("Expected %v, got %v", test.m, m)
		}
	}
}

func TestTransformPanicsOnInvalidValue(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Errorf("Expected panic didn't happen")
		}
	}()
	Transform("foo", map[string]Transformer{})
}
