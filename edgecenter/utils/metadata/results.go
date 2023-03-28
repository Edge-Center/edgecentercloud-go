package metadata

import (
	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/pagination"
)

type commonResult struct {
	edgecloud.Result
}

type Metadata struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	ReadOnly bool   `json:"read_only"`
}

type ResourceMetadataPage struct {
	pagination.LinkedPageBase
}

// ResourceMetadataResult represents the result of a get operation.
type ResourceMetadataResult struct {
	commonResult
}

func ExtractMetadataInto(r pagination.Page, v interface{}) error {
	return r.(ResourceMetadataPage).Result.ExtractIntoSlicePtr(v, "results")
}

// ExtractMetadata accepts a Page struct, specifically a ResourceMetadataPage struct,
// and extracts the elements into a slice of securitygroups metadata structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractMetadata(r pagination.Page) ([]Metadata, error) {
	var s []Metadata
	err := ExtractMetadataInto(r, &s)
	return s, err
}

// ActionResultMetadata represents the result of a creation, delete or update operation(no content).
type ActionResultMetadata struct {
	edgecloud.ErrResult
}

func (r ResourceMetadataResult) Extract() (*Metadata, error) {
	var s Metadata
	err := r.ExtractInto(&s)
	return &s, err
}

// IsEmpty checks whether a ResourceMetadataPage struct is empty.
func (r ResourceMetadataPage) IsEmpty() (bool, error) {
	is, err := ExtractMetadata(r)
	return len(is) == 0, err
}
