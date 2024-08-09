package base

type Model interface {
	GetItems() map[string]any
	GetModelFields() map[string]FieldInfo
}

type BaseModel struct {
	Items       map[string]any
	ModelFields map[string]FieldInfo
}

func (b *BaseModel) GetItems() map[string]any {
	return b.Items
}

func (b *BaseModel) GetModelFields() map[string]FieldInfo {
	return b.ModelFields
}

type FieldInfo struct {
	Annotation          any
	Default             any
	DefaultFactory      func([]any) any
	Alias               string
	AliasPriority       int
	ValidationAlias     string
	SerializationAlias  string
	Title               string
	FieldTitleGenerator func(string, FieldInfo) string
	Description         string
	Examples            []any
	Exclude             bool
	Discriminator       string
	Deprecated          string
	Frozen              bool
	ValidateDefault     bool
	Repr                bool
	Init                bool
	InitVar             bool
	KwOnly              bool
	MetaData            []any
	JsonSchemaExtra     map[string]bool
}
