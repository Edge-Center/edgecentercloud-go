package edgecloud

import (
	"context"
	"fmt"
	"net/http"
)

const (
	snapshotsBasePathV1 = "/v1/snapshots"
)

// SnapshotsService is an interface for creating and managing Snapshots with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/snapshots
type SnapshotsService interface {
	List(context.Context, *SnapshotListOptions) ([]Snapshot, *Response, error)
	Create(context.Context, *SnapshotCreateRequest) (*TaskResponse, *Response, error)
	Delete(context.Context, string) (*TaskResponse, *Response, error)
	Get(context.Context, string) (*Snapshot, *Response, error)
	MetadataUpdate(context.Context, string, *MetadataCreateRequest) (*Snapshot, *Response, error)
}

// SnapshotsServiceOp handles communication with Snapshots methods of the EdgecenterCloud API.
type SnapshotsServiceOp struct {
	client *Client
}

var _ SnapshotsService = &SnapshotsServiceOp{}

// Snapshot represents an EdgecenterCloud Snapshot.
type Snapshot struct {
	Region        string   `json:"region"`
	UpdatedAt     *string  `json:"updated_at"`
	CreatedAt     string   `json:"created_at"`
	Name          string   `json:"name"`
	ID            string   `json:"id"`
	RegionID      int      `json:"region_id"`
	ProjectID     int      `json:"project_id"`
	TaskID        *string  `json:"task_id"`
	Status        string   `json:"status"`
	CreatorTaskID *string  `json:"creator_task_id"`
	Size          int      `json:"size"`
	VolumeID      string   `json:"volume_id"`
	Description   string   `json:"description"`
	Metadata      Metadata `json:"metadata"`
}

// SnapshotListOptions specifies the optional query parameters to List method.
type SnapshotListOptions struct {
	VolumeID          string `url:"volume_id,omitempty"  validate:"omitempty"`
	InstanceID        string `url:"instance_id,omitempty"  validate:"omitempty"`
	ScheduleID        string `url:"schedule_id,omitempty"  validate:"omitempty"`
	LifecyclePolicyID int    `url:"lifecycle_policy_id,omitempty"  validate:"omitempty"`
	Limit             int    `url:"limit,omitempty"  validate:"omitempty"`
	Offset            int    `url:"offset,omitempty"  validate:"omitempty"`
}

type SnapshotCreateRequest struct {
	Description string   `json:"description,omitempty"`
	VolumeID    string   `json:"volume_id" required:"true" validate:"required"`
	Name        string   `json:"name" required:"true" validate:"required"`
	Metadata    Metadata `json:"metadata,omitempty"`
}

// snapshotsRoot represents a Snapshots root.
type snapshotsRoot struct {
	Count     int
	Snapshots []Snapshot `json:"results"`
}

// List get Snapshots.
func (s *SnapshotsServiceOp) List(ctx context.Context, opts *SnapshotListOptions) ([]Snapshot, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(snapshotsBasePathV1)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(snapshotsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Snapshots, resp, err
}

// Create a Snapshot.
func (s *SnapshotsServiceOp) Create(ctx context.Context, reqBody *SnapshotCreateRequest) (*TaskResponse, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(snapshotsBasePathV1)

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

// Delete a Snapshot.
func (s *SnapshotsServiceOp) Delete(ctx context.Context, snapshotID string) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(snapshotID, "snapshotID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(snapshotsBasePathV1), snapshotID)

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

// Get a Snapshot.
func (s *SnapshotsServiceOp) Get(ctx context.Context, snapshotID string) (*Snapshot, *Response, error) {
	if resp, err := isValidUUID(snapshotID, "snapshotID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(snapshotsBasePathV1), snapshotID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	router := new(Snapshot)
	resp, err := s.client.Do(ctx, req, router)
	if err != nil {
		return nil, resp, err
	}

	return router, resp, err
}

// MetadataUpdate updates snapshot metadata.
func (s *SnapshotsServiceOp) MetadataUpdate(ctx context.Context, snapshotID string, reqBody *MetadataCreateRequest) (*Snapshot, *Response, error) {
	if resp, err := isValidUUID(snapshotID, "snapshotID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/%s", s.client.addProjectRegionPath(snapshotsBasePathV1), snapshotID, metadataPath)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, reqBody.Metadata)
	if err != nil {
		return nil, nil, err
	}

	router := new(Snapshot)
	resp, err := s.client.Do(ctx, req, router)
	if err != nil {
		return nil, resp, err
	}

	return router, resp, err
}
