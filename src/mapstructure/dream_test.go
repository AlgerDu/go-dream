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
