package base

import (
	"encoding/base64"
	"fmt"
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

func (b *BaseAction) checkFileUploadable(param string) bool {
	return b.requestSchema.GetModelFields() == nil && BaseFileModel.GetModelFields() == nil
}

func (b *BaseAction) ExecuteAction(
	requestData Request,
	metaData map[string]any,
) (map[string]any, Response) {
	var Error error
	modifiedRequestData := make(map[string]any)
	for param, value := range requestData.GetItems() { // # type: ignore
		modelField, ok := b.requestSchema.GetModelFields()[param]
		if !ok {
			panic(fmt.Errorf("invalid param %v for action %v", param, strings.ToUpper(b.GetToolMergedActionName())))
		}
		annotations := modelField.JsonSchemaExtra
		fileReadable := annotations != nil && annotations["file_readable"]
		svalue, ok := value.(string)
		if isFile(svalue) && fileReadable && ok {
			fileContent, err := os.ReadFile(svalue)
			if err != nil {
				Error = err
				break
			}
			if utf8.Valid(fileContent) { //  # Try decoding	as UTF - 8 to check	if it's normal text
				modifiedRequestData[param] = string(fileContent)
			} else { //# If decoding fails, treat as binary and encode in base64
				modifiedRequestData[param] = base64.StdEncoding.EncodeToString(fileContent)
			}
		} else if b.checkFileUploadable(param) && isFile(svalue) && ok {
			// For uploadable files, we	also need to send the filename
			fileContent, err := os.ReadFile(svalue)
			if err != nil {
				Error = err
				break
			}
			modifiedRequestData[param] = map[string]any{
				"name":    filepath.Base(svalue),
				"content": base64.StdEncoding.EncodeToString(fileContent),
			}
		} else {
			modifiedRequestData[param] = value
		}
	}
	if Error != nil {
		return b.Execute(b.requestSchema, metaData)
	} else {
		return map[string]any{
			"status":  "failure",
			"details": "Error executing action with error: " + Error.Error(),
		}, nil
	}
}
