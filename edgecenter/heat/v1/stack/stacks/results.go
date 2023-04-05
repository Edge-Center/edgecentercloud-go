package stacks

import (
	"time"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/pagination"
)

type commonResult struct {
	edgecloud.Result
}

// CreateResult represents the result of a Create operation.
type CreateResult struct {
	edgecloud.Result
}

// UpdateResult represents the result of a Update operation.
type UpdateResult struct {
	edgecloud.ErrResult
}

// DeleteResult represents the result of a Delete operation.
type DeleteResult struct {
	edgecloud.ErrResult
}

// GetResult represents the result of a get operation. Call its Extract method to interpret it as a Heat stack.
type GetResult struct {
	commonResult
}

// Extract is a function that accepts a result and extracts a heat stack.
func (r commonResult) Extract() (*Stack, error) {
	var s Stack
	err := r.ExtractInto(&s)
	return &s, err
}

// Extract is a function that accepts a result and extracts a heat stack.
func (r CreateResult) Extract() (*CreatedStack, error) {
	var s struct {
		CreatedStack *CreatedStack `json:"stack"`
	}
	err := r.ExtractInto(&s)
	return s.CreatedStack, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// StackList struct.
type StackList struct {
	CreationTime       time.Time        `json:"creation_time"`
	DeletionTime       *time.Time       `json:"deletion_time"`
	UpdatedTime        *time.Time       `json:"updated_time"`
	Description        string           `json:"description"`
	ID                 string           `json:"id"`
	Links              []edgecloud.Link `json:"links"`
	Parent             *string          `json:"parent"`
	StackName          string           `json:"stack_name"`
	StackOwner         *string          `json:"stack_owner"`
	StackStatus        string           `json:"stack_status"`
	StackStatusReason  *string          `json:"stack_status_reason"`
	StackUserProjectID string           `json:"stack_user_project_id"`
	Tags               []string         `json:"tags"`
}

// Stack struct.
type Stack struct {
	*StackList
	Capabilities        []string                 `json:"capabilities"`
	DisableRollback     bool                     `json:"disable_rollback"`
	NotificationTopics  []string                 `json:"notification_topics"`
	TemplateDescription *string                  `json:"template_description"`
	TimeoutMinutes      int                      `json:"timeout_mins"`
	Outputs             []map[string]interface{} `json:"outputs"`
	Parameters          map[string]interface{}   `json:"parameters"`
}

// CreatedStack represents the object extracted from a Create operation.
type CreatedStack struct {
	ID    string           `json:"id"`
	Links []edgecloud.Link `json:"links"`
}

// StackPage is the page returned by a pager when traversing over a
// collection of loadbalancers.
type StackPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of loadbalancers has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r StackPage) NextPageURL() (string, error) {
	var s struct {
		Links []edgecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return edgecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a StackPage struct is empty.
func (r StackPage) IsEmpty() (bool, error) {
	is, err := ExtractStacks(r)
	return len(is) == 0, err
}

// ExtractStack accepts a Page struct, specifically a StackPage struct,
// and extracts the elements into a slice of Stack structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractStacks(r pagination.Page) ([]StackList, error) {
	var s []StackList
	err := ExtractStacksInto(r, &s)
	return s, err
}

func ExtractStacksInto(r pagination.Page, v interface{}) error {
	return r.(StackPage).Result.ExtractIntoSlicePtr(v, "results")
}
