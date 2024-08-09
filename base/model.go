package base

import (
	"fmt"
	"github.com/goccy/go-json"
	jsi "github.com/invopop/jsonschema"
	jsl "github.com/kaptinlin/jsonschema"
	"log"
)

type JsonValue interface {
	Get(string) JsonValue
}

type JsonInt int
type JsonFloat float64
type JsonString string
type JsonBool bool
type JsonArray []JsonValue
type JsonDict map[string]JsonValue
type JsonSchema = JsonDict
type JsonDictError struct{}
type JsonType interface {
	int | float64 | string | bool
}

func (e *JsonDictError) Error() string {
	return fmt.Sprintf("Wrong!!!,because is not a JsonMap")
}

type Model interface {
	ModelValidateJson(schemaJSON []byte) Response
	ModelJsonSchema() JsonSchema
	GetModelFields() JsonDict
	Items() JsonDict
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
	JsonSchemaExtra     JsonDict
	Frozen              bool
	ValidateDefault     bool
	Repr                bool
	Init                bool
	InitVar             bool
	KwOnly              bool
	MetaData            []any
}

type BaseModel struct {
	ModelFields map[string]FieldInfo
}

func (b *BaseModel) ModelJsonSchema() JsonSchema {
	reflector := jsi.Reflector{}
	schema := reflector.Reflect(b.ModelFields)
	schemaJSON, _ := json.Marshal(schema)

	var schemaMap JsonSchema
	_ = json.Unmarshal(schemaJSON, &schemaMap)
	return schemaMap
}

func (b *BaseModel) ModelValidateJson(schemaJSON []byte) *BaseModel {
	reflector := jsi.Reflector{}
	jschema := reflector.Reflect(b.ModelFields)
	schemaJson, _ := json.Marshal(jschema)
	compiler := jsl.NewCompiler()
	schema, err := compiler.Compile(schemaJson)
	if err != nil {
		panic(fmt.Errorf("failed to compile schema: %v", err))
	}

	result := schema.Validate(schemaJSON)
	if result.IsValid() {
		return b
	} else {
		log.Println("The schema is not valid. See errors:")
		details, _ := json.MarshalIndent(result.ToList(), "", "  ")
		log.Println(string(details))
		panic(fmt.Errorf("The schema is not valid. "))
	}
}

func (v JsonInt) Get(key string) JsonValue {
	return JsonDictError{}
}

func (v JsonFloat) Get(key string) JsonValue {
	return JsonDictError{}
}

func (v JsonBool) Get(key string) JsonValue {
	return JsonDictError{}
}

func (v JsonString) Get(key string) JsonValue {
	return JsonDictError{}
}

func (v JsonArray) Get(key string) JsonValue {
	return JsonDictError{}
}

func (v JsonDictError) Get(key string) JsonValue {
	return JsonDictError{}
}

func (v JsonDict) Get(key string) JsonValue {
	value, ok := v[key]
	if ok {
		return value
	} else {
		return nil
	}
}

func ToString(value JsonValue) string {
	return string(value.(JsonString))
}
