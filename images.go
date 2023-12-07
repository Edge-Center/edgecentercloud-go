package edgecloud

import (
	"context"
	"fmt"
	"net/http"
)

const (
	imagesBasePathV1        = "/v1/images"
	bmimagesBasePathV1      = "/v1/bmimages"
	projectimagesBasePathV1 = "/v1/projectimages"
	downloadimageBasePathV1 = "/v1/downloadimage"
)

// ImagesService is an interface for creating and managing Images with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/images
type ImagesService interface {
	List(context.Context, *ImageListOptions) ([]Image, *Response, error)
	Create(context.Context, *ImageCreateRequest) (*TaskResponse, *Response, error)
	Delete(context.Context, string) (*TaskResponse, *Response, error)
	Get(context.Context, string) (*Image, *Response, error)
	Update(context.Context, string, *ImageUpdateRequest) (*Image, *Response, error)
	Upload(context.Context, *ImageUploadRequest) (*TaskResponse, *Response, error)

	ImagesBaremetal
	ImagesProject
	ImagesMetadata
}

type ImagesBaremetal interface {
	ImagesBaremetalList(context.Context, *ImageListOptions) ([]Image, *Response, error)
	ImagesBaremetalCreate(context.Context, *ImageCreateRequest) (*TaskResponse, *Response, error)
}

type ImagesProject interface {
	ImagesProjectList(context.Context) ([]Image, *Response, error)
}

type ImagesMetadata interface {
	MetadataList(context.Context, string) ([]MetadataDetailed, *Response, error)
	MetadataCreate(context.Context, string, *Metadata) (*Response, error)
	MetadataUpdate(context.Context, string, *Metadata) (*Response, error)
	MetadataDeleteItem(context.Context, string, *MetadataItemOptions) (*Response, error)
	MetadataGetItem(context.Context, string, *MetadataItemOptions) (*MetadataDetailed, *Response, error)
}

// ImagesServiceOp handles communication with Images methods of the EdgecenterCloud API.
type ImagesServiceOp struct {
	client *Client
}

var _ ImagesService = &ImagesServiceOp{}

// Image represents an EdgecenterCloud Image.
type Image struct {
	DiskFormat       string           `json:"disk_format"`
	MetadataDetailed MetadataDetailed `json:"metadata_detailed"`
	Metadata         Metadata         `json:"metadata"`
	MinRAM           int              `json:"min_ram"`
	MinDisk          int              `json:"min_disk"`
	OSVersion        string           `json:"os_version"`
	CreatedAt        string           `json:"created_at"`
	UpdatedAt        string           `json:"updated_at"`
	TaskID           string           `json:"task_id"`
	ProjectID        int              `json:"project_id"`
	RegionID         int              `json:"region_id"`
	Region           string           `json:"region"`
	CreatorTaskID    string           `json:"creator_task_id"`
	Status           string           `json:"status"`
	OSType           OSType           `json:"os_type"`
	SSHKey           SSHKey           `json:"ssh_key"`
	OSDistro         string           `json:"os_distro"`
	Visibility       string           `json:"visibility"`
	DisplayOrder     int              `json:"display_order"`
	HWFirmwareType   HWFirmwareType   `json:"hw_firmware_type"`
	Name             string           `json:"name"`
	Size             int              `json:"size"`
	IsBaremetal      bool             `json:"is_baremetal"`
	HWMachineType    HWMachineType    `json:"hw_machine_type"`
	Description      string           `json:"description"`
	ID               string           `json:"id"`
}

type OSType string

const (
	OSTypeLinux   OSType = "linux"
	OSTypeWindows OSType = "windows"
)

type SSHKey string

const (
	SSHKeyAllow    SSHKey = "allow"
	SSHKeyDeny     SSHKey = "deny"
	SSHKeyRequired SSHKey = "required"
)

type HWFirmwareType string

const (
	HWFirmwareTypeBios HWFirmwareType = "bios"
	HWFirmwareTypeUEFI HWFirmwareType = "uefi"
)

type HWMachineType string

const (
	HWMachineTypeI440 HWMachineType = "i440"
	HWMachineTypeQ35  HWMachineType = "q35"
)

type ImageListOptions struct {
	Visibility string `url:"visibility,omitempty"  validate:"omitempty"`
	Private    string `url:"private,omitempty"  validate:"omitempty"`
	MetadataKV string `url:"metadata_kv,omitempty"  validate:"omitempty"`
	MetadataK  string `url:"metadata_k,omitempty"  validate:"omitempty"`
}

type ImageCreateRequest struct {
	SSHKey         SSHKey         `json:"ssh_key"`
	OSType         OSType         `json:"os_type"`
	Name           string         `json:"name" required:"true" validate:"required"`
	IsBaremetal    bool           `json:"is_baremetal"`
	VolumeID       string         `json:"volume_id" required:"true" validate:"required"`
	HWMachineType  HWMachineType  `json:"hw_machine_type"`
	HWFirmwareType HWFirmwareType `json:"hw_firmware_type"`
	Source         string         `json:"source"`
	Metadata       Metadata       `json:"metadata"`
}

type ImageUpdateRequest struct {
	SSHKey         SSHKey         `json:"ssh_key"`
	Name           string         `json:"name"`
	IsBaremetal    bool           `json:"is_baremetal"`
	HWMachineType  HWMachineType  `json:"hw_machine_type"`
	HWFirmwareType HWFirmwareType `json:"hw_firmware_type"`
	OSType         OSType         `json:"os_type"`
	Metadata       Metadata       `json:"metadata"`
}

type ImageUploadRequest struct {
	SSHKey         SSHKey         `json:"ssh_key"`
	OSDistro       string         `json:"os_distro"`
	Name           string         `json:"name" required:"true" validate:"required"`
	URL            string         `json:"url" required:"true" validate:"required"`
	COWFormat      bool           `json:"cow_format"`
	OSVersion      string         `json:"os_version"`
	IsBaremetal    bool           `json:"is_baremetal"`
	HWMachineType  HWMachineType  `json:"hw_machine_type"`
	HWFirmwareType HWFirmwareType `json:"hw_firmware_type"`
	OSType         OSType         `json:"os_type"`
	Metadata       Metadata       `json:"metadata"`
}

// imagesRoot represents an Images root.
type imagesRoot struct {
	Count  int
	Images []Image `json:"results"`
}

// List get images.
func (s *ImagesServiceOp) List(ctx context.Context, opts *ImageListOptions) ([]Image, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(imagesBasePathV1)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(imagesRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Images, resp, err
}

// Create an Image.
func (s *ImagesServiceOp) Create(ctx context.Context, reqBody *ImageCreateRequest) (*TaskResponse, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(imagesBasePathV1)

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

// Get an image.
func (s *ImagesServiceOp) Get(ctx context.Context, imageID string) (*Image, *Response, error) {
	if resp, err := isValidUUID(imageID, "imageID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(imagesBasePathV1), imageID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	image := new(Image)
	resp, err := s.client.Do(ctx, req, image)
	if err != nil {
		return nil, resp, err
	}

	return image, resp, err
}

// Delete an image.
func (s *ImagesServiceOp) Delete(ctx context.Context, imageID string) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(imageID, "imageID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(imagesBasePathV1), imageID)

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

// Update image fields.
func (s *ImagesServiceOp) Update(ctx context.Context, imageID string, reqBody *ImageUpdateRequest) (*Image, *Response, error) {
	if resp, err := isValidUUID(imageID, "imageID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(imagesBasePathV1), imageID)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	image := new(Image)
	resp, err := s.client.Do(ctx, req, image)
	if err != nil {
		return nil, resp, err
	}

	return image, resp, err
}

// Upload an Image.
func (s *ImagesServiceOp) Upload(ctx context.Context, reqBody *ImageUploadRequest) (*TaskResponse, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(downloadimageBasePathV1)

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

// ImagesBaremetalList get images of baremetal instances.
func (s *ImagesServiceOp) ImagesBaremetalList(ctx context.Context, opts *ImageListOptions) ([]Image, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(bmimagesBasePathV1)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(imagesRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Images, resp, err
}

// ImagesBaremetalCreate an Image.
func (s *ImagesServiceOp) ImagesBaremetalCreate(ctx context.Context, reqBody *ImageCreateRequest) (*TaskResponse, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(bmimagesBasePathV1)

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

// ImagesProjectList get images owned by a project.
func (s *ImagesServiceOp) ImagesProjectList(ctx context.Context) ([]Image, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(projectimagesBasePathV1)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(imagesRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Images, resp, err
}

// MetadataList security group detailed metadata items.
func (s *ImagesServiceOp) MetadataList(ctx context.Context, imageID string) ([]MetadataDetailed, *Response, error) {
	if resp, err := isValidUUID(imageID, "imageID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	return metadataList(ctx, s.client, imageID, imagesBasePathV1)
}

// MetadataCreate or update security group metadata.
func (s *ImagesServiceOp) MetadataCreate(ctx context.Context, imageID string, reqBody *Metadata) (*Response, error) {
	if resp, err := isValidUUID(imageID, "imageID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	return metadataCreate(ctx, s.client, imageID, imagesBasePathV1, reqBody)
}

// MetadataUpdate security group metadata.
func (s *ImagesServiceOp) MetadataUpdate(ctx context.Context, imageID string, reqBody *Metadata) (*Response, error) {
	if resp, err := isValidUUID(imageID, "imageID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	return metadataUpdate(ctx, s.client, imageID, imagesBasePathV1, reqBody)
}

// MetadataDeleteItem a security group metadata item by key.
func (s *ImagesServiceOp) MetadataDeleteItem(ctx context.Context, imageID string, opts *MetadataItemOptions) (*Response, error) {
	if resp, err := isValidUUID(imageID, "imageID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	return metadataDeleteItem(ctx, s.client, imageID, imagesBasePathV1, opts)
}

// MetadataGetItem security group detailed metadata.
func (s *ImagesServiceOp) MetadataGetItem(ctx context.Context, imageID string, opts *MetadataItemOptions) (*MetadataDetailed, *Response, error) {
	if resp, err := isValidUUID(imageID, "imageID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	return metadataGetItem(ctx, s.client, imageID, imagesBasePathV1, opts)
}
