package dinfra

import (
	"encoding/json"
	"testing"
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
