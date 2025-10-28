package edgecloud

import (
	"context"
	"net/http"
)

const (
	AvailabilityZoneBasePath = "/v1/availability_zones"
)

// AvailabilityZonesService is an interface for managing Availability Zones with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/Availability-Zones
type AvailabilityZonesService interface {
	List(context.Context) (*AvailabilityZonesList, *Response, error)
}

// AvailabilityZonesServiceOp handles communication with Availability Zones methods of the EdgecenterCloud API.
type AvailabilityZonesServiceOp struct {
	client *Client
}

var _ AvailabilityZonesService = &AvailabilityZonesServiceOp{}

// AvailabilityZonesList represents an EdgecenterCloud availability zones list.
type AvailabilityZonesList struct {
	RegionID          int      `json:"region_id"`
	AvailabilityZones []string `json:"availability_zones"`
}

// List get availability zones in a region.
func (s *AvailabilityZonesServiceOp) List(ctx context.Context) (*AvailabilityZonesList, *Response, error) {
	if resp, err := s.client.ValidateRegion(); err != nil {
		return nil, resp, err
	}

	path := s.client.addRegionPath(AvailabilityZoneBasePath)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	az := new(AvailabilityZonesList)
	resp, err := s.client.Do(ctx, req, az)
	if err != nil {
		return nil, resp, err
	}

	return az, resp, nil
}
