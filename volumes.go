package edgecloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

const (
	volumesBasePathV1 = "/v1/volumes"
)

// VolumesService is an interface for creating and managing Volumes with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/volumes
type VolumesService interface {
	Create(context.Context, *VolumeCreateRequest, *ServicePath) (*TaskResponse, *Response, error)
	Get(context.Context, string, *ServicePath) (*Volume, *Response, error)
	Delete(context.Context, string, *ServicePath) (*TaskResponse, *Response, error)
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
	LimiterStats        []LimiterStats     `json:"limiter_stats"`
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
	Standard  VolumeType = "standard"
	SsdHiIops VolumeType = "ssd_hiiops"
	SsdLocal  VolumeType = "ssd_local"
	Cold      VolumeType = "cold"
	Ultra     VolumeType = "ultra"
)

type VolumeSource string

const (
	NewVolume      VolumeSource = "new-volume"
	Image          VolumeSource = "image"
	Snapshot       VolumeSource = "snapshot"
	ExistingVolume VolumeSource = "existing-volume"
	AppTemplate    VolumeSource = "apptemplate"
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

// volumeRoot represents a Volume root.
type volumeRoot struct {
	Volume *Volume       `json:"volume"`
	Tasks  *TaskResponse `json:"tasks"`
}

// Get individual Volume.
func (s *VolumesServiceOp) Get(ctx context.Context, volumeID string, p *ServicePath) (*Volume, *Response, error) {
	if _, err := uuid.Parse(volumeID); err != nil {
		return nil, nil, NewArgError("volumeID", "should be the correct UUID")
	}

	if p == nil {
		return nil, nil, NewArgError("ServicePath", "cannot be nil")
	}

	path := addServicePath(volumesBasePathV1, p)
	path = fmt.Sprintf("%s/%s", path, volumeID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(volumeRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Volume, resp, err
}

// Create a Volume.
func (s *VolumesServiceOp) Create(ctx context.Context, createRequest *VolumeCreateRequest, p *ServicePath) (*TaskResponse, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	if p == nil {
		return nil, nil, NewArgError("ServicePath", "cannot be nil")
	}

	path := addServicePath(volumesBasePathV1, p)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(volumeRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Tasks, resp, err
}

// Delete the Volume.
func (s *VolumesServiceOp) Delete(ctx context.Context, volumeID string, p *ServicePath) (*TaskResponse, *Response, error) {
	if _, err := uuid.Parse(volumeID); err != nil {
		return nil, nil, NewArgError("volumeID", "should be the correct UUID")
	}

	if p == nil {
		return nil, nil, NewArgError("ServicePath", "cannot be nil")
	}

	path := addServicePath(volumesBasePathV1, p)
	path = fmt.Sprintf("%s/%s", path, volumeID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(volumeRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Tasks, resp, err
}