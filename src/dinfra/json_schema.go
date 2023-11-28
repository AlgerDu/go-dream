package dinfra

type (
	JsonPropertyType string // json schema 支持的属性类型

	// json schema 属性描述
	JsonProperty struct {
		Type        JsonPropertyType `json:"type"`
		Title       *string          `json:"title"`       // 最好是简短的
		Description *string          `json:"description"` // 说明，更长更加详细

		MaxLength *int64  `json:"maxLength"` // string 最大长度，非负
		MinLength *int64  `json:"minLength"` // string 最小长度，非负
		Pattern   *string `json:"pattern"`   // string 必须符合对应的正则表达式

		MultipleOf *int64 `json:"multipleOf"` // number integer 给定数字的倍数
		Minimum    *int64 `json:"minimum"`    // number integer 最小值 >=
		Maximum    *int64 `json:"maximum"`    // number integer 最大值 <=

		Properties map[string]JsonProperty `json:"properties"`
		Required   []string                `json:"required"` // object 必须属性

		Items       *JsonPropertyType `json:"items"`       // array 列表项的说明
		MinItems    *int64            `json:"minItems"`    // array 数组最小长度
		MaxItems    *int64            `json:"maxItems"`    // array 数组最大长度
		UniqueItems *bool             `json:"uniqueItems"` // array 数组每个元素唯一

		Enum    []any `json:"enum"`
		Default any   `json:"default"` // 该值不用于在验证过程中填充缺失值。文档生成器或表单生成器等非验证工具可能会使用此值提示用户如何使用该值。
		Const   any   `json:"const"`   // 常量，固定值
	}

	// json schema
	JsonSchema struct {
		JsonProperty

		Schema *string `json:"$schema"` // 使用的 schema 版本
		ID     *string `json:"$id"`     //
	}
)

var (
	JPT_Number  JsonPropertyType = "number"  // 整数类型
	JPT_Integer JsonPropertyType = "integer" // 数值类型
	JPT_Boolean JsonPropertyType = "boolean" // 布尔
	JPT_String  JsonPropertyType = "string"  // 字符串
	JPT_Object  JsonPropertyType = "object"  // 对象
	JPT_Array   JsonPropertyType = "array"   // 数组
	JPT_Null    JsonPropertyType = "null"    // 空

	JPT_Datetime JsonPropertyType = "datetime" // 时间(扩展定义)
)
