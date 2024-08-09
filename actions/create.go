package actions

import (
	"errors"
	"strings"
)

type CreateFileRequest struct {
	*BaseFileRequest
	//"""Request to create a file."""
	FilePath string
	//"""File path to create in the editor.
	//If file already exists, it will be overwritten"""
}

// @field_validator("file_path")
// @classmethod

func (cfr *CreateFileRequest) ValidateFilePath(v string) (string, error) {
	if strings.TrimSpace(v) == "" {
		return "", errors.New("file name cannot be empty or just whitespace")
	}
	if v == "." || v == ".." {
		return "", errors.New(`file name cannot be "." or ".."`)
	}
	return v, nil
}

type CreateFileResponse struct {
	*BaseFileResponse
	//"""Response to create a file."""
	file string
	//"Path the created file."
	success bool
	//"Whether the file was created successfully"
}

func NewCreateFileResponse(file string, success bool) *CreateFileResponse {
	return &CreateFileResponse{
		NewBaseFileResponse(""),
		file,
		success,
	}
}

type CreateFile struct {
	*BaseFileRequest
	//"""
	//Creates a new file within a shell session.
	//Example:
	//- To create a file, provide the path of the new file. If the path you provide
	//is relative, it will be created relative to the current working directory.
	//- The response will indicate whether the file was created successfully and list any errors.
	//Raises:
	//- ValueError: If the file path is not a string or if the file path is empty.
	//- FileExistsError: If the file already exists.
	//- PermissionError: If the user does not have permission to create the file.
	//- FileNotFoundError: If the directory does not exist.
	//- OSError: If an OS-specific error occurs.
	//"""

	displayName    string             //= "Create a new file"
	requestSchema  *CreateFileRequest //= CreateFileRequest
	responseSchema *ChwdirResponse    //= CreateFileResponse
}

func (cf *CreateFile) ExecuteOnFileManager(
	fileManager FileManager, requestData CreateFileRequest,
) (cfr *CreateFileResponse) {
	cfr = NewCreateFileResponse("", false)
	newfile, err := fileManager.Create(requestData.FilePath)
	if err == nil {
		newfilepath := newfile.Path
		cfr.file = newfilepath
		cfr.success = true
		return cfr
	} else {
		cfr.Error = err
		return cfr
	}
}
