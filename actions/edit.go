package actions

import (
	"errors"
	"filetool/base"
	"fmt"
)

type EditFileRequest struct {
	*BaseFileRequest
	//"""Request to edit a file."""
	FilePath string
	//"The path to the file that will be edited. If not provided, "
	//"THE CURRENTLY OPEN FILE will be edited. If provided, the "
	//"file at the provided path will be OPENED and edited, changing "
	//"the opened file."
	Text string
	//"The text that will replace the specified line range in the file."
	StartLine int
	//"The line number at which the file edit will start (REQUIRED). Inclusive - the start line will be included in the edit.",
	EndLine int
	//"The line number at which the file edit will end (REQUIRED). Inclusive - the end line will be included in the edit.",
}

type EditFileResponse struct {
	*BaseFileResponse
	//"""Response to edit a file."""
	OldText string
	//"The updated changes. If the file was not edited, the original file "
	//"will be returned."
	Error error
	//"Error message if any",
	UpdatedText string
	//"The updated text. If the file was not edited, this will be empty.",
}

func NewEditFileResponse() *EditFileResponse {
	return &EditFileResponse{
		NewBaseFileResponse(""),
		"",
		nil,
		"",
	}
}

type EditFile struct {
	*BaseFileAction
	//"""
	//Use this tools to edit a file on specific line numbers.
	//
	//Please note that THE EDIT COMMAND REQUIRES PROPER INDENTATION.
	//
	//Python files will be checked for syntax errors after the edit.
	//If you'd like to add the line '        print(x)' you must fully write
	//that out, with all those spaces before the code!
	//
	//If the system detects a syntax error, the edit will not be executed.
	//Simply try to edit the file again, but make sure to read the error message
	//and modify the edit command you issue accordingly. Issuing the same command
	//a second time will just lead to the same error message again.
	//
	//If a lint error occurs, the edit will not be applied. Review the error message,
	//adjust your edit accordingly, and try again.
	//Raises:
	//- FileNotFoundError: If the file does not exist.
	//- PermissionError: If the user does not have permission to edit the file.
	//- OSError: If an OS-specific error occurs.
	//Note:
	//This action edits a specific part of the file, if you want to rewrite the
	//complete file, use `write` tool instead.
	//"""
	displayName    string           // = "Edit a file"
	requestSchema  EditFileRequest  // = EditFileRequest
	responseSchema EditFileResponse //= EditFileResponse
}

func (ef *EditFile) ExecuteOnFileManager(
	fileManager FileManager,
	requestData EditFileRequest,
) (efr *EditFileResponse) {
	efr = NewEditFileResponse()
	var file *base.File
	var err error
	if requestData.FilePath == "" {
		file = fileManager.Recent
	} else {
		file, err = fileManager.Open(requestData.FilePath)
	}
	if file == nil {
		err = fmt.Errorf("file not found: %v", requestData.FilePath)
	}
	response := file.WriteAndRunLint(
		requestData.Text,
		requestData.StartLine,
		requestData.EndLine,
	)
	if response.Error != nil && len(response.Error.Error()) > 0 {
		efr.Error = errors.New("No Update, found error: " + response.Error.Error())
		return
	}
	if err == nil {
		efr.OldText = response.ReplacedText
		efr.UpdatedText = response.ReplacedWith
		return
	} else {
		efr.Error = err
		return
	}
}
