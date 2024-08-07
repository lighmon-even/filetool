package actions

import "github.com/lighmon-even/filetool/base"

type FileManager = base.FileManager

type data struct {
	FileManagers map[string]FileManager
}

type FileRequest interface {
	base.Request
	GetFileManagerID() string
}

type BaseFileRequest struct {
	//"ID of the file manager where the file will be opened, if not "
	//"provided the recent file manager will be used to execute the action"
	*base.BaseModel
	FileManagerId string
}

func (bfr *BaseFileRequest) GetFileManagerId() string {
	return bfr.FileManagerId
}

type FileResponse interface {
	base.Response
}

func NewBaseFileRequest(fileManagerId string) *BaseFileRequest {
	return &BaseFileRequest{FileManagerId: fileManagerId}
}

type BaseFileResponse struct {
	//"Error message if the action failed"
	*base.BaseModel
	Error error
	//"Current working directory of the file manager."
	CurrentWorkingDirectory string
}

func NewBaseFileResponse(currentWorkingDirectory string) *BaseFileResponse {
	return &BaseFileResponse{
		Error:                   nil,
		CurrentWorkingDirectory: currentWorkingDirectory,
	}
}

type FileAction interface {
	base.Action
	ExecuteOnFileManager(fileManager FileManager, requestData BaseFileRequest) BaseFileResponse
}

type BaseFileAction struct {
	*base.BaseAction
}

func NewBaseFileAction() *BaseFileAction {
	return &BaseFileAction{}
}

func (s *BaseFileAction) ExecuteOnFileManager(fileManager FileManager, requestData FileRequest) (map[string]any, FileResponse) {
	return nil, NewBaseFileResponse("")
}

func (s *BaseFileAction) Execute(requestData base.Request, authorisationData map[string]any) (map[string]any, base.Response) {
	//workPlace := authorisationData["workspace"]
	//fileManagers := workPlace.FileManagers // type: ignore
	//fileManager := fileManagers[requestData.GetFileManagerId()]
	//resp := s.ExecuteOnFileManager(fileManager, requestData)
	//resp.CurrentWorkingDirectory = fileManager.WorkingDir
	return nil, nil
}
