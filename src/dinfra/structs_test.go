package dinfra

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
)

type (
	human struct {
		Name string
	}

	student struct {
		*human

		Class string
	}
)

func TestStruct_Zero(t *testing.T) {
	schema := &JsonSchema{
		JsonProperty: &JsonProperty{
			Name: "test",
			Properties: map[string]*JsonProperty{
				"task": {
					Name: "task",
				},
			},
		},
	}

	m, _ := StructToMap(schema)
	mJson, _ := json.Marshal(m)
	t.Log(string(mJson))
}

func TestStructToMap_Zero(t *testing.T) {

	student := &student{
		human: &human{
			Name: "lihua",
		},
		Class: "1",
	}

	data := &map[string]any{}

	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   data,
		Squash:   true,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		t.Error(err)
	}

	err = decoder.Decode(student)
	if err != nil {
		t.Error(err)
	}

	t.Log(data)
}

func TestMapToStruct_Zero(t *testing.T) {
	data := map[string]any{
		"name":  "lihua",
		"class": "1",
	}

	student := &student{
		human: &human{},
	}

	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   student,
		Squash:   true,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		t.Error(err)
	}

	err = decoder.Decode(data)
	if err != nil {
		t.Error(err)
	}

	t.Log(student.Name)
	t.Log(student.Class)
}

func TestMapToStruct_Base(t *testing.T) {
	from := map[string]any{
		"name":  "lihua",
		"class": "1",
	}

	student, err := MapToStruct(from, &student{human: &human{}})
	if err != nil {
		t.Error(err)
	}

	t.Log(student.Name)
	t.Log(student.Class)
}
