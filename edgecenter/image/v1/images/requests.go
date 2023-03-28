package images

import (
	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/image/v1/images/types"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/task/v1/tasks"
	"github.com/Edge-Center/edgecentercloud-go/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToImageListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct {
	Private    bool              `q:"private" validate:"omitempty"`
	Visibility types.Visibility  `q:"visibility" validate:"omitempty"`
	MetadataK  string            `q:"metadata_k" validate:"omitempty"`
	MetadataKV map[string]string `q:"metadata_kv" validate:"omitempty"`
}

// ToImageListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToImageListQuery() (string, error) {
	q, err := edgecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToImageCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create an image.
type CreateOpts struct {
	Name           string                `json:"name" required:"true" validate:"required"`
	HwMachineType  types.HwMachineType   `json:"hw_machine_type" validate:"required,enum"`
	SshKey         types.SshKeyType      `json:"ssh_key" validate:"required,enum"`
	OSType         types.OSType          `json:"os_type" validate:"required,enum"`
	IsBaremetal    *bool                 `json:"is_baremetal,omitempty"`
	HwFirmwareType types.HwFirmwareType  `json:"hw_firmware_type" validate:"required,enum"`
	Source         types.ImageSourceType `json:"source" validate:"required,enum"`
	VolumeID       string                `json:"volume_id" required:"true" validate:"required"`
}

// Validate CreateOpts.
func (opts CreateOpts) Validate() error {
	return edgecloud.Validate.Struct(opts)
}

// ToImageCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToImageCreateMap() (map[string]interface{}, error) {
	return edgecloud.BuildRequestBody(opts, "")
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToImageUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to create an image.
type UpdateOpts struct {
	HwMachineType  types.HwMachineType  `json:"hw_machine_type" validate:"required,enum"`
	SshKey         types.SshKeyType     `json:"ssh_key" validate:"required,enum"`
	Name           string               `json:"name" required:"true"`
	OSType         types.OSType         `json:"os_type" validate:"required,enum"`
	IsBaremetal    *bool                `json:"is_baremetal,omitempty"`
	HwFirmwareType types.HwFirmwareType `json:"hw_firmware_type" validate:"required,enum"`
}

// Validate UpdateOpts.
func (opts UpdateOpts) Validate() error {
	return edgecloud.Validate.Struct(opts)
}

// ToImageUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToImageUpdateMap() (map[string]interface{}, error) {
	return edgecloud.BuildRequestBody(opts, "")
}

// UploadOptsBuilder allows extensions to add additional parameters to the Upload request.
type UploadOptsBuilder interface {
	ToImageUploadMap() (map[string]interface{}, error)
}

// UploadOpts represents options used to upload an image.
type UploadOpts struct {
	OsVersion      string               `json:"os_version,omitempty"`
	HwMachineType  types.HwMachineType  `json:"hw_machine_type" validate:"required,enum"`
	SshKey         types.SshKeyType     `json:"ssh_key" validate:"required,enum"`
	Name           string               `json:"name" required:"true" validate:"required"`
	OsDistro       string               `json:"os_distro,omitempty"`
	OSType         types.OSType         `json:"os_type" validate:"required,enum"`
	URL            string               `json:"url" required:"true" validate:"required,url"`
	IsBaremetal    *bool                `json:"is_baremetal,omitempty"`
	HwFirmwareType types.HwFirmwareType `json:"hw_firmware_type" validate:"required,enum"`
	CowFormat      bool                 `json:"cow_format"`
	Metadata       map[string]string    `json:"metadata,omitempty"`
}

// Validate UploadOpts.
func (opts UploadOpts) Validate() error {
	return edgecloud.Validate.Struct(opts)
}

// ToImageUploadMap builds a request body from UploadOpts.
func (opts UploadOpts) ToImageUploadMap() (map[string]interface{}, error) {
	return edgecloud.BuildRequestBody(opts, "")
}

func List(client *edgecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToImageListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ImagePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific image based on its unique ID.
func Get(client *edgecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(client, id)
	_, r.Err = client.Get(url, &r.Body, nil)
	return
}

func ListAll(client *edgecloud.ServiceClient, opts ListOptsBuilder) ([]Image, error) {
	pages, err := List(client, opts).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractImages(pages)
	if err != nil {
		return nil, err
	}

	return all, nil

}

// Create an image.
func Create(client *edgecloud.ServiceClient, opts CreateOptsBuilder) (r tasks.Result) {
	b, err := opts.ToImageCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), b, &r.Body, nil)
	return
}

// Delete an image.
func Delete(client *edgecloud.ServiceClient, imageID string) (r tasks.Result) {
	url := deleteURL(client, imageID)
	_, r.Err = client.DeleteWithResponse(url, &r.Body, nil)
	return
}

// Update accepts a UpdateOpts struct and updates an existing image using the values provided.
func Update(client *edgecloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	url := updateURL(client, id)
	b, err := opts.ToImageUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Patch(url, b, &r.Body, nil)
	return
}

// Upload accepts a UploadOpts struct and upload an image using the values provided.
func Upload(client *edgecloud.ServiceClient, opts UploadOptsBuilder) (r tasks.Result) {
	b, err := opts.ToImageUploadMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(uploadURL(client), b, &r.Body, nil)
	return
}
