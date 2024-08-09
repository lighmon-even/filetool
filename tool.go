package filetool

import (
	"fmt"
	"github.com/lighmon-even/filetool/actions"
	. "github.com/lighmon-even/filetool/base"
)

var (
	ChangeWorkingDirectory = actions.NewChangeWorkingDirectory()
	CreateFile             = actions.NewCreateFile()
	EditFile               = actions.NewEditFile()
)

type FileTool struct {
	BaseTool
}

//File I/O tool.

func (ft *FileTool) Actions() ([]Action, error) {
	//Return the list of actions.
	return []Action{
		//OpenFile,
		EditFile,
		CreateFile,
		//Scroll,
		//ListFiles,
		//SearchWord,
		//FindFile,
		//Write,
		ChangeWorkingDirectory,
		//GitClone,
		//GitRepoTree,
		//GitPatch,
	}, nil
}

func (ft *FileTool) GetAction(name string) (Action, error) {
	acts, err := ft.Actions()
	if err != nil {
		return nil, err
	}
	for _, action := range acts {
		instance := action
		if instance.GetToolMergedActionName() != name {
			continue
		}
		return instance, nil
	}
	return nil, fmt.Errorf("no action found with name %v", name)
}

func (ft *FileTool) Triggers() ([]any, error) {
	//Return the list of triggers.
	return []any{}, nil
}
