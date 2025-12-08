package edgecloud

import (
	"context"
	"fmt"
	"net/http"
)

const (
	MkaaSClustersBasePathV2 = "/internal/mkaas/v2/clusters"
)

// MkaaSService is an interface for creating and managing mkaas clusters with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/mkaas
type MkaaSService interface {
	MkaasClusters
	MkaasPools
}

type MkaasClusters interface {
	ClusterCreate(context.Context, MkaaSClusterCreateRequest) (*TaskResponse, *Response, error)
	ClustersList(context.Context, *MkaaSClusterListOptions) ([]MkaaSCluster, *Response, error)
	ClusterGet(context.Context, int) (*MkaaSCluster, *Response, error)
	ClusterDelete(context.Context, int) (*TaskResponse, *Response, error)
}

type MkaasPools interface {
	PoolCreate(ctx context.Context, clusterID int, reqBody MkaaSPoolCreateRequest) (*TaskResponse, *Response, error)
	PoolsList(ctx context.Context, clusterID int, opts *MkaaSPoolListOptions) ([]MkaaSPool, *Response, error)
	PoolGet(ctx context.Context, clusterID, poolID int) (*MkaaSPool, *Response, error)
	PoolUpdate(ctx context.Context, clusterID, poolID int, reqBody MkaaSPoolUpdateRequest) (*TaskResponse, *Response, error)
	PoolDelete(ctx context.Context, clusterID, poolID int) (*TaskResponse, *Response, error)
}

// MkaasServiceOp handles communication with mkaas methods of the EdgecenterCloud API.
type MkaasServiceOp struct {
	client *Client
}

var _ MkaaSService = &MkaasServiceOp{}

type MkaaSClustersRoot struct {
	Count    int
	Clusters []MkaaSCluster `json:"results"`
}

type MkaaSPoolsRoot struct {
	Count int
	Pools []MkaaSPool `json:"results"`
}

type MkaaSClusterListOptions struct {
	Name   string `url:"name,omitempty"`
	Status string `url:"status,omitempty"`
	Limit  int    `url:"limit,omitempty"`
	Offset int    `url:"offset,omitempty"`
}

type MkaaSClusterCreateRequest struct {
	Name           string                    `json:"name"`
	SSHKeyPairName string                    `json:"ssh_keypair_name"`
	NetworkID      string                    `json:"network_id"`
	SubnetID       string                    `json:"subnet_id"`
	ControlPlane   ControlPlaneCreateRequest `json:"control_plane"`
	Pools          []MkaaSPoolCreateRequest  `json:"pools"`
}

// MkaaSCluster represents an EdgecenterCloud MkaaS Cluster.
type MkaaSCluster struct {
	ID             int          `json:"id"`
	RegionID       int          `json:"region_id"`
	ProjectID      int          `json:"project_id"`
	SSHKeypairName string       `json:"ssh_keypair_name"`
	Name           string       `json:"name"`
	NetworkID      string       `json:"network_id"`
	SubnetID       string       `json:"subnet_id"`
	ControlPlane   ControlPlane `json:"control_plane"`
	Pools          []MkaaSPool  `json:"pools"`
	InternalIP     string       `json:"internal_ip"`
	ExternalIP     string       `json:"external_ip"`
	Existed        string       `json:"existed,omitempty"` // Duration string (e.g., "237h36m46.703341967s")
	Created        string       `json:"created"`
	Processing     bool         `json:"processing"`
	Status         string       `json:"status"`
	State          string       `json:"state"`
}

type ControlPlaneCreateRequest struct {
	Flavor     string     `json:"flavor"`
	NodeCount  int        `json:"node_count"`
	VolumeSize int        `json:"volume_size"`
	VolumeType VolumeType `json:"volume_type"`
	Version    string     `json:"version"`
}

// ControlPlane configuration.
type ControlPlane struct {
	Flavor     string     `json:"flavor"`
	NodeCount  int        `json:"node_count"`
	VolumeSize int        `json:"volume_size"`
	VolumeType VolumeType `json:"volume_type"`
	Version    string     `json:"version"`
}

type MkaaSPoolListOptions struct {
	Name   string `url:"name,omitempty"`
	Status string `url:"status,omitempty"`
	Limit  int    `url:"limit,omitempty"`
	Offset int    `url:"offset,omitempty"`
}

type MkaaSPoolCreateRequest struct {
	Name            string            `json:"name,omitempty"`
	Flavor          string            `json:"flavor,omitempty"`
	MaxNodeCount    *int              `json:"max_node_count,omitempty"`
	MinNodeCount    *int              `json:"min_node_count,omitempty"`
	NodeCount       int               `json:"node_count,omitempty"`
	VolumeSize      int               `json:"volume_size,omitempty"`
	SecurityGroupID *string           `json:"security_group_id,omitempty"`
	VolumeType      VolumeType        `json:"volume_type,omitempty"`
	Labels          map[string]string `json:"labels,omitempty"`
	Taints          []MkaaSTaint      `json:"taints,omitempty"`
}

type MkaaSPoolUpdateRequest struct {
	Name            *string           `json:"name,omitempty"`
	Flavor          *string           `json:"flavor,omitempty"`
	MaxNodeCount    *int              `json:"max_node_count,omitempty"`
	MinNodeCount    *int              `json:"min_node_count,omitempty"`
	NodeCount       *int              `json:"node_count,omitempty"`
	VolumeSize      *int              `json:"volume_size,omitempty"`
	SecurityGroupID *string           `json:"security_group_id,omitempty"`
	VolumeType      *VolumeType       `json:"volume_type,omitempty"`
	Labels          map[string]string `json:"labels,omitempty"`
	Taints          []MkaaSTaint      `json:"taints,omitempty"`
}

type MkaaSPool struct {
	ID           int               `json:"id"`
	Name         string            `json:"name"`
	Flavor       string            `json:"flavor"`
	MaxNodeCount int               `json:"max_node_count"`
	MinNodeCount int               `json:"min_node_count"`
	NodeCount    int               `json:"node_count"`
	VolumeSize   int               `json:"volume_size"`
	VolumeType   VolumeType        `json:"volume_type"`
	Labels       map[string]string `json:"labels"`
	Taints       []MkaaSTaint      `json:"taints"`
	State        string            `json:"state"`
	Status       string            `json:"status"`
}

// MkaaSTaint configuration for nodes.
type MkaaSTaint struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Effect string `json:"effect"`
}

func (m *MkaasServiceOp) ClusterCreate(ctx context.Context, reqBody MkaaSClusterCreateRequest) (*TaskResponse, *Response, error) {
	if resp, err := m.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := m.client.addProjectRegionPath(MkaaSClustersBasePathV2)

	req, err := m.client.NewRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	tasks := new(TaskResponse)
	resp, err := m.client.Do(ctx, req, tasks)
	if err != nil {
		return nil, resp, err
	}

	return tasks, resp, err
}

func (m *MkaasServiceOp) ClustersList(ctx context.Context, opts *MkaaSClusterListOptions) ([]MkaaSCluster, *Response, error) {
	if resp, err := m.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := m.client.addProjectRegionPath(MkaaSClustersBasePathV2)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := m.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(MkaaSClustersRoot)
	resp, err := m.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Clusters, resp, err
}

func (m *MkaasServiceOp) ClusterGet(ctx context.Context, clusterID int) (*MkaaSCluster, *Response, error) {
	if resp, err := m.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%d", m.client.addProjectRegionPath(MkaaSClustersBasePathV2), clusterID)

	req, err := m.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	cluster := new(MkaaSCluster)
	resp, err := m.client.Do(ctx, req, cluster)
	if err != nil {
		return nil, resp, err
	}

	return cluster, resp, err
}

func (m *MkaasServiceOp) ClusterDelete(ctx context.Context, clusterID int) (*TaskResponse, *Response, error) {
	if resp, err := m.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%d", m.client.addProjectRegionPath(MkaaSClustersBasePathV2), clusterID)

	req, err := m.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	tasks := new(TaskResponse)
	resp, err := m.client.Do(ctx, req, tasks)
	if err != nil {
		return nil, resp, err
	}

	return tasks, resp, err
}

func (m *MkaasServiceOp) PoolCreate(ctx context.Context, clusterID int, reqBody MkaaSPoolCreateRequest) (*TaskResponse, *Response, error) {
	if resp, err := m.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%d/pools", m.client.addProjectRegionPath(MkaaSClustersBasePathV2), clusterID)

	req, err := m.client.NewRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	tasks := new(TaskResponse)
	resp, err := m.client.Do(ctx, req, tasks)
	if err != nil {
		return nil, resp, err
	}

	return tasks, resp, err
}

func (m *MkaasServiceOp) PoolUpdate(ctx context.Context, clusterID, poolID int, reqBody MkaaSPoolUpdateRequest) (*TaskResponse, *Response, error) {
	if resp, err := m.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%d/pools/%d", m.client.addProjectRegionPath(MkaaSClustersBasePathV2), clusterID, poolID)

	req, err := m.client.NewRequest(ctx, http.MethodPatch, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	tasks := new(TaskResponse)
	resp, err := m.client.Do(ctx, req, tasks)
	if err != nil {
		return nil, resp, err
	}

	return tasks, resp, err
}

func (m *MkaasServiceOp) PoolsList(ctx context.Context, clusterID int, opts *MkaaSPoolListOptions) ([]MkaaSPool, *Response, error) {
	if resp, err := m.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%d/pools", m.client.addProjectRegionPath(MkaaSClustersBasePathV2), clusterID)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := m.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(MkaaSPoolsRoot)
	resp, err := m.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Pools, resp, err
}

func (m *MkaasServiceOp) PoolGet(ctx context.Context, clusterID, poolID int) (*MkaaSPool, *Response, error) {
	if resp, err := m.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%d/pools/%d", m.client.addProjectRegionPath(MkaaSClustersBasePathV2), clusterID, poolID)

	req, err := m.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	pool := new(MkaaSPool)
	resp, err := m.client.Do(ctx, req, pool)
	if err != nil {
		return nil, resp, err
	}

	return pool, resp, err
}

func (m *MkaasServiceOp) PoolDelete(ctx context.Context, clusterID, poolID int) (*TaskResponse, *Response, error) {
	if resp, err := m.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%d/pools/%d", m.client.addProjectRegionPath(MkaaSClustersBasePathV2), clusterID, poolID)

	req, err := m.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	tasks := new(TaskResponse)
	resp, err := m.client.Do(ctx, req, tasks)
	if err != nil {
		return nil, resp, err
	}

	return tasks, resp, err
}
