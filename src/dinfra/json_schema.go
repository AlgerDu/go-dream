package dinfra

type (
	JsonPropertyType string

	JsonProperty struct {
		Type        JsonPropertyType        `json:"type"`
		Title       string                  `json:"title"`
		Description string                  `json:"description"`
		Properties  map[string]JsonProperty `json:"properties"`
		Required    bool                    `json:"required"`
		Minimum     int64                   `json:"minimum"`
		Maximum     int64                   `json:"maximum"`
		MaxLength   int64                   `json:"maxLength"`
		MinLength   int64                   `json:"minLength"`
		Pattern     string                  `json:"pattern"`
		Enum        []any                   `json:"enum"`
	}

	JsonSchema struct {
		JsonProperty
	}
)

var (
	JPT_Number  JsonPropertyType = "number"
	JPT_Integer JsonPropertyType = "integer"
	JPT_Boolean JsonPropertyType = "boolean"
	JPT_String  JsonPropertyType = "string"
	JPT_Object  JsonPropertyType = "object"
	JPT_Array   JsonPropertyType = "array"
	JPT_Date    JsonPropertyType = "date"
)
