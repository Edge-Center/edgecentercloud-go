package servergroups

import (
	"encoding/json"
	"fmt"
)

type ServerGroupPolicy string

const (
	AffinityPolicy     ServerGroupPolicy = "affinity"
	AntiAffinityPolicy ServerGroupPolicy = "anti-affinity"
)

func (s ServerGroupPolicy) String() string {
	return string(s)
}

func (s ServerGroupPolicy) List() []ServerGroupPolicy {
	return []ServerGroupPolicy{AffinityPolicy, AntiAffinityPolicy}
}

func (s ServerGroupPolicy) StringList() []string {
	lst := s.List()
	strings := make([]string, 0, len(lst))
	for _, x := range lst {
		strings = append(strings, x.String())
	}
	return strings
}

func (s ServerGroupPolicy) IsValid() error {
	switch s {
	case AffinityPolicy, AntiAffinityPolicy:
		return nil
	}
	return fmt.Errorf("invalid ServerGroupPolicy type: %v", s)
}

func (s ServerGroupPolicy) ValidOrNil() (*ServerGroupPolicy, error) {
	if s.String() == "" {
		return nil, nil
	}
	err := s.IsValid()
	if err != nil {
		return &s, err
	}
	return &s, nil
}

// UnmarshalJSON - implements Unmarshaler interface for ServerGroupPolicy.
func (s *ServerGroupPolicy) UnmarshalJSON(data []byte) error {
	var sg string
	if err := json.Unmarshal(data, &sg); err != nil {
		return err
	}
	v := ServerGroupPolicy(sg)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*s = v
	return nil
}

// MarshalJSON - implements Marshaler interface for ServerGroupPolicy.
func (s *ServerGroupPolicy) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
