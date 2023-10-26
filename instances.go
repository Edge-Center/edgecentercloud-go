package edgecloud

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

const (
	instancesBasePathV1       = "/v1/instances"
	instancesBasePathV2       = "/v2/instances"
	instanceMetadataPath      = "metadata"
	instancesCheckLimitsPath  = "check_limits"
	instancesChangeFlavorPath = "changeflavor"
)

// InstancesService is an interface for creating and managing Instances with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/instances
type InstancesService interface {
	Get(context.Context, string) (*Instance, *Response, error)
	Create(context.Context, *InstanceCreateRequest) (*TaskResponse, *Response, error)
	Delete(context.Context, string, *InstanceDeleteOptions) (*TaskResponse, *Response, error)

	CheckLimits(context.Context, *InstanceCheckLimitsRequest) (*map[string]int, *Response, error)

	UpdateFlavor(context.Context, string, *InstanceFlavorUpdateRequest) (*TaskResponse, *Response, error)

	InstanceMetadata
}

type InstanceMetadata interface {
	MetadataGet(context.Context, string) (*MetadataDetailed, *Response, error)
	MetadataCreate(context.Context, string, *MetadataCreateRequest) (*Response, error)
}

// InstancesServiceOp handles communication with Instances methods of the EdgecenterCloud API.
type InstancesServiceOp struct {
	client *Client
}

var _ InstancesService = &InstancesServiceOp{}

// Instance represents an EdgecenterCloud Instance.
type Instance struct {
	ID               string                       `json:"instance_id"`
	Name             string                       `json:"instance_name"`
	Addresses        map[string][]InstanceAddress `json:"addresses"`
	CreatedAt        string                       `json:"instance_created,omitempty"`
	CreatorTaskID    string                       `json:"creator_task_id,omitempty"`
	Description      string                       `json:"instance_description,omitempty"`
	Flavor           *Flavor                      `json:"flavor"`
	KeypairName      string                       `json:"keypair_name,omitempty"`
	Metadata         Metadata                     `json:"metadata"`
	MetadataDetailed []MetadataDetailed           `json:"metadata_detailed,omitempty"`
	ProjectID        int                          `json:"project_id"`
	Region           string                       `json:"region"`
	RegionID         int                          `json:"region_id"`
	SecurityGroups   []InstanceSecurityGroup      `json:"security_groups"`
	Status           string                       `json:"status,omitempty"` // todo: need to implement new status type
	TaskID           string                       `json:"task_id"`
	TaskState        string                       `json:"task_state,omitempty"`
	VMState          string                       `json:"vm_state,omitempty"` // todo: need to implement new vm_state type
	Volumes          []InstanceVolume             `json:"volumes"`
}

// InstanceVolume represent an instance volume struct.
type InstanceVolume struct {
	ID                  string `json:"id"`
	DeleteOnTermination bool   `json:"delete_on_termination"`
}

// InstanceSecurityGroup represent an instance firewall struct.
type InstanceSecurityGroup struct {
	Name string `json:"name"`
}

type AddressType string

const (
	AddressTypeFixed    AddressType = "fixed"
	AddressTypeFloating AddressType = "floating"
)

// InstanceAddress represent an instance network struct.
type InstanceAddress struct {
	Type       string `json:"type"`
	SubnetName string `json:"subnet_name"`
	SubnetID   string `json:"subnet_id"`
	Address    net.IP `json:"addr"`
}

type InterfaceType string

const (
	AnySubnetInterfaceType InterfaceType = "any_subnet"
	ExternalInterfaceType  InterfaceType = "external"
	ReservedFixedIPType    InterfaceType = "reserved_fixed_ip"
	SubnetInterfaceType    InterfaceType = "subnet"
)

type InstanceInterface struct {
	Type           InterfaceType                  `json:"type,omitempty" validate:"omitempty,enum"`
	NetworkID      string                         `json:"network_id,omitempty" validate:"rfe=Type:subnet;any_subnet,omitempty,uuid4"`
	FloatingIP     *InterfaceFloatingIP           `json:"floating_ip,omitempty" validate:"omitempty,dive"`
	PortID         string                         `json:"port_id,omitempty" validate:"rfe=Type:reserved_fixed_ip,allowed_without_all=NetworkID SubnetID,omitempty,uuid4"`
	SubnetID       string                         `json:"subnet_id,omitempty" validate:"rfe=Type:subnet,omitempty,uuid4"`
	SecurityGroups []InstanceSecurityGroupsCreate `json:"security_groups"`
}

// InstanceSecurityGroupsCreate represent an instance firewall create struct.
type InstanceSecurityGroupsCreate struct {
	ID string `json:"id"`
}

// InstanceVolumeCreate represent a instance volume create struct.
type InstanceVolumeCreate struct {
	Source        VolumeSource `json:"source" required:"true" validate:"required,enum"`
	BootIndex     int          `json:"boot_index"`
	TypeName      VolumeType   `json:"type_name,omitempty" validate:"omitempty"`
	Size          int          `json:"size,omitempty" validate:"rfe=Source:image;new-volume,sfe=Source:snapshot;existing-volume"`
	Name          string       `json:"name,omitempty" validate:"omitempty"`
	AttachmentTag string       `json:"attachment_tag,omitempty" validate:"omitempty"`
	ImageID       string       `json:"image_id,omitempty" validate:"rfe=Source:image,sfe=Source:snapshot;apptemplate;existing-volume;new-volume,allowed_without_all=SnapshotID VolumeID,omitempty,uuid4"`
	VolumeID      string       `json:"volume_id,omitempty" validate:"rfe=Source:existing-volume,sfe=Source:image;shapshot;apptemplate;new-volume,allowed_without_all=ImageID SnapshotID,omitempty,uuid4"`
	SnapshotID    string       `json:"snapshot_id,omitempty" validate:"rfe=Source:snapshot,sfe=Source:image;existing-volume;new-volume;apptemplate,allowed_without_all=ImageID VolumeID,omitempty,uuid4"`
	AppTemplateID string       `json:"apptemplate_id,omitempty" validate:"rfe=Source:apptemplate,sfe=Source:image;existing-volume;new-volume;snapshot,allowed_without_all=ImageID VolumeID,omitempty,uuid4"`
	Metadata      Metadata     `json:"metadata,omitempty" validate:"omitempty,dive"`
}

// InstanceCreateRequest represents a request to create an Instance.
type InstanceCreateRequest struct {
	Names          []string                       `json:"names,omitempty" validate:"required_without=NameTemplates"`
	Flavor         string                         `json:"flavor" required:"true"`
	NameTemplates  []string                       `json:"name_templates,omitempty" validate:"required_without=Names"`
	KeypairName    string                         `json:"keypair_name,omitempty"`
	UserData       string                         `json:"user_data,omitempty" validate:"omitempty,base64"`
	Username       string                         `json:"username,omitempty" validate:"omitempty,required_with=Password"`
	Password       string                         `json:"password,omitempty" validate:"omitempty,required_with=Username"`
	Interfaces     []InstanceInterface            `json:"interfaces" required:"true" validate:"required,dive"`
	SecurityGroups []InstanceSecurityGroupsCreate `json:"security_groups,omitempty" validate:"omitempty,dive,uuid4"`
	Metadata       Metadata                       `json:"metadata,omitempty" validate:"omitempty,dive"`
	Configuration  map[string]interface{}         `json:"configuration,omitempty" validate:"omitempty,dive"`
	ServerGroupID  string                         `json:"servergroup_id,omitempty" validate:"omitempty,uuid4"`
	AllowAppPorts  bool                           `json:"allow_app_ports,omitempty"`
	Volumes        []InstanceVolumeCreate         `json:"volumes" required:"true" validate:"required,dive"`
}

// InstanceDeleteOptions specifies the optional query parameters to Delete method.
type InstanceDeleteOptions struct {
	Volumes          []string `url:"volumes,omitempty" validate:"omitempty,dive,uuid4" delimiter:"comma"`
	DeleteFloatings  bool     `url:"delete_floatings,omitempty"  validate:"omitempty,allowed_without=FloatingIPs"`
	FloatingIPs      []string `url:"floatings,omitempty" validate:"omitempty,allowed_without=DeleteFloatings,dive,uuid4" delimiter:"comma"`
	ReservedFixedIPs []string `url:"reserved_fixed_ips,omitempty" validate:"omitempty,dive,uuid4" delimiter:"comma"`
}

type InstanceCheckLimitsVolume struct {
	Source     VolumeSource `json:"source" required:"true" validate:"required,enum"`
	TypeName   VolumeType   `json:"type_name,omitempty" validate:"omitempty"`
	Size       int          `json:"size,omitempty" validate:"omitempty"`
	SnapshotID string       `json:"snapshot_id,omitempty" validate:"omitempty"`
	ImageID    string       `json:"image_id,omitempty" validate:"omitempty"`
}

// InstanceCheckLimitsRequest represents a request to check the limits of an instance.
type InstanceCheckLimitsRequest struct {
	Names         []string                    `json:"names,omitempty" validate:"required_without=NameTemplates"`
	NameTemplates []string                    `json:"name_templates,omitempty" validate:"required_without=Names"`
	Flavor        string                      `json:"flavor,omitempty"`
	Interfaces    []InstanceInterface         `json:"interfaces,omitempty" required:"true" validate:"required,dive"`
	Volumes       []InstanceCheckLimitsVolume `json:"volumes,omitempty" required:"true" validate:"required,dive"`
}

// InstanceFlavorUpdateRequest represents a request to change the flavor of the instance.
type InstanceFlavorUpdateRequest struct {
	FlavorID string `json:"flavor_id" required:"true" validate:"required"`
}

// Get individual Instance.
func (s *InstancesServiceOp) Get(ctx context.Context, instanceID string) (*Instance, *Response, error) {
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addServicePath(instancesBasePathV1), instanceID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	instance := new(Instance)
	resp, err := s.client.Do(ctx, req, instance)
	if err != nil {
		return nil, resp, err
	}

	return instance, resp, err
}

// Create an Instance.
func (s *InstancesServiceOp) Create(ctx context.Context, createRequest *InstanceCreateRequest) (*TaskResponse, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addServicePath(instancesBasePathV2)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
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

// Delete the Instance.
func (s *InstancesServiceOp) Delete(ctx context.Context, instanceID string, opts *InstanceDeleteOptions) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addServicePath(instancesBasePathV1), instanceID)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
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

// MetadataGet instance detailed metadata (tags).
func (s *InstancesServiceOp) MetadataGet(ctx context.Context, instanceID string) (*MetadataDetailed, *Response, error) {
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addServicePath(instancesBasePathV1)
	path = fmt.Sprintf("%s/%s/%s", path, instanceID, instanceMetadataPath)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	metadata := new(MetadataDetailed)
	resp, err := s.client.Do(ctx, req, metadata)
	if err != nil {
		return nil, resp, err
	}

	return metadata, resp, err
}

// MetadataCreate instance metadata (tags).
func (s *InstancesServiceOp) MetadataCreate(ctx context.Context, instanceID string, metadata *MetadataCreateRequest) (*Response, error) {
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	path := s.client.addServicePath(instancesBasePathV1)
	path = fmt.Sprintf("%s/%s/%s", path, instanceID, instanceMetadataPath)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, metadata)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// CheckLimits check a quota for instance creation.
func (s *InstancesServiceOp) CheckLimits(ctx context.Context, checkLimitsRequest *InstanceCheckLimitsRequest) (*map[string]int, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addServicePath(instancesBasePathV2)
	path = fmt.Sprintf("%s/%s", path, instancesCheckLimitsPath)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, checkLimitsRequest)
	if err != nil {
		return nil, nil, err
	}

	limits := new(map[string]int)
	resp, err := s.client.Do(ctx, req, limits)
	if err != nil {
		return nil, resp, err
	}

	return limits, resp, nil
}

// UpdateFlavor changes the flavor of the server instance.
func (s *InstancesServiceOp) UpdateFlavor(ctx context.Context, instanceID string, instanceFlavorUpdateRequest *InstanceFlavorUpdateRequest) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addServicePath(instancesBasePathV1)
	path = fmt.Sprintf("%s/%s/%s", path, instanceID, instancesChangeFlavorPath)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, instanceFlavorUpdateRequest)
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
