package edgecloud

import (
	"context"
	"fmt"
	"net/http"
)

const (
	volumesBasePathV1 = "/v1/volumes"
)

const (
	volumesRetype = "retype"
	volumesExtend = "extend"
	volumesAttach = "attach"
	volumesDetach = "detach"
	volumesRevert = "revert"
)

// VolumesService is an interface for creating and managing Volumes with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/volumes
type VolumesService interface {
	List(context.Context, *VolumeListOptions) ([]Volume, *Response, error)
	Create(context.Context, *VolumeCreateRequest) (*TaskResponse, *Response, error)
	Get(context.Context, string) (*Volume, *Response, error)
	Delete(context.Context, string) (*TaskResponse, *Response, error)
	ChangeType(context.Context, string, *VolumeChangeTypeRequest) (*Volume, *Response, error)
	Extend(context.Context, string, *VolumeExtendSizeRequest) (*TaskResponse, *Response, error)
	Rename(context.Context, string, *Name) (*Volume, *Response, error)
	Attach(context.Context, string, *VolumeAttachRequest) (*Volume, *Response, error)
	Detach(context.Context, string, *VolumeDetachRequest) (*Volume, *Response, error)
	Revert(context.Context, string) (*TaskResponse, *Response, error)

	VolumeMetadata
}

type VolumeMetadata interface {
	MetadataList(context.Context, string) ([]MetadataDetailed, *Response, error)
	MetadataCreate(context.Context, string, *MetadataCreateRequest) (*Response, error)
	MetadataUpdate(context.Context, string, *MetadataCreateRequest) (*Response, error)
	MetadataDeleteItem(context.Context, string, *MetadataItemOptions) (*Response, error)
	MetadataGetItem(context.Context, string, *MetadataItemOptions) (*MetadataDetailed, *Response, error)
}

// VolumesServiceOp handles communication with Volumes methods of the EdgecenterCloud API.
type VolumesServiceOp struct {
	client *Client
}

var _ VolumesService = &VolumesServiceOp{}

// Volume represents an EdgecenterCloud Volume.
type Volume struct {
	ID                  string             `json:"id"`
	Name                string             `json:"name"`
	Status              string             `json:"status"` // todo: need to implement volume status type
	Size                int                `json:"size"`
	CreatedAt           string             `json:"created_at"`
	UpdatedAt           string             `json:"updated_at"`
	VolumeType          VolumeType         `json:"volume_type"`
	Device              string             `json:"device"`
	InstanceID          string             `json:"instance_id"`
	Bootable            bool               `json:"bootable"`
	CreatorTaskID       string             `json:"creator_task_id"`
	TaskID              string             `json:"task_id"`
	Metadata            Metadata           `json:"metadata"`
	MetadataDetailed    []MetadataDetailed `json:"metadata_detailed,omitempty"`
	SnapshotIDs         []string           `json:"snapshot_ids"`
	Region              string             `json:"region"`
	RegionID            int                `json:"region_id"`
	ProjectID           int                `json:"project_id"`
	Attachments         []Attachment       `json:"attachments"`
	VolumeImageMetadata Metadata           `json:"volume_image_metadata"`
	LimiterStats        LimiterStats       `json:"limiter_stats"`
	AvailabilityZone    string             `json:"availability_zone"`
}

// LimiterStats represents a limiter_stats structure.
type LimiterStats struct {
	IopsBaseLimit  int `json:"iops_base_limit"`
	IopsBurstLimit int `json:"iops_burst_limit"`
	MBpsBaseLimit  int `json:"MBps_base_limit"`
	MBpsBurstLimit int `json:"MBps_burst_limit"`
}

// Attachment represents an attachment structure.
type Attachment struct {
	ServerID     string `json:"server_id"`
	InstanceName string `json:"instance_name"`
	AttachmentID string `json:"attachment_id"`
	VolumeID     string `json:"volume_id"`
	Device       string `json:"device"`
	AttachedAt   string `json:"attached_at"`
}

type VolumeType string

const (
	VolumeTypeStandard  VolumeType = "standard"
	VolumeTypeSsdHiIops VolumeType = "ssd_hiiops"
	VolumeTypeSsdLocal  VolumeType = "ssd_local"
	VolumeTypeCold      VolumeType = "cold"
	VolumeTypeUltra     VolumeType = "ultra"
)

type VolumeSource string

const (
	VolumeSourceNewVolume      VolumeSource = "new-volume"
	VolumeSourceImage          VolumeSource = "image"
	VolumeSourceSnapshot       VolumeSource = "snapshot"
	VolumeSourceExistingVolume VolumeSource = "existing-volume"
	VolumeSourceAppTemplate    VolumeSource = "apptemplate"
)

// VolumeCreateRequest represents a request to create a Volume.
type VolumeCreateRequest struct {
	AttachmentTag        string       `json:"attachment_tag,omitempty" validate:"omitempty,required_with=InstanceIDToAttachTo"`
	ImageID              string       `json:"image_id,omitempty" validate:"rfe=Source:image,allowed_without=SnapshotID,omitempty,uuid4"`
	InstanceIDToAttachTo string       `json:"instance_id_to_attach_to,omitempty" validate:"omitempty,uuid4"`
	LifeCyclePolicyIDs   []int        `json:"lifecycle_policy_ids,omitempty"`
	Metadata             Metadata     `json:"metadata,omitempty" validate:"omitempty,dive"`
	Name                 string       `json:"name" required:"true" validate:"required"`
	Size                 int          `json:"size,omitempty"`
	SnapshotID           string       `json:"snapshot_id,omitempty" validate:"rfe=Source:snapshot,allowed_without=ImageID,omitempty,uuid4"`
	Source               VolumeSource `json:"source" required:"true" validate:"required,enum"`
	TypeName             VolumeType   `json:"type_name" required:"true" validate:"required,enum"`
}

// VolumeListOptions specifies the optional query parameters to List method.
type VolumeListOptions struct {
	InstanceID     string `url:"instance_id,omitempty"  validate:"omitempty"`
	ClusterID      string `url:"cluster_id,omitempty"  validate:"omitempty"`
	Limit          int    `url:"limit,omitempty"  validate:"omitempty"`
	Offset         int    `url:"offset,omitempty"  validate:"omitempty"`
	Bootable       bool   `url:"bootable,omitempty"  validate:"omitempty"`
	HasAttachments bool   `url:"has_attachments,omitempty"  validate:"omitempty"`
	IDPart         string `url:"id_part,omitempty"  validate:"omitempty"`
	NamePart       string `url:"name_part,omitempty"  validate:"omitempty"`
	MetadataKV     string `url:"metadata_kv,omitempty"  validate:"omitempty"`
	MetadataK      string `url:"metadata_k,omitempty"  validate:"omitempty"`
}

// VolumeChangeTypeRequest represents a request to change a Volume type.
type VolumeChangeTypeRequest struct {
	VolumeType VolumeType `json:"volume_type" required:"true" validate:"required,enum"`
}

// VolumeExtendSizeRequest represents a request to extend a Volume size.
type VolumeExtendSizeRequest struct {
	Size int `json:"size" required:"true" validate:"required"`
}

// VolumeAttachRequest represents a request to attach a Volume to Instance.
type VolumeAttachRequest struct {
	InstanceID    string `json:"instance_id" required:"true" validate:"required"`
	AttachmentTag string `json:"attachment_tag,omitempty" validate:"omitempty"`
}

// VolumeDetachRequest represents a request to detach a Volume from an Instance.
type VolumeDetachRequest struct {
	InstanceID string `json:"instance_id" required:"true" validate:"required"`
}

// volumesRoot represents a Volume root.
type volumesRoot struct {
	Count  int
	Volume []Volume `json:"results"`
}

// List get volumes.
func (s *VolumesServiceOp) List(ctx context.Context, opts *VolumeListOptions) ([]Volume, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(volumesBasePathV1)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(volumesRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Volume, resp, err
}

// Get individual Volume.
func (s *VolumesServiceOp) Get(ctx context.Context, volumeID string) (*Volume, *Response, error) {
	if resp, err := isValidUUID(volumeID, "volumeID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(volumesBasePathV1), volumeID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	volume := new(Volume)
	resp, err := s.client.Do(ctx, req, volume)
	if err != nil {
		return nil, resp, err
	}

	return volume, resp, err
}

// Create a Volume.
func (s *VolumesServiceOp) Create(ctx context.Context, reqBody *VolumeCreateRequest) (*TaskResponse, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(volumesBasePathV1)

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

// Delete the Volume.
func (s *VolumesServiceOp) Delete(ctx context.Context, volumeID string) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(volumeID, "volumeID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(volumesBasePathV1), volumeID)

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

// ChangeType of the volume.
func (s *VolumesServiceOp) ChangeType(ctx context.Context, volumeID string, reqBody *VolumeChangeTypeRequest) (*Volume, *Response, error) {
	if resp, err := isValidUUID(volumeID, "volumeID"); err != nil {
		return nil, resp, err
	}

	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(volumesBasePathV1), volumeID, volumesRetype)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	volume := new(Volume)
	resp, err := s.client.Do(ctx, req, volume)
	if err != nil {
		return nil, resp, err
	}

	return volume, resp, err
}

// Extend the volume size.
func (s *VolumesServiceOp) Extend(ctx context.Context, volumeID string, reqBody *VolumeExtendSizeRequest) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(volumeID, "volumeID"); err != nil {
		return nil, resp, err
	}

	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(volumesBasePathV1), volumeID, volumesExtend)

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

// Rename the volume.
func (s *VolumesServiceOp) Rename(ctx context.Context, volumeID string, reqBody *Name) (*Volume, *Response, error) {
	if resp, err := isValidUUID(volumeID, "volumeID"); err != nil {
		return nil, resp, err
	}

	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(volumesBasePathV1), volumeID)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	volume := new(Volume)
	resp, err := s.client.Do(ctx, req, volume)
	if err != nil {
		return nil, resp, err
	}

	return volume, resp, err
}

// Attach the volume.
func (s *VolumesServiceOp) Attach(ctx context.Context, volumeID string, reqBody *VolumeAttachRequest) (*Volume, *Response, error) {
	if resp, err := isValidUUID(volumeID, "volumeID"); err != nil {
		return nil, resp, err
	}

	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(volumesBasePathV1), volumeID, volumesAttach)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	volume := new(Volume)
	resp, err := s.client.Do(ctx, req, volume)
	if err != nil {
		return nil, resp, err
	}

	return volume, resp, err
}

// Detach the volume.
func (s *VolumesServiceOp) Detach(ctx context.Context, volumeID string, reqBody *VolumeDetachRequest) (*Volume, *Response, error) {
	if resp, err := isValidUUID(volumeID, "volumeID"); err != nil {
		return nil, resp, err
	}

	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(volumesBasePathV1), volumeID, volumesDetach)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	volume := new(Volume)
	resp, err := s.client.Do(ctx, req, volume)
	if err != nil {
		return nil, resp, err
	}

	return volume, resp, err
}

// Revert a volume to its last snapshot.
func (s *VolumesServiceOp) Revert(ctx context.Context, volumeID string) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(volumeID, "volumeID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(volumesBasePathV1), volumeID, volumesRevert)

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

// MetadataList volume detailed metadata items.
func (s *VolumesServiceOp) MetadataList(ctx context.Context, volumeID string) ([]MetadataDetailed, *Response, error) {
	if resp, err := isValidUUID(volumeID, "volumeID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	return metadataList(ctx, s.client, volumeID, volumesBasePathV1)
}

// MetadataCreate or update volume metadata.
func (s *VolumesServiceOp) MetadataCreate(ctx context.Context, volumeID string, reqBody *MetadataCreateRequest) (*Response, error) {
	if resp, err := isValidUUID(volumeID, "volumeID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	return metadataCreate(ctx, s.client, volumeID, volumesBasePathV1, reqBody)
}

// MetadataUpdate volume metadata.
func (s *VolumesServiceOp) MetadataUpdate(ctx context.Context, volumeID string, reqBody *MetadataCreateRequest) (*Response, error) {
	if resp, err := isValidUUID(volumeID, "volumeID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	return metadataUpdate(ctx, s.client, volumeID, volumesBasePathV1, reqBody)
}

// MetadataDeleteItem a volume metadata item by key.
func (s *VolumesServiceOp) MetadataDeleteItem(ctx context.Context, volumeID string, opts *MetadataItemOptions) (*Response, error) {
	if resp, err := isValidUUID(volumeID, "volumeID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	return metadataDeleteItem(ctx, s.client, volumeID, volumesBasePathV1, opts)
}

// MetadataGetItem volume detailed metadata.
func (s *VolumesServiceOp) MetadataGetItem(ctx context.Context, volumeID string, opts *MetadataItemOptions) (*MetadataDetailed, *Response, error) {
	if resp, err := isValidUUID(volumeID, "volumeID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	return metadataGetItem(ctx, s.client, volumeID, volumesBasePathV1, opts)
}
