package mapstructure

import (
	"reflect"
	"testing"
)

// 新添加的一些单元测试；

type (
	People struct {
		Name string
	}

	Student struct {
		*People
		Class string
	}
)

func TestDecode_EmbeddedPointerWithConfigSquash(t *testing.T) {
	input := &Student{
		People: &People{
			Name: "lihua",
		},
		Class: "1",
	}

	var result map[string]any
	err := Decode(input, &result, WithSquash)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := map[string]any{
		"Name":  "lihua",
		"Class": "1",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("bad: %#v", result)
	}
}

func TestDecode_CustomConvertName(t *testing.T) {
	input := &Student{
		People: &People{
			Name: "lihua",
		},
		Class: "1",
	}
	var result map[string]any
	err := Decode(input, &result, func(config *DecoderConfig) {
		config.ConvertName = func(fieldName string) string {
			switch fieldName {
			case "People":
				return "people"
			case "Name":
				return "name"
			case "Class":
				return "class"
			default:
				return fieldName
			}
		}
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := map[string]any{
		"people": map[string]any{
			"name": "lihua",
		},
		"class": "1",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("bad: %#v", result)
	}
}
