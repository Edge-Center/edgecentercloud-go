package ports

import (
	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/instance/v1/instances"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/reservedfixedip/v1/reservedfixedips"
)

type commonResult struct {
	edgecloud.Result
}

// Extract is a function that accepts a result and extracts a security group resource.
func (r commonResult) Extract() (*instances.Interface, error) {
	var s instances.Interface
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// UpdateResult represents the result of a enable or disable operation. Call its Extract
// method to interpret it as a Interface.
type UpdateResult struct {
	commonResult
}

type assignResult struct {
	edgecloud.Result
}

type AssignResult struct {
	assignResult
}

// Extract is a function that accepts a result and extracts a security group resource.
func (r assignResult) Extract() (*InstancePort, error) {
	var s InstancePort
	err := r.ExtractInto(&s)
	return &s, err
}

func (r assignResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

type InstancePort struct {
	NetworkID           string                                 `json:"network_id"`
	AllowedAddressPairs []reservedfixedips.AllowedAddressPairs `json:"allowed_address_pairs"`
	InstanceID          string                                 `json:"instance_id"`
	PortID              string                                 `json:"port_id"`
}
