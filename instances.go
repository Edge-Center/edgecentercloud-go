package edgecloud

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

const (
	instancesBasePathV1 = "/v1/instances"
	instancesBasePathV2 = "/v2/instances"
)

const (
	instancesCheckLimits           = "check_limits"
	instancesChangeFlavor          = "changeflavor"
	instancesAvailableFlavors      = "available_flavors"
	instancesAvailableNames        = "available_names"
	instancesPorts                 = "ports"
	instancesStart                 = "start"
	instancesStop                  = "stop"
	instancesPowercycle            = "powercycle"
	instancesReboot                = "reboot"
	instancesSuspend               = "suspend"
	instancesResume                = "resume"
	instancesMetrics               = "metrics"
	instancesInstances             = "instances"
	instancesSecurityGroups        = "securitygroups"
	instancesAddSecurityGroup      = "addsecuritygroup"
	instancesDelSecurityGroup      = "delsecuritygroup"
	instancesGetConsole            = "get_console"
	instancesAttachInterface       = "attach_interface"
	instancesDetachInterface       = "detach_interface"
	instancesInterfaces            = "interfaces"
	instancesPutIntoServerGroup    = "put_into_servergroup"
	instancesRemoveFromServerGroup = "remove_from_servergroup"
)

// InstancesService is an interface for creating and managing Instances with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/instances
type InstancesService interface {
	List(context.Context, *InstanceListOptions) ([]Instance, *Response, error)
	Get(context.Context, string) (*Instance, *Response, error)
	Create(context.Context, *InstanceCreateRequest) (*TaskResponse, *Response, error)
	Delete(context.Context, string, *InstanceDeleteOptions) (*TaskResponse, *Response, error)
	CheckLimits(context.Context, *InstanceCheckLimitsRequest) (*map[string]int, *Response, error)
	AvailableNames(context.Context) (*InstanceAvailableNames, *Response, error)
	Rename(context.Context, string, *Name) (*Instance, *Response, error)
	PortsList(context.Context, string) ([]InstancePort, *Response, error)
	MetricsList(context.Context, string, *InstanceMetricsListRequest) ([]InstanceMetrics, *Response, error)
	GetConsole(context.Context, string) (*RemoteConsole, *Response, error)
	AttachInterface(context.Context, string, *InstanceAttachInterfaceRequest) (*TaskResponse, *Response, error)
	DetachInterface(context.Context, string, *InstanceDetachInterfaceRequest) (*TaskResponse, *Response, error)
	InterfaceList(context.Context, string) ([]RouterInterface, *Response, error)
	PutIntoServerGroup(context.Context, string, *InstancePutIntoServerGroupRequest) (*TaskResponse, *Response, error)
	RemoveFromServerGroup(context.Context, string) (*TaskResponse, *Response, error)

	InstanceAction
	InstanceFlavor
	InstanceSecurityGroup
	InstanceMetadata
}

type InstanceAction interface {
	InstanceStart(context.Context, string) (*Instance, *Response, error)
	InstanceStop(context.Context, string) (*Instance, *Response, error)
	InstancePowercycle(context.Context, string) (*Instance, *Response, error)
	InstanceReboot(context.Context, string) (*Instance, *Response, error)
	InstanceSuspend(context.Context, string) (*Instance, *Response, error)
	InstanceResume(context.Context, string) (*Instance, *Response, error)
}

type InstanceFlavor interface {
	UpdateFlavor(context.Context, string, *InstanceFlavorUpdateRequest) (*TaskResponse, *Response, error)
	AvailableFlavors(context.Context, *InstanceCheckFlavorVolumeRequest, *FlavorsOptions) ([]Flavor, *Response, error)
	AvailableFlavorsToResize(context.Context, string, *FlavorsOptions) ([]Flavor, *Response, error)
}

type InstanceSecurityGroup interface {
	FilterBySecurityGroup(context.Context, string) ([]Instance, *Response, error)
	SecurityGroupList(context.Context, string) ([]IDName, *Response, error)
	SecurityGroupAssign(context.Context, string, *AssignSecurityGroupRequest) (*Response, error)
	SecurityGroupUnAssign(context.Context, string, *AssignSecurityGroupRequest) (*Response, error)
}

type InstanceMetadata interface {
	MetadataGet(context.Context, string) (*MetadataDetailed, *Response, error)
	MetadataList(context.Context, string) ([]MetadataDetailed, *Response, error)
	MetadataCreate(context.Context, string, *Metadata) (*Response, error)
	MetadataUpdate(context.Context, string, *Metadata) (*Response, error)
	MetadataDeleteItem(context.Context, string, *MetadataItemOptions) (*Response, error)
	MetadataGetItem(context.Context, string, *MetadataItemOptions) (*MetadataDetailed, *Response, error)
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
	SecurityGroups   []Name                       `json:"security_groups"`
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
	Type           InterfaceType        `json:"type,omitempty" validate:"omitempty,enum"`
	NetworkID      string               `json:"network_id,omitempty" validate:"rfe=Type:subnet;any_subnet,omitempty,uuid4"`
	FloatingIP     *InterfaceFloatingIP `json:"floating_ip,omitempty" validate:"omitempty,dive"`
	PortID         string               `json:"port_id,omitempty" validate:"rfe=Type:reserved_fixed_ip,allowed_without_all=NetworkID SubnetID,omitempty,uuid4"`
	SubnetID       string               `json:"subnet_id,omitempty" validate:"rfe=Type:subnet,omitempty,uuid4"`
	SecurityGroups []ID                 `json:"security_groups"`
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
	Names          []string               `json:"names,omitempty" validate:"required_without=NameTemplates"`
	Flavor         string                 `json:"flavor" required:"true"`
	NameTemplates  []string               `json:"name_templates,omitempty" validate:"required_without=Names"`
	KeypairName    string                 `json:"keypair_name,omitempty"`
	UserData       string                 `json:"user_data,omitempty" validate:"omitempty,base64"`
	Username       string                 `json:"username,omitempty" validate:"omitempty,required_with=Password"`
	Password       string                 `json:"password,omitempty" validate:"omitempty,required_with=Username"`
	Interfaces     []InstanceInterface    `json:"interfaces" required:"true" validate:"required,dive"`
	SecurityGroups []ID                   `json:"security_groups,omitempty" validate:"omitempty,dive,uuid4"`
	Metadata       Metadata               `json:"metadata,omitempty" validate:"omitempty,dive"`
	Configuration  map[string]interface{} `json:"configuration,omitempty" validate:"omitempty,dive"`
	ServerGroupID  string                 `json:"servergroup_id,omitempty" validate:"omitempty,uuid4"`
	AllowAppPorts  bool                   `json:"allow_app_ports,omitempty"`
	Volumes        []InstanceVolumeCreate `json:"volumes" required:"true" validate:"required,dive"`
}

// InstanceDeleteOptions specifies the optional query parameters to Delete method.
type InstanceDeleteOptions struct {
	Volumes          []string `url:"volumes,omitempty" validate:"omitempty,dive,uuid4" delimiter:"comma"`
	DeleteFloatings  bool     `url:"delete_floatings,omitempty"  validate:"omitempty,allowed_without=FloatingIPs"`
	FloatingIPs      []string `url:"floatings,omitempty" validate:"omitempty,allowed_without=DeleteFloatings,dive,uuid4" delimiter:"comma"`
	ReservedFixedIPs []string `url:"reserved_fixed_ips,omitempty" validate:"omitempty,dive,uuid4" delimiter:"comma"`
}

// InstanceListOptions specifies the optional query parameters to List method.
type InstanceListOptions struct {
	IncludeBaremetal  bool   `url:"include_baremetal,omitempty"  validate:"omitempty"`
	IncludeK8S        bool   `url:"include_k8s,omitempty"  validate:"omitempty"`
	ExcludeSecgroup   string `url:"exclude_secgroup,omitempty"  validate:"omitempty"`
	AvailableFloating string `url:"available_floating,omitempty"  validate:"omitempty"`
	Name              string `url:"name,omitempty"  validate:"omitempty"`
	FlavorID          string `url:"flavor_id,omitempty"  validate:"omitempty"`
	Limit             int    `url:"limit,omitempty"  validate:"omitempty"`
	Offset            int    `url:"offset,omitempty"  validate:"omitempty"`
	Status            string `url:"status,omitempty"  validate:"omitempty"`
	ChangesSince      string `url:"changes-since,omitempty"  validate:"omitempty"`
	ChangesBefore     string `url:"changes-before,omitempty"  validate:"omitempty"`
	IP                string `url:"ip,omitempty"  validate:"omitempty"`
	UUID              string `url:"uuid,omitempty"  validate:"omitempty"`
	MetadataKV        string `url:"metadata_kv,omitempty"  validate:"omitempty"`
	MetadataK         string `url:"metadata_k,omitempty"  validate:"omitempty"`
	OrderBy           string `url:"order_by,omitempty"  validate:"omitempty"`
}

// instancesRoot represents an Instance root.
type instancesRoot struct {
	Count     int
	Instances []Instance `json:"results"`
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

// InstanceCheckFlavorVolumeRequest represents a request to get flavors of the instance.
type InstanceCheckFlavorVolumeRequest struct {
	Volumes []InstanceVolumeCreate `json:"volumes" required:"true" validate:"required,dive"`
}

type InstanceAvailableNames struct {
	AllowedBMNameWinTemplates []string `json:"allowed_bm_name_win_templates"`
	NameTemplatesLimited      bool     `json:"name_templates_limited"`
	AllowedNameWinTemplates   []string `json:"allowed_name_win_templates"`
	AllowedBMNameTemplates    []string `json:"allowed_bm_name_templates"`
	CustomNameAllowed         bool     `json:"custom_name_allowed"`
	AllowedNameTemplates      []string `json:"allowed_name_templates"`
}

type InstancePort struct {
	Name           string   `json:"name"`
	ID             string   `json:"id"`
	SecurityGroups []IDName `json:"security_groups"`
}

// InstanceMetricsListRequest represents a request to get a Instance Metrics list.
type InstanceMetricsListRequest struct {
	TimeUnit     TimeUnit `json:"time_unit" required:"true" validate:"required,name"`
	TimeInterval int      `json:"time_interval" required:"true" validate:"required,name"`
}

// InstanceMetrics represents an EdgecenterCloud Instance metrics.
type InstanceMetrics struct {
	Disks             []DiskMetrics `json:"disks"`
	CPUUtil           int           `json:"cpu_util"`
	NetworkPpsIngress int           `json:"network_pps_ingress"`
	NetworkBpsIngress int           `json:"network_Bps_ingress"`
	NetworkPpsEgress  int           `json:"network_pps_egress"`
	NetworkBpsEgress  int           `json:"network_Bps_egress"`
	Time              string        `json:"time"`
	MemoryUtil        int           `json:"memory_util"`
}

type DiskMetrics struct {
	DiskName      string `json:"disk_name"`
	DiskIOpsRead  int    `json:"disk_iops_read"`
	DiskIOpsWrite int    `json:"disk_iops_write"`
	DiskBpsWrite  int    `json:"disk_Bps_write"`
	DiskBpsRead   int    `json:"disk_Bps_read"`
}

type AssignSecurityGroupRequest struct {
	Name                    string                    `json:"name"`
	PortsSecurityGroupNames []PortsSecurityGroupNames `json:"ports_security_group_names"`
}

type PortsSecurityGroupNames struct {
	SecurityGroupNames []string `json:"security_group_names"`
	PortID             string   `json:"port_id"`
}

type RemoteConsole struct {
	URL      string `json:"url"`
	Type     string `json:"type"`
	Protocol string `json:"protocol"`
}

// InstanceAttachInterfaceRequest represents a request to attach Interface to the Instance.
type InstanceAttachInterfaceRequest struct {
	Type           InterfaceType `json:"type"`
	SecurityGroups []ID          `json:"security_groups"`
	SubnetID       string        `json:"subnet_id,omitempty" validate:"rfe=Type:subnet,omitempty,uuid4"`
	NetworkID      string        `json:"network_id,omitempty" validate:"rfe=Type:any_subnet,omitempty,uuid4"`
	PortID         string        `json:"port_id,omitempty" validate:"rfe=Type:reserved_fixed_ip,allowed_without_all=NetworkID SubnetID,omitempty,uuid4"`
}

// InstanceDetachInterfaceRequest represents a request to detach Interface from the Instance.
type InstanceDetachInterfaceRequest struct {
	PortID    string `json:"port_id,omitempty" validate:"omitempty"`
	IPAddress string `json:"ip_address,omitempty" validate:"omitempty"`
}

// InstancePutIntoServerGroupRequest represents a request to put an Interface into the Server Group.
type InstancePutIntoServerGroupRequest struct {
	ServerGroupID string `json:"servergroup_id"`
}

// instanceFlavorsRoot represents an Instance Flavors root.
type instanceFlavorsRoot struct {
	Count   int
	Flavors []Flavor `json:"results"`
}

// instancePortsRoot represents an Instance Flavors root.
type instancePortsRoot struct {
	Count         int
	InstancePorts []InstancePort `json:"results"`
}

// instanceMetricsRoot represents an Instance Metrics root.
type instanceMetricsRoot struct {
	Count           int
	InstanceMetrics []InstanceMetrics `json:"results"`
}

// instanceRouterInterfaceRoot represents an Instance RouterInterface root.
type instanceRouterInterfaceRoot struct {
	Count            int
	RouterInterfaces []RouterInterface `json:"results"`
}

// List get instances.
func (s *InstancesServiceOp) List(ctx context.Context, opts *InstanceListOptions) ([]Instance, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(instancesBasePathV1)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
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

// Get individual Instance.
func (s *InstancesServiceOp) Get(ctx context.Context, instanceID string) (*Instance, *Response, error) {
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(instancesBasePathV1), instanceID)

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
func (s *InstancesServiceOp) Create(ctx context.Context, reqBody *InstanceCreateRequest) (*TaskResponse, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(instancesBasePathV2)

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

// Delete the Instance.
func (s *InstancesServiceOp) Delete(ctx context.Context, instanceID string, opts *InstanceDeleteOptions) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(instancesBasePathV1), instanceID)
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

	path := s.client.addProjectRegionPath(instancesBasePathV1)
	path = fmt.Sprintf("%s/%s/%s", path, instanceID, metadataPath)

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

// MetadataList load balancer detailed metadata items.
func (s *InstancesServiceOp) MetadataList(ctx context.Context, instanceID string) ([]MetadataDetailed, *Response, error) {
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	return metadataList(ctx, s.client, instanceID, instancesBasePathV1)
}

// MetadataUpdate load balancer metadata.
func (s *InstancesServiceOp) MetadataUpdate(ctx context.Context, instanceID string, reqBody *Metadata) (*Response, error) {
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	return metadataUpdate(ctx, s.client, instanceID, instancesBasePathV1, reqBody)
}

// MetadataDeleteItem a load balancer metadata item by key.
func (s *InstancesServiceOp) MetadataDeleteItem(ctx context.Context, instanceID string, opts *MetadataItemOptions) (*Response, error) {
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	return metadataDeleteItem(ctx, s.client, instanceID, instancesBasePathV2, opts)
}

// MetadataGetItem load balancer detailed metadata.
func (s *InstancesServiceOp) MetadataGetItem(ctx context.Context, instanceID string, opts *MetadataItemOptions) (*MetadataDetailed, *Response, error) {
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	return metadataGetItem(ctx, s.client, instanceID, instancesBasePathV2, opts)
}

// MetadataCreate instance metadata (tags).
func (s *InstancesServiceOp) MetadataCreate(ctx context.Context, instanceID string, metadata *Metadata) (*Response, error) {
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	path := s.client.addProjectRegionPath(instancesBasePathV1)
	path = fmt.Sprintf("%s/%s/%s", path, instanceID, metadataPath)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, metadata)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// CheckLimits check a quota for instance creation.
func (s *InstancesServiceOp) CheckLimits(ctx context.Context, reqBody *InstanceCheckLimitsRequest) (*map[string]int, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(instancesBasePathV2)
	path = fmt.Sprintf("%s/%s", path, instancesCheckLimits)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
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
func (s *InstancesServiceOp) UpdateFlavor(ctx context.Context, instanceID string, reqBody *InstanceFlavorUpdateRequest) (*TaskResponse, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(instancesBasePathV1)
	path = fmt.Sprintf("%s/%s/%s", path, instanceID, instancesChangeFlavor)

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

// AvailableFlavors get flavors for an instance by volume config.
func (s *InstancesServiceOp) AvailableFlavors(ctx context.Context, reqBody *InstanceCheckFlavorVolumeRequest, opts *FlavorsOptions) ([]Flavor, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path, err := addOptions(s.client.addProjectRegionPath(instancesBasePathV1), opts)
	if err != nil {
		return nil, nil, err
	}

	path = fmt.Sprintf("%s/%s", path, instancesAvailableFlavors)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	flavors := new(instanceFlavorsRoot)
	resp, err := s.client.Do(ctx, req, flavors)
	if err != nil {
		return nil, resp, err
	}

	return flavors.Flavors, resp, err
}

// AvailableFlavorsToResize Get flavors to resize into.
func (s *InstancesServiceOp) AvailableFlavorsToResize(ctx context.Context, instanceID string, opts *FlavorsOptions) ([]Flavor, *Response, error) {
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(instancesBasePathV1), instanceID, instancesAvailableFlavors)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	flavors := new(instanceFlavorsRoot)
	resp, err := s.client.Do(ctx, req, flavors)
	if err != nil {
		return nil, resp, err
	}

	return flavors.Flavors, resp, err
}

// AvailableNames get instance naming restrictions that are applied to specified project and region.
func (s *InstancesServiceOp) AvailableNames(ctx context.Context) (*InstanceAvailableNames, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(instancesBasePathV1), instancesAvailableNames)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	instanceAvailableNames := new(InstanceAvailableNames)
	resp, err := s.client.Do(ctx, req, instanceAvailableNames)
	if err != nil {
		return nil, resp, err
	}

	return instanceAvailableNames, resp, err
}

// Rename the Instance.
func (s *InstancesServiceOp) Rename(ctx context.Context, instanceID string, reqBody *Name) (*Instance, *Response, error) {
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return nil, resp, err
	}

	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(instancesBasePathV1), instanceID)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, reqBody)
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

// PortsList get network ports.
func (s *InstancesServiceOp) PortsList(ctx context.Context, instanceID string) ([]InstancePort, *Response, error) {
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(instancesBasePathV1), instanceID, instancesPorts)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(instancePortsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.InstancePorts, resp, err
}

// InstanceStart start the instance.
func (s *InstancesServiceOp) InstanceStart(ctx context.Context, instanceID string) (*Instance, *Response, error) {
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(instancesBasePathV1), instanceID, instancesStart)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
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

// InstanceStop stop the instance.
func (s *InstancesServiceOp) InstanceStop(ctx context.Context, instanceID string) (*Instance, *Response, error) {
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(instancesBasePathV1), instanceID, instancesStop)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
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

// InstancePowercycle powercycle the instance.
func (s *InstancesServiceOp) InstancePowercycle(ctx context.Context, instanceID string) (*Instance, *Response, error) {
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(instancesBasePathV1), instanceID, instancesPowercycle)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
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

// InstanceReboot reboot the instance.
func (s *InstancesServiceOp) InstanceReboot(ctx context.Context, instanceID string) (*Instance, *Response, error) {
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(instancesBasePathV1), instanceID, instancesReboot)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
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

// InstanceSuspend suspend the instance.
func (s *InstancesServiceOp) InstanceSuspend(ctx context.Context, instanceID string) (*Instance, *Response, error) {
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(instancesBasePathV1), instanceID, instancesSuspend)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
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

// InstanceResume resume the instance.
func (s *InstancesServiceOp) InstanceResume(ctx context.Context, instanceID string) (*Instance, *Response, error) {
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(instancesBasePathV1), instanceID, instancesResume)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
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

// MetricsList get instance metrics.
func (s *InstancesServiceOp) MetricsList(ctx context.Context, instanceID string, reqBody *InstanceMetricsListRequest) ([]InstanceMetrics, *Response, error) {
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(instancesBasePathV1), instanceID, instancesMetrics)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	root := new(instanceMetricsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.InstanceMetrics, resp, err
}

// FilterBySecurityGroup returns a list of instances with the filter by the security group.
func (s *InstancesServiceOp) FilterBySecurityGroup(ctx context.Context, securityGroupID string) ([]Instance, *Response, error) {
	if resp, err := isValidUUID(securityGroupID, "securityGroupID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(instancesBasePathV1), securityGroupID, instancesInstances)

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

// SecurityGroupList returns a list of instance security groups.
func (s *InstancesServiceOp) SecurityGroupList(ctx context.Context, instanceID string) ([]IDName, *Response, error) {
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(instancesBasePathV1), instanceID, instancesSecurityGroups)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(idNameRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.IDNames, resp, err
}

// SecurityGroupAssign the security group to the server.
func (s *InstancesServiceOp) SecurityGroupAssign(ctx context.Context, instanceID string, reqBody *AssignSecurityGroupRequest) (*Response, error) {
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return resp, err
	}

	if reqBody == nil {
		return nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(instancesBasePathV1), instanceID, instancesAddSecurityGroup)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// SecurityGroupUnAssign the security group to the server.
func (s *InstancesServiceOp) SecurityGroupUnAssign(ctx context.Context, instanceID string, reqBody *AssignSecurityGroupRequest) (*Response, error) {
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return resp, err
	}

	if reqBody == nil {
		return nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(instancesBasePathV1), instanceID, instancesDelSecurityGroup)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// GetConsole get an Instance console URL.
func (s *InstancesServiceOp) GetConsole(ctx context.Context, instanceID string) (*RemoteConsole, *Response, error) {
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(instancesBasePathV1), instanceID, instancesGetConsole)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	remoteConsole := new(RemoteConsole)
	resp, err := s.client.Do(ctx, req, remoteConsole)
	if err != nil {
		return nil, resp, err
	}

	return remoteConsole, resp, err
}

// AttachInterface to the instance.
func (s *InstancesServiceOp) AttachInterface(ctx context.Context, instanceID string, reqBody *InstanceAttachInterfaceRequest) (*TaskResponse, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(instancesBasePathV1)
	path = fmt.Sprintf("%s/%s/%s", path, instanceID, instancesAttachInterface)

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

// DetachInterface from the instance.
func (s *InstancesServiceOp) DetachInterface(ctx context.Context, instanceID string, reqBody *InstanceDetachInterfaceRequest) (*TaskResponse, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(instancesBasePathV1)
	path = fmt.Sprintf("%s/%s/%s", path, instanceID, instancesDetachInterface)

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

// InterfaceList returns a list of network interfaces attached to the Instance.
func (s *InstancesServiceOp) InterfaceList(ctx context.Context, instanceID string) ([]RouterInterface, *Response, error) {
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(instancesBasePathV1), instanceID, instancesInterfaces)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(instanceRouterInterfaceRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.RouterInterfaces, resp, err
}

// PutIntoServerGroup put an instance into server group.
func (s *InstancesServiceOp) PutIntoServerGroup(ctx context.Context, instanceID string, reqBody *InstancePutIntoServerGroupRequest) (*TaskResponse, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(instancesBasePathV1)
	path = fmt.Sprintf("%s/%s/%s", path, instanceID, instancesPutIntoServerGroup)

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

// RemoveFromServerGroup remove an instance from server group.
func (s *InstancesServiceOp) RemoveFromServerGroup(ctx context.Context, instanceID string) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(instanceID, "instanceID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(instancesBasePathV1)
	path = fmt.Sprintf("%s/%s/%s", path, instanceID, instancesRemoveFromServerGroup)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
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
