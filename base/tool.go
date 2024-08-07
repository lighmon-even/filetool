package base

import (
	"errors"
	"fmt"
)

// Tool Abstraction for local tools
type Tool interface {
	GetAction(name string) Action
	Actions() ([]Action, error)
	Triggers() ([]any, error)
}
type BaseTool struct {
	Name string
}

// GetAction Get action object.
func (bt *BaseTool) GetAction(name string) (Action, error) {
	actions, err := bt.Actions()
	if err != nil {
		return nil, err
	}
	for _, action := range actions {
		instance := action
		if instance.GetToolMergedActionName() != name {
			continue
		}
		return instance, nil
	}
	return nil, fmt.Errorf("no action found with name %v", name)
}

func (bt *BaseTool) Actions() ([]Action, error) {
	return []Action{}, errors.New("this method should be overridden by subclasses")
}

func (bt *BaseTool) Triggers() ([]any, error) {
	return []any{}, errors.New("this method should be overridden by subclasses")
}
