package base

import (
	"reflect"
)

type Request interface {
}

type Response interface {
}

type FileModel struct {
	Name    string `json:"name"`
	Content string `json:"content"`
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

func (b *BaseAction) Execute(requestData Request, authorisationData map[string]interface{}) (map[string]interface{}, Response) {
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

//func (b *BaseAction)checkFileUploadable(param string) bool {
//	return (
//		b.requestSchema.ModelJsonSchema()
//	.get("properties",
//	{
//	})
//	.get(param,
//	{
//	})
//	.get("allOf",[
//	{
//	}])[0]
//.get("properties", {})) == FileModel.ModelJsonSchema().get("properties")
//|| b.requestSchema.ModelJsonSchema()
//.get("properties", {})
//.get(param, {})
//.get("properties", {})==FileModel.ModelJsonSchema().get("properties")
//}
//func (b *BaseAction)ExecuteAction(
//	requestData Request,
//	metaData map[string]interface{},
//) (map[string]interface{}, Response){
//	modified_request_data: t.Dict[str, t.Union[str, t.Dict[str, str]]] = {}
//	for param, value in request_data.items():  # type: ignore
//	if param not in self.request_schema.model_fields:
//	raise ValueError(
//	f"Invalid param `{param}` for action `{self.get_tool_merged_action_name().upper()}`"
//)
//	annotations = self.request_schema.model_fields[param].json_schema_extra
//	file_readable = annotations is not None and annotations.get(  # type: ignore
//	"file_readable", False
//)
//
//	if file_readable and isinstance(value, str) and os.path.isfile(value):
//	with open(value, "rb") as file:
//	file_content = file.read()
//	file_content.decode(
//	"utf-8"
//)  # Try decoding as UTF-8 to check if it's normal text
//	modified_request_data[param] = file_content.decode("utf-8")
//	except UnicodeDecodeError:
//	# If decoding fails, treat as binary and encode in base64
//	modified_request_data[param] = base64.b64encode(
//	file_content
//).decode("utf-8")
//	elif (
//	self._check_file_uploadable(param=param)
//	and isinstance(value, str)
//	and os.path.isfile(value)
//):
//	# For uploadable files, we also need to send the  filename
//	with open(value, "rb") as file:
//	file_content = file.read()
//
//	modified_request_data[param] = {
//	"name": os.path.basename(value),
//	"content": base64.b64encode(file_content).decode("utf-8"),
//}
//	else:
//	modified_request_data[param] = value
//
//	return self.execute(
//	request_data=self.request_schema.model_validate_json(
//	json_data=json.dumps(
//	modified_request_data,
//)
//),
//	authorisation_data=metadata,
//)
//	except json.JSONDecodeError as e:
//	return {
//	"status": "failure",
//	"details": f"Could not parse response with error: {e}. Please contact the tool developer.",
//}
//	except Exception as e:
//	self.logger.error(
//	"Error while executing `%s` with parameters `%s`; Error: %s",
//	self.display_name,
//	request_data,
//	traceback.format_exc(),
//)
//	return {
//	"status": "failure",
//	"details": "Error executing action with error: " + str(e),
//}
//
//}
//}
