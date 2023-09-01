package pools

import (
	"net/http"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/instance/v1/instances"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/task/v1/tasks"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/volume/v1/volumes"
	"github.com/Edge-Center/edgecentercloud-go/pagination"
)

// Get retrieves a specific cluster pool based on its unique ID.
func Get(client *edgecloud.ServiceClient, clusterID, poolID string) (r GetResult) {
	var response *http.Response
	response, r.Err = client.Get(getURL(client, clusterID, poolID), &r.Body, &edgecloud.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	if r.Err == nil {
		r.Header = response.Header
	}
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToClusterPoolsListQuery() (string, error)
}

// ListOpts is used to filter and sort the pool of a cluster when using List.
type ListOpts struct {
	// Pagination marker for large data sets. (UUID field from node group).
	Marker int `q:"marker"`
	// Maximum number of resources to return in a single page.
	Limit int `q:"limit"`
	// Column to sort results by. Default: id.
	SortKey string `q:"sort_key"`
	// Direction to sort. "asc" or "desc". Default: asc.
	SortDir string `q:"sort_dir"`
	// List all pools with the specified role.
	Role string `q:"role"`
	// Details
	Detail bool `q:"detail"`
}

// ToClusterPoolsListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToClusterPoolsListQuery() (string, error) {
	q, err := edgecloud.BuildQueryString(opts)
	return q.String(), err
}

// List makes a request to the Magnum API to retrieve pools
// belonging to the given cluster. The request can be modified to
// filter or sort the list using the options available in ListOpts.
//
// Use the AllPages method of the returned Pager to ensure that
// all pools are returned (for example when using the Limit option
// to limit the number of node groups returned per page).
//
// Not all node group fields are returned in a list request.
// Only the fields UUID, Name, FlavorID, ImageID,
// NodeCount, Role, IsDefault, Status and StackID
// are returned, all other fields are omitted
// and will have their zero value when extracted.
func List(client *edgecloud.ServiceClient, clusterID string, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client, clusterID)
	if opts != nil {
		query, err := opts.ToClusterPoolsListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ClusterPoolPage{LinkedPageBase: pagination.LinkedPageBase{PageResult: r}, ClusterID: &clusterID}
	})
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToClusterPoolCreateMap() (map[string]interface{}, error)
}

// CreateOpts is used to set available fields upon node group creation.
//
// If unset, some fields have defaults or will inherit from the cluster value.
type CreateOpts struct {
	Name             string             `json:"name" required:"true" validate:"required"`
	FlavorID         string             `json:"flavor_id" required:"true" validate:"required"`
	NodeCount        *int               `json:"node_count,omitempty" validate:"omitempty,gt=0"`
	DockerVolumeSize *int               `json:"docker_volume_size,omitempty" validate:"omitempty,gt=0"`
	DockerVolumeType volumes.VolumeType `json:"docker_volume_type,omitempty" validate:"omitempty,enum"`
	MinNodeCount     int                `json:"min_node_count,omitempty" validate:"omitempty,gt=0,ltefield=NodeCount"`
	MaxNodeCount     *int               `json:"max_node_count,omitempty" validate:"omitempty,gt=0,gtefield=MinNodeCount,gtefield=NodeCount"`
	Labels           map[string]string  `json:"labels,omitempty"`
	ImageID          string             `json:"image_id,omitempty"`
}

// Validate CreateOpts.
func (opts CreateOpts) Validate() error {
	return edgecloud.TranslateValidationError(edgecloud.Validate.Struct(opts))
}

// ToClusterPoolCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToClusterPoolCreateMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return edgecloud.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and creates a new cluster Pool using the values provided.
func Create(client *edgecloud.ServiceClient, clusterID string, opts CreateOptsBuilder) (r tasks.Result) {
	b, err := opts.ToClusterPoolCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	var result *http.Response
	result, r.Err = client.Post(createURL(client, clusterID), b, &r.Body, &edgecloud.RequestOpts{
		OkCodes: []int{http.StatusOK, http.StatusCreated},
	})

	if r.Err == nil {
		r.Header = result.Header
	}

	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToClusterPoolUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to update a pool.
type UpdateOpts struct {
	Name         string `json:"name,omitempty" validate:"required_without_all=MinNodeCount MaxNodeCount,omitempty"`
	MinNodeCount int    `json:"min_node_count,omitempty" validate:"required_without_all=Name MaxNodeCount,omitempty,gt=0"`
	MaxNodeCount int    `json:"max_node_count,omitempty" validate:"required_without_all=Name MixNodeCount,omitempty,gt=0,gtefield=MinNodeCount"`
}

// Validate UpdateOpts.
func (opts UpdateOpts) Validate() error {
	return edgecloud.TranslateValidationError(edgecloud.Validate.Struct(opts))
}

// ToClusterPoolUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToClusterPoolUpdateMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return edgecloud.BuildRequestBody(opts, "")
}

// Update accepts a UpdateOpts struct and updates an existing Pool using the values provided.
func Update(client *edgecloud.ServiceClient, clusterID, poolID string, opts UpdateOptsBuilder) (r tasks.Result) {
	b, err := opts.ToClusterPoolUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	var result *http.Response
	result, r.Err = client.Patch(updateURL(client, clusterID, poolID), b, &r.Body, &edgecloud.RequestOpts{
		OkCodes: []int{http.StatusOK, http.StatusCreated},
	})

	if r.Err == nil {
		r.Header = result.Header
	}

	return
}

// Delete accepts a unique ID and deletes the cluster Pool associated with it.
func Delete(client *edgecloud.ServiceClient, clusterID, nodeGroupID string) (r tasks.Result) {
	var result *http.Response
	result, r.Err = client.DeleteWithResponse(deleteURL(client, clusterID, nodeGroupID), &r.Body, nil)
	r.Header = result.Header
	return
}

// ListAll is a convenience function that returns all cluster pools.
func ListAll(client *edgecloud.ServiceClient, clusterID string, opts ListOptsBuilder) ([]ClusterPoolList, error) {
	pages, err := List(client, clusterID, opts).AllPages()
	if err != nil {
		return nil, err
	}

	return ExtractClusterPools(pages, &clusterID)
}

// Instances returns a Pager which allows you to iterate over a collection of pool instances.
func Instances(client *edgecloud.ServiceClient, clusterID string, id string) pagination.Pager {
	url := instancesURL(client, clusterID, id)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return instances.InstancePage{LinkedPageBase: pagination.LinkedPageBase{PageResult: r}}
	})
}

// InstancesAll returns all pool instances.
func InstancesAll(client *edgecloud.ServiceClient, clusterID, id string) ([]instances.Instance, error) {
	page, err := Instances(client, clusterID, id).AllPages()
	if err != nil {
		return nil, err
	}
	return instances.ExtractInstances(page)
}

// Volumes returns a Pager which allows you to iterate over a collection of pool instances.
func Volumes(client *edgecloud.ServiceClient, clusterID string, id string) pagination.Pager {
	url := volumesURL(client, clusterID, id)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return volumes.VolumePage{LinkedPageBase: pagination.LinkedPageBase{PageResult: r}}
	})
}

// VolumesAll returns all pool volumes.
func VolumesAll(client *edgecloud.ServiceClient, clusterID, id string) ([]volumes.Volume, error) {
	page, err := Volumes(client, clusterID, id).AllPages()
	if err != nil {
		return nil, err
	}
	return volumes.ExtractVolumes(page)
}
