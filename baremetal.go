package edgecloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

const (
	bmInstancesBasePathV1     = "/v1/bminstances"
	bmCapacityBasePathV1      = "/v1/bmcapacity"
	bmCheckLimitsSupPath      = "check_limits"
	bmRebuildSubPath          = "rebuild"
	bmAvailableFlavorsSubPath = "available_flavors"
)

// BareMetalService is an interface for creating and managing bare metal Instances with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/instances
type BareMetalService interface {
	BareMetalListInstances(context.Context, *BareMetalInstancesListOpts) ([]Instance, *Response, error)
	BareMetalCreateInstance(context.Context, *BareMetalServerCreateRequest) (*TaskResponse, *Response, error)
	BareMetalRebuildInstance(context.Context, string, *BareMetalRebuildRequest) (*TaskResponse, *Response, error)
	BareMetalListFlavors(context.Context, *BareMetalFlavorsOpts, *BareMetalFlavorsRequest) ([]BareMetalFlavor, *Response, error)
	BareMetalGetCountAvailableNodes(context.Context) (*BareMetalCapacity, *Response, error)
	BareMetalCheckQuotasForInstanceCreation(context.Context, *BareMetalQuotaCheckRequest) (Quota, *Response, error)
}

// BareMetalInstancesListOpts allows the filtering and sorting of paginated collections through the API.
type BareMetalInstancesListOpts struct {
	Name                    string `url:"name,omitempty" validate:"omitempty"`
	ProfileName             string `url:"profile_name,omitempty" validate:"omitempty"`
	OnlyWithFixedExternalIP bool   `url:"only_with_fixed_external_ip,omitempty" validate:"omitempty"`
	Limit                   int    `url:"limit,omitempty"  validate:"omitempty"`
	Offset                  int    `url:"offset,omitempty"  validate:"omitempty"`
	FlavorID                string `url:"flavor_id,omitempty"  validate:"omitempty"`
	Status                  string `url:"status,omitempty" validate:"omitempty"`
	ChangesBefore           string `url:"changes-before,omitempty" validate:"omitempty"`
	IP                      string `url:"ip,omitempty"  validate:"omitempty"`
	UUID                    string `url:"uuid,omitempty"  validate:"omitempty"`
	MetadataKV              string `url:"metadata_kv,omitempty"  validate:"omitempty"`
	MetadataK               string `url:"metadata_k,omitempty"  validate:"omitempty"`
	OrderBy                 string `url:"order_by,omitempty"  validate:"omitempty"`
}

// BareMetalRebuildOpts allows the filtering and sorting of paginated collections through the API.
type BareMetalRebuildOpts struct {
	InstanceID string `json:"instance_id" required:"true" validate:"uuid4"`
}

type BareMetalInterfaceOpts struct {
	Type       InterfaceType        `json:"type" validate:"omitempty,enum"`
	NetworkID  string               `json:"network_id,omitempty" validate:"rfe=Type:subnet,omitempty,uuid4"`
	SubnetID   string               `json:"subnet_id,omitempty" validate:"rfe=Type:subnet,omitempty,uuid4"`
	PortID     string               `json:"port_id,omitempty" validate:"rfe=Type:reserved_fixed_ip,allowed_without_all=NetworkID SubnetID,omitempty,uuid4"`
	FloatingIP *InterfaceFloatingIP `json:"floating_ip,omitempty" validate:"omitempty,dive"`
}

// BareMetalServerCreateRequest represents a request to create an bare metal server.
type BareMetalServerCreateRequest struct {
	KeypairName   string                   `json:"keypair_name,omitempty"`
	AppTemplateID string                   `json:"apptemplate_id,omitempty"`
	Flavor        string                   `json:"flavor" required:"true"`
	Metadata      Metadata                 `json:"metadata,omitempty" validate:"omitempty,dive"`
	NameTemplates []string                 `json:"name_templates,omitempty" validate:"required_without=Names"`
	Username      string                   `json:"username,omitempty" validate:"omitempty,required_with=Password"`
	Password      string                   `json:"password,omitempty" validate:"omitempty,required_with=Username"`
	Names         []string                 `json:"names,omitempty" validate:"required_without=NameTemplates"`
	Interfaces    []BareMetalInterfaceOpts `json:"interfaces" required:"true" validate:"required,dive"`
	ImageID       string                   `json:"image_id" validate:"omitempty"`
	UserData      string                   `json:"user_data,omitempty" validate:"omitempty,base64"`
	AppConfig     map[string]interface{}   `json:"app_config,omitempty" validate:"omitempty,dive"`
}
type BareMetalRebuildRequest struct {
	ImageID string `json:"image_id" validate:"omitempty"`
}

type BareMetalFlavorsOpts struct {
	IncludePrices bool `url:"include_prices,omitempty"`
}

type BareMetalFlavorsRequest struct {
	ImageID string `json:"image_id" url:"image_id,omitempty"`
}

type bareMetalFlavorsRoot struct {
	Count   int               `json:"count"`
	Flavors []BareMetalFlavor `json:"results"`
}

// BareMetalFlavor  represents an EdgecenterCloud BareMetalFlavor.
type BareMetalFlavor struct {
	FlavorName          string              `json:"flavor_name"`
	PricePerMonth       int                 `json:"price_per_month"`
	Disabled            bool                `json:"disabled"`
	HardwareDescription HardwareDescription `json:"hardware_description,omitempty"`
	CurrencyCode        string              `json:"currency_code"`
	PriceStatus         string              `json:"price_status"`
	FlavorID            string              `json:"flavor_id"`
	PricePerHour        int                 `json:"price_per_hour"`
	ResourceClass       string              `json:"resource_class"`
}

type BareMetalCapacity struct {
	Capacity map[string]int `json:"capacity"`
}

type BareMetalQuotaCheckRequest struct {
	Flavor     string                   `json:"flavor" required:"true"`
	Interfaces []BareMetalInterfaceOpts `json:"interfaces" required:"true" validate:"required,dive"`
}

func (s *InstancesServiceOp) BareMetalCheckQuotasForInstanceCreation(ctx context.Context, reqBody *BareMetalQuotaCheckRequest) (Quota, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}
	path := s.client.addProjectRegionPath(bmInstancesBasePathV1)
	path = fmt.Sprintf("%s/%s", path, bmCheckLimitsSupPath)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	quotas := new(Quota)
	resp, err := s.client.Do(ctx, req, quotas)
	if err != nil {
		return nil, resp, err
	}

	return *quotas, resp, err
}

func (s *InstancesServiceOp) BareMetalGetCountAvailableNodes(ctx context.Context) (*BareMetalCapacity, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(bmCapacityBasePathV1)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	capacity := new(BareMetalCapacity)
	resp, err := s.client.Do(ctx, req, capacity)
	if err != nil {
		return nil, resp, err
	}

	return capacity, resp, err
}

func (s *InstancesServiceOp) BareMetalListFlavors(ctx context.Context, opts *BareMetalFlavorsOpts, reqBody *BareMetalFlavorsRequest) ([]BareMetalFlavor, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}
	var err error
	path := s.client.addProjectRegionPath(bmInstancesBasePathV1)
	path = fmt.Sprintf("%s/%s", path, bmAvailableFlavorsSubPath)
	if opts != nil {
		path, err = addOptions(path, opts)
		if err != nil {
			return nil, nil, err
		}
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	root := new(bareMetalFlavorsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Flavors, resp, err
}

func (s *InstancesServiceOp) BareMetalRebuildInstance(ctx context.Context, instanceID string, reqBody *BareMetalRebuildRequest) (*TaskResponse, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(bmInstancesBasePathV1)
	path = fmt.Sprintf("%s/%s/%s", path, instanceID, bmRebuildSubPath)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	tasks := new(TaskResponse)
	resp, err := s.client.Do(ctx, req, tasks)
	if err != nil {
		return nil, resp, err
	}

	return tasks, resp, err
}

func (s *InstancesServiceOp) BareMetalCreateInstance(ctx context.Context, reqBody *BareMetalServerCreateRequest) (*TaskResponse, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(bmInstancesBasePathV1)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	tasks := new(TaskResponse)
	resp, err := s.client.Do(ctx, req, tasks)
	if err != nil {
		return nil, resp, err
	}

	return tasks, resp, err
}

func (s *InstancesServiceOp) BareMetalListInstances(ctx context.Context, opts *BareMetalInstancesListOpts) ([]Instance, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	validate := validator.New()
	err := validate.Struct(opts)
	if err != nil {
		return nil, nil, err
	}

	path := s.client.addProjectRegionPath(bmInstancesBasePathV1)
	if opts != nil {
		path, err = addOptions(path, opts)
		if err != nil {
			return nil, nil, err
		}
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(instancesRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Instances, resp, err
}
