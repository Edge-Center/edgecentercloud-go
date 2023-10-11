package edgecloud

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/google/uuid"
)

const (
	instancesBasePathV1 = "/v1/instances"
	instancesBasePathV2 = "/v2/instances"
)

// InstancesService is an interface for creating and managing Instances with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/instances
type InstancesService interface {
	Get(context.Context, string, *ServicePath) (*Instance, *Response, error)
	Create(context.Context, *InstanceCreateRequest, *ServicePath) (*TaskResponse, *Response, error)
	Delete(context.Context, string, *ServicePath) (*TaskResponse, *Response, error)
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
	Metadata         map[string]interface{}       `json:"metadata"`
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

type FloatingIPSource string

const (
	NewFloatingIP      FloatingIPSource = "new"
	ExistingFloatingIP FloatingIPSource = "existing"
)

type InterfaceFloatingIP struct {
	Source             FloatingIPSource `json:"source" validate:"required,enum"`
	ExistingFloatingID string           `json:"existing_floating_id" validate:"rfe=Source:existing,sfe=Source:new,omitempty,uuid"`
}

type InstanceInterface struct {
	Type       InterfaceType        `json:"type,omitempty" validate:"omitempty,enum"`
	NetworkID  string               `json:"network_id,omitempty" validate:"rfe=Type:subnet;any_subnet,omitempty,uuid4"`
	FloatingIP *InterfaceFloatingIP `json:"floating_ip,omitempty" validate:"omitempty,dive"`
	PortID     string               `json:"port_id,omitempty" validate:"rfe=Type:reserved_fixed_ip,allowed_without_all=NetworkID SubnetID,omitempty,uuid4"`
	SubnetID   string               `json:"subnet_id,omitempty" validate:"rfe=Type:subnet,omitempty,uuid4"`
}

// InstanceSecurityGroupsCreate represent an instance firewall create struct.
type InstanceSecurityGroupsCreate struct {
	ID string `json:"id"`
}

type VolumeSource string

const (
	NewVolume      VolumeSource = "new-volume"
	Image          VolumeSource = "image"
	Snapshot       VolumeSource = "snapshot"
	ExistingVolume VolumeSource = "existing-volume"
	AppTemplate    VolumeSource = "apptemplate"
)

// InstanceVolumeCreate represent a instance volume create struct.
type InstanceVolumeCreate struct {
	Source        VolumeSource           `json:"source" required:"true" validate:"required,enum"`
	BootIndex     int                    `json:"boot_index"`
	TypeName      VolumeType             `json:"type_name,omitempty" validate:"omitempty"`
	Size          int                    `json:"size,omitempty" validate:"rfe=Source:image;new-volume,sfe=Source:snapshot;existing-volume"`
	Name          string                 `json:"name,omitempty" validate:"omitempty"`
	AttachmentTag string                 `json:"attachment_tag,omitempty" validate:"omitempty"`
	ImageID       string                 `json:"image_id,omitempty" validate:"rfe=Source:image,sfe=Source:snapshot;apptemplate;existing-volume;new-volume,allowed_without_all=SnapshotID VolumeID,omitempty,uuid4"`
	VolumeID      string                 `json:"volume_id,omitempty" validate:"rfe=Source:existing-volume,sfe=Source:image;shapshot;apptemplate;new-volume,allowed_without_all=ImageID SnapshotID,omitempty,uuid4"`
	SnapshotID    string                 `json:"snapshot_id,omitempty" validate:"rfe=Source:snapshot,sfe=Source:image;existing-volume;new-volume;apptemplate,allowed_without_all=ImageID VolumeID,omitempty,uuid4"`
	AppTemplateID string                 `json:"apptemplate_id,omitempty" validate:"rfe=Source:apptemplate,sfe=Source:image;existing-volume;new-volume;snapshot,allowed_without_all=ImageID VolumeID,omitempty,uuid4"`
	Metadata      map[string]interface{} `json:"metadata,omitempty" validate:"omitempty,dive"`
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
	Metadata       map[string]interface{}         `json:"metadata,omitempty" validate:"omitempty,dive"`
	Configuration  map[string]interface{}         `json:"configuration,omitempty" validate:"omitempty,dive"`
	ServerGroupID  string                         `json:"servergroup_id,omitempty" validate:"omitempty,uuid4"`
	AllowAppPorts  bool                           `json:"allow_app_ports,omitempty"`
	Volumes        []InstanceVolumeCreate         `json:"volumes" required:"true" validate:"required,dive"`
}

// instanceRoot represents an Instance root.
type instanceRoot struct {
	Instance *Instance     `json:"instance"`
	Tasks    *TaskResponse `json:"tasks"`
}

// Get individual Instance.
func (s *InstancesServiceOp) Get(ctx context.Context, instanceID string, p *ServicePath) (*Instance, *Response, error) {
	if _, err := uuid.Parse(instanceID); err != nil {
		return nil, nil, NewArgError("instanceID", "should be the correct UUID")
	}

	if p == nil {
		return nil, nil, NewArgError("ServicePath", "cannot be nil")
	}

	path := addServicePath(instancesBasePathV1, p)
	path = fmt.Sprintf("%s/%s", path, instanceID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(instanceRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Instance, resp, err
}

// Create an Instance.
func (s *InstancesServiceOp) Create(ctx context.Context, createRequest *InstanceCreateRequest, p *ServicePath) (*TaskResponse, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	if p == nil {
		return nil, nil, NewArgError("ServicePath", "cannot be nil")
	}

	path := addServicePath(instancesBasePathV2, p)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(instanceRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Tasks, resp, err
}

// Delete the Instance.
func (s *InstancesServiceOp) Delete(ctx context.Context, instanceID string, p *ServicePath) (*TaskResponse, *Response, error) {
	if _, err := uuid.Parse(instanceID); err != nil {
		return nil, nil, NewArgError("instanceID", "should be the correct UUID")
	}

	if p == nil {
		return nil, nil, NewArgError("ServicePath", "cannot be nil")
	}

	path := addServicePath(instancesBasePathV1, p)
	path = fmt.Sprintf("%s/%s", path, instanceID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(instanceRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Tasks, resp, err
}
