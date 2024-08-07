package actions

type ChwdirRequest struct {
	*BaseFileRequest
	path string
}

func NewChwdirRequest(id string, path string) *ChwdirRequest {
	return &ChwdirRequest{
		BaseFileRequest: NewBaseFileRequest(id),
		path:            path}
}

//"""Request to change the current working directory."""
//"The path to change the current working directory to. "
//"Can be absolute, relative to the current working directory, or use '..' to navigate up the directory tree.",

type ChwdirResponse struct {
	*BaseFileResponse
}

func NewChwdirResponse(cwd string) *ChwdirResponse {
	return &ChwdirResponse{
		NewBaseFileResponse(cwd),
	}
}

//"""Response to change the current working directory."""

type ChangeWorkingDirectory struct {
	*BaseFileAction
	displayName    string          // = "Change Working Directory"
	requestSchema  *ChwdirRequest  //= ChwdirRequest
	responseSchema *ChwdirResponse //= ChwdirResponse
}

func NewChangeWorkingDirectory() *ChangeWorkingDirectory {
	return &ChangeWorkingDirectory{
		NewBaseFileAction(),
		"Change Working Directory",
		NewChwdirRequest("", ""),
		NewChwdirResponse(""),
	}
}

//"""
//Changes the current working directory of the file manager to the specified path.
//The commands after this action will be executed in this new directory.
//Can result in:
//- PermissionError: If the user doesn't have permission to access the directory.
//- FileNotFoundError: If the directory or any parent directory does not exist.
//- RuntimeError: If the path cannot be resolved due to a loop or other issues.
//"""

func (cwd *ChangeWorkingDirectory) ExecuteOnFileManager(
	fileManager FileManager, requestData ChwdirRequest) *ChwdirResponse {
	err := fileManager.Chdir(requestData.path)
	ncr := NewChwdirResponse("")
	ncr.Error = err
	return ncr
}
