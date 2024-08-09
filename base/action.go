package base

import (
	"encoding/base64"
	"fmt"
	"github.com/goccy/go-json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"unicode/utf8"
)

type Request interface {
	Model
}

type Response interface {
	Model
}

type FileModel struct {
	*BaseModel
	Name    string `json:"name"`
	Content []byte `json:"content"`
}

var BaseFileModel = FileModel{
	Name:    "",
	Content: nil,
}

type Action interface {
	ToolName() string
	SetToolName(string)
	ActionName() string
	DisplayName() string
	SetDisplayName(string)
	Tags() []string
	SetTags([]string)
	RequestSchema() Request
	SetRequestSchema(request Request)
	ResponseSchema() Response
	SetResponseSchema(response Response)
	Execute(Request, map[string]any) (map[string]any, Response)
	RequiredScopes() []string
	GetToolMergedActionName() string
	checkFileUploadable(param string) bool
	ExecuteAction(Request, map[string]any,
	) (map[string]any, Response)
}

type BaseAction struct {
	toolName         string
	historyMaintains bool
	displayName      string   // Add an internal variable to hold the display name
	requestSchema    Request  // Placeholder for request schema
	responseSchema   Response // Placeholder for response schema
	tags             []string // Placeholder for tags

	//# For workspace
	RunOnShell bool
	Requires   []string // List of python dependencies
	Module     string   // File where this tool is defined
}

// ToolName @property
func (b *BaseAction) ToolName() string {
	return b.toolName
}

// SetToolName ToolName 的 setter 方法
func (b *BaseAction) SetToolName(value string) {
	b.toolName = value
}

// ActionName 的 getter 方法
func (b *BaseAction) ActionName() string {
	return reflect.TypeOf(b).Elem().Name()
}

// DisplayName 的 getter 方法
func (b *BaseAction) DisplayName() string {
	return b.displayName
}

// SetDisplayName DisplayName 的 setter 方法
func (b *BaseAction) SetDisplayName(value string) {
	b.displayName = value
}

// Tags 的 getter 方法
func (b *BaseAction) Tags() []string {
	return b.tags
}

// SetTags Tags 的 setter 方法
func (b *BaseAction) SetTags(value []string) {
	b.tags = value
}

// RequestSchema 的 getter 方法
func (b *BaseAction) RequestSchema() Request {
	return b.requestSchema
}

// SetRequestSchema RequestSchema 的 setter 方法
func (b *BaseAction) SetRequestSchema(value Request) {
	b.requestSchema = value
}

// ResponseSchema 的 getter 方法
func (b *BaseAction) ResponseSchema() Response {
	return b.responseSchema
}

// SetResponseSchema ResponseSchema 的 setter 方法
func (b *BaseAction) SetResponseSchema(value Response) {
	b.responseSchema = value
}

func (b *BaseAction) Execute(requestData Request, authorisationData map[string]any) (map[string]any, Response) {
	return nil, nil
}

// RequiredScopes @property
func (b *BaseAction) RequiredScopes() []string {
	return []string{}
}

func (b *BaseAction) GetToolMergedActionName() string {
	return b.toolName + b.ActionName()
}

//func (b *BaseAction)GetActionSchema() {
//
//	//requestSchemaJson := b.requestSchema.ModelJsonSchema(byAlias:False)
//	//modifiedProperties := requestSchemaJson.get("properties", {})
//	for _, details := range modifiedProperties.items() {
//		if details.get("file_readable", False) {
//			details["oneOf"] =
//			{
//				{
//					"type": details.get("type"),
//					"description": details.get("description", ""),
//				}, {
//				"type": "string",
//					"format": "file-path",
//					"description": f
//				"File path to {details.get('description', '')}",
//			},
//			}
//			del
//			details["type"] //  # Remove original type to avoid conflict in oneOf
//		}
//	}
//	requestSchemaJson["properties"] = modifiedProperties
//	actionSchema =
//	{
//		"appKey": b.toolName,
//		"appName": b.toolName,
//		"logo": "empty",
//		"appId": generateHashedAppId(b.toolName),
//		"name": b.GetToolMergedActionName(),
//		"display_name": b.displayName,
//		"tags": b.tags,
//		"enabled": True,
//		"description": (
//		b.__class__.__doc__
//		if self.__class__.__doc__
//		else
//		b.ActionName()
//		),
//		"parameters": json.Loads(
//		jsonref.dumps(
//			jsonref.replace_refs(
//				b.requestSchema.ModelJsonSchema(byAlias: False),
//	lazyLoad:
//		False,
//	)
//	)
//	),
//		"response": json.Loads(
//		jsonref.dumps(
//			jsonref.replaceRefs(
//				b.responseSchema.ModelJsonSchema(),
//				lazyLoad: False,
//	)
//	)
//	),
//	}
//	return actionSchema
//}

func (b *BaseAction) checkFileUploadable(param string) bool {
	/*
		return b.requestSchema.ModelJsonSchema().Get("properties").Get(param).Get("allOf").(JsonArray)[0].Get("properties") ==
			BaseFileModel.ModelJsonSchema().Get("properties") ||
			b.requestSchema.ModelJsonSchema().Get("properties").Get(param).Get("properties") ==
				BaseFileModel.ModelJsonSchema().Get("properties")
	*/
	return true
}

func (b *BaseAction) ExecuteAction(
	requestData Request,
	metaData map[string]any,
) (map[string]any, Response) {
	modifiedRequestData := make(map[string]any)
	for param, value := range requestData.Items() { // # type: ignore
		if _, ok := b.requestSchema.GetModelFields()[param]; !ok {
			panic(fmt.Errorf("invalid param %v for action %v", param, strings.ToUpper(b.GetToolMergedActionName())))
		}
		annotations := b.requestSchema.GetModelFields()[param].Get("JsonSchemaExtra")
		fileReadable := bool(annotations != nil && annotations.Get("file_readable").(JsonBool))
		if ok, _ := isFile(string(value.(JsonString))); fileReadable && isInstance(value, "JsonString") && ok {
			fileContent, err := os.ReadFile(ToString(value))
			if err != nil {
			}
			if utf8.Valid(fileContent) { //  # Try decoding	as UTF - 8 to check	if it's normal text
				modifiedRequestData[param] = string(fileContent)
			} else { //# If decoding fails, treat as binary and encode in base64
				modifiedRequestData[param] = base64.StdEncoding.EncodeToString(fileContent)
			}
		} else if b.checkFileUploadable(param) && isInstance(value, "string") && ok {
			// For uploadable files, we	also need to send the filename
			fileContent, err := os.ReadFile(ToString(value))
			if err != nil {
			}
			modifiedRequestData[param] = map[string]string{
				"name":    filepath.Base(ToString(value)),
				"content": base64.StdEncoding.EncodeToString(fileContent),
			}
		} else {
			modifiedRequestData[param] = value
		}
	}
	instance, _ := json.Marshal(modifiedRequestData)
	res1, res2 := b.Execute(b.requestSchema.ModelValidateJson(instance), metaData)

	if r := recover(); r != nil {
		return map[string]any{
			"status":  "failure",
			"details": "Error executing action with error: " + r.(error).Error(),
		}, nil
	} else {
		return res1, res2
	}
}

func isInstance(value any, s string) bool {
	return reflect.TypeOf(value).Name() == s
}
