package actions

import "filetool/base"

type FileManager = base.FileManager

type data struct {
	FileManagers map[string]FileManager
}

type BaseFileRequest struct {
	//"ID of the file manager where the file will be opened, if not "
	//"provided the recent file manager will be used to execute the action"
	FileManagerId string
}

func NewBaseFileRequest(fileManagerId string) *BaseFileRequest {
	return &BaseFileRequest{FileManagerId: fileManagerId}
}

type BaseFileResponse struct {
	//"Error message if the action failed"
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
	ExecuteOnFileManager(fileManager FileManager, requestData BaseFileRequest) BaseFileResponse
}

type BaseFileAction struct{}

func NewBaseFileAction() *BaseFileAction {
	return &BaseFileAction{}
}

func (s *BaseFileAction) ExecuteOnFileManager(fileManager FileManager, requestData BaseFileRequest) BaseFileResponse {
	return BaseFileResponse{}
}
func (s *BaseFileAction) Execute(requestData BaseFileRequest, authorisationData map[string]data) BaseFileResponse {
	fileManagers := authorisationData["workspace"].FileManagers // type: ignore
	fileManager := fileManagers[requestData.FileManagerId]
	resp := s.ExecuteOnFileManager(fileManager, requestData)
	resp.CurrentWorkingDirectory = fileManager.WorkingDir
	return resp
}
