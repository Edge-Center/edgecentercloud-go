package edgecloud

import (
	"context"
	"fmt"
	"net/http"
)

const (
	regionsBasePath = "/v1/regions"
)

// RegionsService is an interface for creating and managing Regions with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/regions
type RegionsService interface {
	List(context.Context, *RegionListOptions) ([]Region, *Response, error)
	Get(context.Context, string, *RegionGetOptions) (*Region, *Response, error)
}

// RegionsServiceOp handles communication with Regions methods of the EdgecenterCloud API.
type RegionsServiceOp struct {
	client *Client
}

var _ RegionsService = &RegionsServiceOp{}

// Region represents a EdgecenterCloud Region configuration.
type Region struct {
	HasKVM               bool               `json:"has_kvm"`
	HasK8S               bool               `json:"has_k8s"`
	DisplayName          string             `json:"display_name"`
	HasBaremetal         bool               `json:"has_baremetal"`
	VLANPhysicalNetwork  string             `json:"vlan_physical_network"`
	EndpointType         RegionEndpointType `json:"endpoint_type"`
	Zone                 Zone               `json:"zone"`
	SpiceProxyURL        string             `json:"spice_proxy_url"`
	KeystoneID           int                `json:"keystone_id"`
	KeystoneName         string             `json:"keystone_name"`
	Keystone             Keystone           `json:"keystone"`
	MetricsDatabaseID    int                `json:"metrics_database_id"`
	SerialProxyURL       string             `json:"serial_proxy_url"`
	ID                   int                `json:"id"`
	TaskID               *string            `json:"task_id"`
	ExternalNetworkID    string             `json:"external_network_id"`
	AccessLevel          string             `json:"access_level"`
	CreatedOn            string             `json:"created_on"`
	AvailableVolumeTypes []string           `json:"available_volume_types"`
	Country              string             `json:"country"`
	NoVNSProxyURL        string             `json:"novnc_proxy_url"`
	K8SMasterFlavorID    string             `json:"k8s_master_flavor_id"`
	State                RegionState        `json:"state"`
	MetricsDatabase      MetricsDatabase    `json:"metrics_database"`
}

type Zone string

const (
	ZoneAPAC         Zone = "APAC"
	ZoneEMEA         Zone = "EMEA"
	ZoneAmericas     Zone = "AMERICAS"
	ZoneRussiaAndCIS Zone = "RUSSIA_AND_CIS"
)

type RegionEndpointType string

const (
	RegionEndpointTypePublic   RegionEndpointType = "public"
	RegionEndpointTypeInternal RegionEndpointType = "internal"
	RegionEndpointTypeAdmin    RegionEndpointType = "admin"
)

type Keystone struct {
	ID                        int           `json:"id"`
	URL                       string        `json:"url"`
	State                     KeystoneState `json:"state"`
	KeystoneFederatedDomainID string        `json:"keystone_federated_domain_id"`
	CreatedOn                 string        `json:"created_on"`
	AdminPassword             string        `json:"admin_password"`
}

type KeystoneState string

const (
	KeystoneStateNew               KeystoneState = "NEW"
	KeystoneStateInitializedFailed KeystoneState = "INITIALIZATION_FAILED"
	KeystoneStateInitialized       KeystoneState = "INITIALIZED"
	KeystoneStateDeleted           KeystoneState = "DELETED"
)

type RegionState string

const (
	RegionStateNew            RegionState = "NEW"
	RegionStateInactive       RegionState = "INACTIVE"
	RegionStateActive         RegionState = "ACTIVE"
	RegionStateMaintenance    RegionState = "MAINTENANCE"
	RegionStateDeleting       RegionState = "DELETING"
	RegionStateDeletionFailed RegionState = "DELETION_FAILED"
	RegionStateDeleted        RegionState = "DELETED"
)

type MetricsDatabase struct {
	ID int `json:"id"`
}

// RegionGetOptions specifies the optional query parameters to Get method.
type RegionGetOptions struct {
	ShowVolumeTypes bool `url:"show_volume_types,omitempty"  validate:"omitempty"`
}

// RegionListOptions specifies the optional query parameters to List method.
type RegionListOptions struct {
	ShowVolumeTypes bool `url:"show_volume_types,omitempty"  validate:"omitempty"`
}

// regionsRoot represents a Region root.
type regionsRoot struct {
	Count  int
	Region []Region `json:"results"`
}

// List get regions.
func (s *RegionsServiceOp) List(ctx context.Context, opts *RegionListOptions) ([]Region, *Response, error) {
	path, err := addOptions(regionsBasePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(regionsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Region, resp, err
}

// Get retrieves a single region by its ID.
func (s *RegionsServiceOp) Get(ctx context.Context, regionID string, opts *RegionGetOptions) (*Region, *Response, error) {
	path := fmt.Sprintf("%s/%s", regionsBasePath, regionID)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	project := new(Region)
	resp, err := s.client.Do(ctx, req, project)
	if err != nil {
		return nil, resp, err
	}

	return project, resp, err
}
