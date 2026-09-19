package table

import (
	"reflect"
	"strings"
)

// Field represents a single struct field's header name and value.
type Field struct {
	// Header is the name of the struct field.
	Header string
	// Value is the value of the struct field.
	Value any
}

// StructHeaders extracts the column header names for a generic struct type T.
func StructHeaders[T any]() []string {
	var zero T

	headers := make([]string, 0)
	for _, f := range StructFields(zero) {
		headers = append(headers, f.Header)
	}

	return headers
}

// StructFields inspects a struct or a pointer to a struct and extracts its exported fields.
// It checks for a `table` tag first, falling back to the struct field name.
func StructFields(s any) []Field {
	if s == nil {
		return make([]Field, 0)
	}

	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	if t.Kind() == reflect.Pointer {
		if v.IsNil() {
			return make([]Field, 0)
		}
		v = v.Elem()
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return make([]Field, 0)
	}

	var result []Field

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		header := field.Name
		if tag := field.Tag.Get("table"); tag != "" && tag != "-" {
			header = strings.Split(tag, ",")[0]
		}

		value := v.Field(i).Interface()

		result = append(result, Field{
			Header: header,
			Value:  value,
		})
	}

	return result
}
