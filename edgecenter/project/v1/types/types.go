package types

import (
	"encoding/json"
	"fmt"
)

type ProjectState string

const (
	ProjectStateActive   ProjectState = "ACTIVE"
	ProjectStateDeleted  ProjectState = "DELETED"
	ProjectStateDeleting ProjectState = "DELETING"
)

func (rs ProjectState) IsValid() error {
	switch rs {
	case ProjectStateActive,
		ProjectStateDeleted,
		ProjectStateDeleting:
		return nil
	}
	return fmt.Errorf("invalid ProjectState type: %v", rs)
}

func (rs ProjectState) ValidOrNil() (*ProjectState, error) {
	if rs.String() == "" {
		return nil, nil
	}
	err := rs.IsValid()
	if err != nil {
		return &rs, err
	}
	return &rs, nil
}

func (rs ProjectState) String() string {
	return string(rs)
}

func (rs ProjectState) List() []ProjectState {
	return []ProjectState{
		ProjectStateActive,
		ProjectStateDeleted,
		ProjectStateDeleting,
	}
}

func (rs ProjectState) StringList() []string {
	lst := rs.List()
	strings := make([]string, 0, len(lst))
	for _, x := range lst {
		strings = append(strings, x.String())
	}
	return strings
}

// UnmarshalJSON - implements Unmarshaler interface for ProjectState.
func (rs *ProjectState) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := ProjectState(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*rs = v
	return nil
}

// MarshalJSON - implements Marshaler interface for ProjectState.
func (rs *ProjectState) MarshalJSON() ([]byte, error) {
	return json.Marshal(rs.String())
}
