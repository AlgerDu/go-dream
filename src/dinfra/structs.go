package dinfra

import (
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type (
	CanToMap interface {
		ToMap() map[string]any
	}
)

func MapToStruct[DstType any](from map[string]any) (*DstType, error) {
	to := new(DstType)
	err := mapstructure.Decode(from, to)
	return to, err
}

func StructToMap(from any) (map[string]any, error) {
	v := reflect.ValueOf(from)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ToMap only accepts struct or struct pointer; got %T", v)
	}

	return structToMap(v)
}

func fieldToAny(fieldV reflect.Value) (any, bool) {
	fieldVKind := fieldV.Kind()

	if fieldVKind == reflect.Func ||
		(fieldVKind == reflect.Slice && fieldV.IsNil()) ||
		(fieldVKind == reflect.Map && fieldV.IsNil()) ||
		(fieldVKind == reflect.Pointer && fieldV.IsNil()) ||
		(fieldVKind == reflect.Interface && fieldV.IsNil()) {
		return nil, true
	}

	if fieldVKind == reflect.Pointer {
		fieldV = fieldV.Elem()
		fieldVKind = fieldV.Kind()
	}

	switch fieldVKind {
	case reflect.Struct:
		tmp, _ := structToMap(fieldV)
		return tmp, false
	case reflect.Map:
		tmp, _ := mapToMap(fieldV)
		return tmp, false
	case reflect.Slice:
		tmp, _ := sliceToAny(fieldV)
		return tmp, false
	default:
		return fieldV.Interface(), false
	}
}

func structToMap(v reflect.Value) (map[string]any, error) {
	vt := v.Type()
	out := map[string]any{}
	for i := 0; i < v.NumField(); i++ {
		tmp, skip := fieldToAny(v.Field(i))
		if skip {
			continue
		}

		field := vt.Field(i)
		fieldName := Strings_FirstToLower(field.Name)

		if field.Anonymous {
			for k, v := range tmp.(map[string]any) {
				if _, exist := out[k]; !exist {
					out[k] = v
				}
			}
		} else {
			out[fieldName] = tmp
		}
	}

	return out, nil
}

func mapToMap(v reflect.Value) (map[string]any, error) {
	out := map[string]any{}

	keys := v.MapKeys()
	for _, key := range keys {
		keyValue := v.MapIndex(key)
		tmp, skip := fieldToAny(keyValue)
		if skip {
			continue
		}
		out[key.String()] = tmp
	}
	return out, nil
}

func sliceToAny(v reflect.Value) (any, error) {
	out := []any{}

	for i := 0; i < v.Len(); i++ {
		tmp, skip := fieldToAny(v.Index(i))
		if skip {
			continue
		}
		out = append(out, tmp)
	}
	return out, nil
}
