package instances

import (
	"log"
	"net/http"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/flavor/v1/flavors"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/instance/v1/types"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/task/v1/tasks"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/volume/v1/volumes"
	"github.com/Edge-Center/edgecentercloud-go/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToInstanceListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct {
	ExcludeSecGroup   string            `q:"exclude_secgroup"`
	AvailableFloating bool              `q:"available_floating"`
	IncludeBaremetal  bool              `q:"include_baremetal"`
	Name              string            `q:"name"`
	FlavorID          string            `q:"flavor_id"`
	Limit             int               `q:"limit" validate:"omitempty,gt=0"`
	Offset            int               `q:"offset" validate:"omitempty,gt=0"`
	Metadata          map[string]string `q:"metadata_kv" validate:"omitempty"`
}

// ToInstanceListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToInstanceListQuery() (string, error) {
	if err := edgecloud.ValidateStruct(opts); err != nil {
		return "", err
	}
	q, err := edgecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// DeleteOptsBuilder allows extensions to add additional parameters to the Delete request.
type DeleteOptsBuilder interface {
	ToInstanceDeleteQuery() (string, error)
}

// DeleteOpts Set parameters for delete operation.
type DeleteOpts struct {
	Volumes         []string `q:"volumes" validate:"omitempty,dive,uuid4" delimiter:"comma"`
	DeleteFloatings bool     `q:"delete_floatings" validate:"omitempty,allowed_without=FloatingIPs"`
	FloatingIPs     []string `q:"floatings" validate:"omitempty,allowed_without=DeleteFloatings,dive,uuid4" delimiter:"comma"`
}

// ToInstanceDeleteQuery formats a DeleteOpts into a query string.
func (opts DeleteOpts) ToInstanceDeleteQuery() (string, error) {
	if err := opts.Validate(); err != nil {
		return "", err
	}
	q, err := edgecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

func (opts *DeleteOpts) Validate() error {
	return edgecloud.ValidateStruct(opts)
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToInstanceCreateMap() (map[string]interface{}, error)
}

// CreateVolumeOpts represents options used to create a volume.
type CreateVolumeOpts struct {
	Source        types.VolumeSource `json:"source" required:"true" validate:"required,enum"`
	BootIndex     int                `json:"boot_index"`
	Size          int                `json:"size,omitempty" validate:"rfe=Source:image;new-volume,sfe=Source:snapshot;existing-volume"`
	TypeName      volumes.VolumeType `json:"type_name,omitempty" validate:"omitempty"`
	AttachmentTag string             `json:"attachment_tag,omitempty" validate:"omitempty"`
	Name          string             `json:"name,omitempty" validate:"omitempty"`
	ImageID       string             `json:"image_id,omitempty" validate:"rfe=Source:image,sfe=Source:snapshot;existing-volume;new-volume,allowed_without_all=SnapshotID VolumeID,omitempty,uuid4"`
	SnapshotID    string             `json:"snapshot_id,omitempty" validate:"rfe=Source:snapshot,sfe=Source:image;existing-volume;new-volume,allowed_without_all=ImageID VolumeID,omitempty,uuid4"`
	VolumeID      string             `json:"volume_id,omitempty" validate:"rfe=Source:existing-volume,sfe=Source:image;shapshot;new-volume,allowed_without_all=ImageID SnapshotID,omitempty,uuid4"`
}

func (opts *CreateVolumeOpts) Validate() error {
	return edgecloud.ValidateStruct(opts)
}

type CreateNewInterfaceFloatingIPOpts struct {
	Source             types.FloatingIPSource `json:"source" validate:"required,enum"`
	ExistingFloatingID string                 `json:"existing_floating_id" validate:"rfe=Source:existing,sfe=Source:new,omitempty,uuid"`
}

// Validate CreateNewInterfaceFloatingIPOpts.
func (opts CreateNewInterfaceFloatingIPOpts) Validate() error {
	return edgecloud.ValidateStruct(opts)
}

type InterfaceOpts struct {
	Type       types.InterfaceType               `json:"type,omitempty" validate:"omitempty,enum"`
	NetworkID  string                            `json:"network_id,omitempty" validate:"rfe=Type:any_subnet,omitempty,uuid4"`
	SubnetID   string                            `json:"subnet_id,omitempty" validate:"rfe=Type:subnet,omitempty,uuid4"`
	PortID     string                            `json:"port_id,omitempty" validate:"rfe=Type:reserved_fixed_ip,allowed_without_all=NetworkID SubnetID,omitempty,uuid4"`
	IpAddress  string                            `json:"ip_address,omitempty" validate:"allowed_without_all=Type NetworkID SubnetID FloatingIP,omitempty"`
	FloatingIP *CreateNewInterfaceFloatingIPOpts `json:"floating_ip,omitempty" validate:"omitempty,dive"`
}

// Validate InterfaceOpts.
func (opts InterfaceOpts) Validate() error {
	return edgecloud.ValidateStruct(opts)
}

type InterfaceInstanceCreateOpts struct {
	InterfaceOpts
	SecurityGroups []edgecloud.ItemID `json:"security_groups"`
}

// ToInterfaceActionMap builds a request body from InterfaceOpts.
func (opts InterfaceInstanceCreateOpts) ToInterfaceActionMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	result, err := edgecloud.BuildRequestBody(opts, "")
	log.Println(result)
	return result, err
}

// Validate InterfaceInstanceCreateOpts.
func (opts InterfaceInstanceCreateOpts) Validate() error {
	return edgecloud.ValidateStruct(opts)
}

// CreateOpts represents options used to create a instance.
type CreateOpts struct {
	Flavor         string                        `json:"flavor" required:"true"`
	Names          []string                      `json:"names,omitempty" validate:"required_without=NameTemplates"`
	NameTemplates  []string                      `json:"name_templates,omitempty" validate:"required_without=Names"`
	Volumes        []CreateVolumeOpts            `json:"volumes" required:"true" validate:"required,dive"`
	Interfaces     []InterfaceInstanceCreateOpts `json:"interfaces" required:"true" validate:"required,dive"`
	SecurityGroups []edgecloud.ItemID            `json:"security_groups,omitempty" validate:"omitempty,dive,uuid4"`
	Keypair        string                        `json:"keypair_name,omitempty"`
	Password       string                        `json:"password" validate:"omitempty,required_with=Username"`
	Username       string                        `json:"username" validate:"omitempty,required_with=Password"`
	UserData       string                        `json:"user_data" validate:"omitempty,base64"`
	Metadata       *MetadataSetOpts              `json:"metadata,omitempty" validate:"omitempty,dive"`
	Configuration  *MetadataSetOpts              `json:"configuration,omitempty" validate:"omitempty,dive"`
	AllowAppPorts  bool                          `json:"allow_app_ports,omitempty"`
	ServerGroupID  string                        `json:"servergroup_id,omitempty" validate:"omitempty,uuid4"`
}

// Validate CreateOpts.
func (opts CreateOpts) Validate() error {
	return edgecloud.ValidateStruct(opts)
}

// ToInstanceCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToInstanceCreateMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	var err error
	var metadata map[string]interface{}
	var configuration map[string]interface{}
	if opts.Metadata != nil {
		metadata, err = opts.Metadata.ToMetadataMap()
		if err != nil {
			return nil, err
		}
	}
	if opts.Configuration != nil {
		configuration, err = opts.Configuration.ToMetadataMap()
		if err != nil {
			return nil, err
		}
	}
	mp, err := edgecloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	if len(metadata) > 0 {
		mp["metadata"] = metadata
	} else {
		delete(mp, "metadata")
	}
	if len(configuration) > 0 {
		mp["configuration"] = configuration
	} else {
		delete(mp, "configuration")
	}
	return mp, nil
}

// RenameInstanceOptsBuilder allows extensions to add parameters to rename instance request.
type RenameInstanceOptsBuilder interface {
	ToRenameInstanceActionMap() (map[string]interface{}, error)
}

type RenameInstanceOpts struct {
	Name string `json:"name" required:"true" validate:"required"`
}

// Validate RenameInstanceOpts.
func (opts RenameInstanceOpts) Validate() error {
	return edgecloud.ValidateStruct(opts)
}

// ToRenameInstanceActionMap builds a request body from RenameInstanceOpts.
func (opts RenameInstanceOpts) ToRenameInstanceActionMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return edgecloud.BuildRequestBody(opts, "")
}

// SecurityGroupOptsBuilder allows extensions to add parameters to the security groups request.
type SecurityGroupOptsBuilder interface {
	ToSecurityGroupActionMap() (map[string]interface{}, error)
}

type SecurityGroupOpts struct {
	Name                    string                   `json:"name,omitempty"`
	PortsSecurityGroupNames []PortSecurityGroupNames `json:"ports_security_group_names,omitempty"`
}

type PortSecurityGroupNames struct {
	PortID             *string  `json:"port_id"`
	SecurityGroupNames []string `json:"security_group_names"`
}

// Validate SecurityGroupOpts.
func (opts SecurityGroupOpts) Validate() error {
	return edgecloud.ValidateStruct(opts)
}

// ToSecurityGroupActionMap builds a request body from SecurityGroupOpts.
func (opts SecurityGroupOpts) ToSecurityGroupActionMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return edgecloud.BuildRequestBody(opts, "")
}

// InterfaceOptsBuilder allows extensions to add parameters to the interface request.
type InterfaceOptsBuilder interface {
	ToInterfaceActionMap() (map[string]interface{}, error)
}

// ToInterfaceActionMap builds a request body from InterfaceOpts.
func (opts InterfaceOpts) ToInterfaceActionMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return edgecloud.BuildRequestBody(opts, "")
}

// ChangeFlavorOptsBuilder builds parameters or change flavor request.
type ChangeFlavorOptsBuilder interface {
	ToChangeFlavorActionMap() (map[string]interface{}, error)
}

type ChangeFlavorOpts struct {
	FlavorID string `json:"flavor_id" required:"true" validate:"required"`
}

func (opts ChangeFlavorOpts) Validate() error {
	return edgecloud.ValidateStruct(opts)
}

// ToChangeFlavorActionMap builds a request body from ChangeFlavorOpts.
func (opts ChangeFlavorOpts) ToChangeFlavorActionMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return edgecloud.BuildRequestBody(opts, "")
}

// ServerGroupOptsBuilder builds parameters or change server group request.
type ServerGroupOptsBuilder interface {
	ToServerGroupActionMap() (map[string]interface{}, error)
}

type ServerGroupOpts struct {
	ServerGroupID string `json:"servergroup_id" required:"true" validate:"required"`
}

func (opts ServerGroupOpts) Validate() error {
	return edgecloud.ValidateStruct(opts)
}

// ToServerGroupActionMap builds a request body from ServerGroupOpts.
func (opts ServerGroupOpts) ToServerGroupActionMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return edgecloud.BuildRequestBody(opts, "")
}

func List(client *edgecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToInstanceListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return InstancePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific instance based on its unique ID.
func Get(client *edgecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(client, id)
	_, r.Err = client.Get(url, &r.Body, nil)
	return
}

// ListInterfaces retrieves network interfaces for instance.
func ListInterfaces(client *edgecloud.ServiceClient, id string) pagination.Pager {
	url := interfacesListURL(client, id)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return InstanceInterfacePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAll is a convenience function that returns all instances.
func ListAll(client *edgecloud.ServiceClient, opts ListOptsBuilder) ([]Instance, error) {
	pages, err := List(client, opts).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractInstances(pages)
	if err != nil {
		return nil, err
	}

	return all, nil
}

// ListInterfacesAll is a convenience function that returns all instance interfaces.
func ListInterfacesAll(client *edgecloud.ServiceClient, id string) ([]Interface, error) {
	pages, err := ListInterfaces(client, id).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractInstanceInterfaces(pages)
	if err != nil {
		return nil, err
	}

	return all, nil
}

// ListSecurityGroups retrieves security groups interfaces for instance.
func ListSecurityGroups(client *edgecloud.ServiceClient, id string) pagination.Pager {
	url := securityGroupsListURL(client, id)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return InstanceSecurityGroupPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListSecurityGroupsAll is a convenience function that returns all instance security groups.
func ListSecurityGroupsAll(client *edgecloud.ServiceClient, id string) ([]edgecloud.ItemIDName, error) {
	pages, err := ListSecurityGroups(client, id).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractInstanceSecurityGroups(pages)
	if err != nil {
		return nil, err
	}

	return all, nil
}

// ListPorts retrieves ports for instance.
func ListPorts(client *edgecloud.ServiceClient, id string) pagination.Pager {
	url := portsListURL(client, id)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return InstancePortsPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListPortsAll is a convenience function that returns all instance ports.
func ListPortsAll(client *edgecloud.ServiceClient, id string) ([]InstancePorts, error) {
	pages, err := ListPorts(client, id).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractInstancePorts(pages)
	if err != nil {
		return nil, err
	}

	return all, nil
}

// RenameInstance rename instance.
func RenameInstance(client *edgecloud.ServiceClient, id string, opts RenameInstanceOptsBuilder) (r GetResult) {
	b, err := opts.ToRenameInstanceActionMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Patch(renameInstanceURL(client, id), b, &r.Body, &edgecloud.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	return
}

// AssignSecurityGroup adds a security groups to the instance.
func AssignSecurityGroup(client *edgecloud.ServiceClient, id string, opts SecurityGroupOptsBuilder) (r SecurityGroupActionResult) {
	b, err := opts.ToSecurityGroupActionMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(addSecurityGroupsURL(client, id), b, nil, &edgecloud.RequestOpts{
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}

// UnAssignSecurityGroup removes a security groups from the instance.
func UnAssignSecurityGroup(client *edgecloud.ServiceClient, id string, opts SecurityGroupOptsBuilder) (r SecurityGroupActionResult) {
	b, err := opts.ToSecurityGroupActionMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(deleteSecurityGroupsURL(client, id), b, nil, &edgecloud.RequestOpts{
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}

// AddServerGroup adds a server group to the instance.
func AddServerGroup(client *edgecloud.ServiceClient, id string, opts ServerGroupOptsBuilder) (r tasks.Result) {
	b, err := opts.ToServerGroupActionMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(putIntoServerGroupURL(client, id), b, &r.Body, &edgecloud.RequestOpts{
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}

// RemoveServerGroup removes a server group from the instance.
func RemoveServerGroup(client *edgecloud.ServiceClient, id string) (r tasks.Result) {
	_, r.Err = client.Post(removeFromServerGroupURL(client, id), nil, &r.Body, &edgecloud.RequestOpts{
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}

// AttachInterface adds a interface to the instance.
func AttachInterface(client *edgecloud.ServiceClient, id string, opts InterfaceOptsBuilder) (r tasks.Result) {
	b, err := opts.ToInterfaceActionMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(attachInterfaceURL(client, id), b, &r.Body, &edgecloud.RequestOpts{
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}

// DetachInterface removes a interface from the instance.
func DetachInterface(client *edgecloud.ServiceClient, id string, opts InterfaceOptsBuilder) (r tasks.Result) {
	b, err := opts.ToInterfaceActionMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(detachInterfaceURL(client, id), b, &r.Body, &edgecloud.RequestOpts{
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}

// Create creates an instance.
func Create(client *edgecloud.ServiceClient, opts CreateOptsBuilder) (r tasks.Result) {
	b, err := opts.ToInstanceCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), b, &r.Body, nil)
	return
}

func Delete(client *edgecloud.ServiceClient, instanceID string, opts DeleteOptsBuilder) (r tasks.Result) {
	url := deleteURL(client, instanceID)
	if opts != nil {
		query, err := opts.ToInstanceDeleteQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	_, r.Err = client.DeleteWithResponse(url, &r.Body, nil)
	return
}

// Start instance.
func Start(client *edgecloud.ServiceClient, id string) (r UpdateResult) {
	_, r.Err = client.Post(startInstanceURL(client, id), nil, &r.Body, nil)
	return
}

// Stop instance.
func Stop(client *edgecloud.ServiceClient, id string) (r UpdateResult) {
	_, r.Err = client.Post(stopInstanceURL(client, id), nil, &r.Body, nil)
	return
}

// PowerCycle instance.
func PowerCycle(client *edgecloud.ServiceClient, id string) (r UpdateResult) {
	_, r.Err = client.Post(powerCycleInstanceURL(client, id), nil, &r.Body, nil)
	return
}

// Reboot instance.
func Reboot(client *edgecloud.ServiceClient, id string) (r UpdateResult) {
	_, r.Err = client.Post(rebootInstanceURL(client, id), nil, &r.Body, nil)
	return
}

// Suspend instance.
func Suspend(client *edgecloud.ServiceClient, id string) (r UpdateResult) {
	_, r.Err = client.Post(suspendInstanceURL(client, id), nil, &r.Body, nil)
	return
}

// Resume instance.
func Resume(client *edgecloud.ServiceClient, id string) (r UpdateResult) {
	_, r.Err = client.Post(resumeInstanceURL(client, id), nil, &r.Body, nil)
	return
}

// Resize instance.
func Resize(client *edgecloud.ServiceClient, id string, opts ChangeFlavorOptsBuilder) (r tasks.Result) {
	b, err := opts.ToChangeFlavorActionMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(changeFlavorInstanceURL(client, id), b, &r.Body, nil)
	return
}

// ListMetricsOptsBuilder builds parameters or change flavor request.
type ListMetricsOptsBuilder interface {
	ToListMetricsMap() (map[string]interface{}, error)
}

type ListMetricsOpts struct {
	TimeUnit     types.MetricsTimeUnit `json:"time_unit" required:"true" validate:"required,enum"`
	TimeInterval int                   `json:"time_interval" required:"true" validate:"required"`
}

func (opts ListMetricsOpts) Validate() error {
	return edgecloud.ValidateStruct(opts)
}

// ToListMetricsMap builds a request body from ListMetricsOpts.
func (opts ListMetricsOpts) ToListMetricsMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return edgecloud.BuildRequestBody(opts, "")
}

// ListInstanceMetrics retrieves instance's metrics.
func ListInstanceMetrics(client *edgecloud.ServiceClient, id string, opts ListMetricsOptsBuilder) (r ListMetricsResult) {
	b, err := opts.ToListMetricsMap()
	if err != nil {
		return
	}
	_, r.Err = client.Post(listInstanceMetricsURL(client, id), b, &r.Body, nil)
	return
}

func MetadataList(client *edgecloud.ServiceClient, id string) pagination.Pager {
	url := metadataURL(client, id)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return MetadataPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

func MetadataListAll(client *edgecloud.ServiceClient, id string) ([]Metadata, error) {
	pages, err := MetadataList(client, id).AllPages()
	if err != nil {
		return nil, err
	}
	all, err := ExtractMetadata(pages)
	if err != nil {
		return nil, err
	}
	return all, nil
}

// MetadataOptsBuilder allows extensions to add additional parameters to the metadata Create and Update request.
type MetadataOptsBuilder interface {
	ToMetadataMap() (string, error)
}

// MetadataOpts Set parameters for Create or Update operation.
type MetadataOpts struct {
	Key   string `json:"key" validate:"required,max=255"`
	Value string `json:"value" validate:"required,max=255"`
}

// MetadataSetOpts Set parameters for Create or Update operation.
type MetadataSetOpts struct {
	Metadata []MetadataOpts `json:"metadata" validate:"required,min=1,dive"`
}

// Validate MetadataOpts.
func (opts MetadataOpts) Validate() error {
	return edgecloud.ValidateStruct(opts)
}

// Validate MetadataSetOpts.
func (opts MetadataSetOpts) Validate() error {
	return edgecloud.ValidateStruct(opts)
}

// ToMetadataMap builds a request body from MetadataSetOpts.
func (opts MetadataSetOpts) ToMetadataMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	for _, md := range opts.Metadata {
		m[md.Key] = md.Value
	}
	return m, nil
}

// MetadataCreate creates a metadata for an instance.
func MetadataCreate(client *edgecloud.ServiceClient, id string, opts MetadataSetOpts) (r MetadataActionResult) {
	b, err := opts.ToMetadataMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(metadataURL(client, id), b, nil, &edgecloud.RequestOpts{
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}

// MetadataUpdate updates a metadata for an instance.
func MetadataUpdate(client *edgecloud.ServiceClient, id string, opts MetadataSetOpts) (r MetadataActionResult) {
	b, err := opts.ToMetadataMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(metadataURL(client, id), b, nil, &edgecloud.RequestOpts{
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}

// MetadataDelete deletes defined metadata key for an instance.
func MetadataDelete(client *edgecloud.ServiceClient, id string, key string) (r MetadataActionResult) {
	_, r.Err = client.Delete(metadataDetailsURL(client, id, key), &edgecloud.RequestOpts{
		OkCodes: []int{http.StatusNoContent, http.StatusOK},
	})
	return
}

// MetadataGet gets defined metadata key for an instance.
func MetadataGet(client *edgecloud.ServiceClient, id string, key string) (r MetadataResult) {
	url := metadataDetailsURL(client, id, key)
	_, r.Err = client.Get(url, &r.Body, nil)
	return
}

// ListAvailableFlavors get available flavors for the instance to resize into.
func ListAvailableFlavors(client *edgecloud.ServiceClient, id string, opts flavors.ListOptsBuilder) (r flavors.ListResult) {
	url := listAvailableFlavorsURL(client, id)
	if opts != nil {
		query, err := opts.ToFlavorListQuery()
		if err != nil {
			return
		}
		url += query
	}
	_, r.Err = client.Post(url, nil, &r.Body, nil)
	return
}

// GetSpiceConsole retrieves a specific spice console based on instance unique ID.
func GetSpiceConsole(client *edgecloud.ServiceClient, id string) (r RemoteConsoleResult) {
	url := getSpiceConsoleURL(client, id)
	_, r.Err = client.Get(url, &r.Body, nil)
	return
}

// GetInstanceConsole retrieves a specific spice console based on instance unique ID.
func GetInstanceConsole(client *edgecloud.ServiceClient, id string) (r RemoteConsoleResult) {
	url := getInstanceConsoleURL(client, id)
	_, r.Err = client.Get(url, &r.Body, nil)
	return
}

// ListInstanceLocationOptsBuilder allows extensions to add additional parameters to the ListInstanceLocation request.
type ListInstanceLocationOptsBuilder interface {
	ToListInstanceLocationQuery() (string, error)
}

// ListInstanceLocationOpts set parameters for search instance location operation.
type ListInstanceLocationOpts struct {
	Name string `q:"name"`
	ID   string `q:"id"`
}

// ToListInstanceLocationQuery formats a ListInstanceLocationOpts into a query string.
func (opts ListInstanceLocationOpts) ToListInstanceLocationQuery() (string, error) {
	if err := opts.Validate(); err != nil {
		return "", err
	}
	q, err := edgecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// Validate ListInstanceLocationOpts.
func (opts *ListInstanceLocationOpts) Validate() error {
	return edgecloud.ValidateStruct(opts)
}

// ListInstanceLocation get flavors available for the instance to resize into.
func ListInstanceLocation(client *edgecloud.ServiceClient, opts ListInstanceLocationOptsBuilder) (r SearchLocationResult) {
	url := listInstanceLocationURL(client)
	if opts != nil {
		query, err := opts.ToListInstanceLocationQuery()
		if err != nil {
			return
		}
		url += query
	}
	_, r.Err = client.Get(url, &r.Body, nil)
	return
}
