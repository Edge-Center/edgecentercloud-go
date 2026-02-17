package edgecloud

import (
	"context"
	"fmt"
	"net/http"
)

const (
	MKaaSClustersBasePathV2 = "/mkaas/v2/clusters"
)

// MKaaSService is an interface for creating and managing mkaas clusters with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/mkaas
type MKaaSService interface {
	MKaaSClusters
	MKaaSPools
}

type MKaaSClusters interface {
	ClusterCreate(context.Context, MKaaSClusterCreateRequest) (*TaskResponse, *Response, error)
	ClustersList(context.Context, *MKaaSClusterListOptions) ([]MKaaSCluster, *Response, error)
	ClusterGet(context.Context, int) (*MKaaSCluster, *Response, error)
	ClusterUpdate(ctx context.Context, clusterID int, reqBody MKaaSClusterUpdateRequest) (*TaskResponse, *Response, error)
	ClusterDelete(context.Context, int) (*TaskResponse, *Response, error)
}

type MKaaSPools interface {
	PoolCreate(ctx context.Context, clusterID int, reqBody MKaaSPoolCreateRequest) (*TaskResponse, *Response, error)
	PoolsList(ctx context.Context, clusterID int, opts *MKaaSPoolListOptions) ([]MKaaSPool, *Response, error)
	PoolGet(ctx context.Context, clusterID, poolID int) (*MKaaSPool, *Response, error)
	PoolUpdateName(ctx context.Context, clusterID, poolID int, reqBody MKaaSPoolUpdateNameRequest) (*TaskResponse, *Response, error)
	PoolUpdateNodeCount(ctx context.Context, clusterID, poolID int, reqBody MKaaSPoolUpdateScaleRequest) (*TaskResponse, *Response, error)
	PoolUpdateSecurityGroups(ctx context.Context, clusterID, poolID int,
		reqBody MKaaSPoolUpdateSecurityGroupsRequest) (*TaskResponse, *Response, error)
	PoolDelete(ctx context.Context, clusterID, poolID int) (*TaskResponse, *Response, error)
}

// MKaaSServiceOp handles communication with mkaas methods of the EdgecenterCloud API.
type MKaaSServiceOp struct {
	client *Client
}

var _ MKaaSService = &MKaaSServiceOp{}

type MKaaSClustersRoot struct {
	Count    int
	Clusters []MKaaSCluster `json:"results"`
}

type MKaaSPoolsRoot struct {
	Count int
	Pools []MKaaSPool `json:"results"`
}

type MKaaSClusterListOptions struct {
	Name   string `url:"name,omitempty"`
	Status string `url:"status,omitempty"`
	Limit  int    `url:"limit,omitempty"`
	Offset int    `url:"offset,omitempty"`
}

type MKaaSClusterCreateRequest struct {
	Name                     string                    `json:"name"`
	SSHKeyPairName           string                    `json:"ssh_keypair_name"`
	NetworkID                string                    `json:"network_id"`
	SubnetID                 string                    `json:"subnet_id"`
	PodSubnet                *string                   `json:"pod_subnet,omitempty"`
	ServiceSubnet            *string                   `json:"service_subnet,omitempty"`
	PublishKubeAPIToInternet bool                      `json:"publish_kube_api_to_internet"`
	ControlPlane             ControlPlaneCreateRequest `json:"control_plane"`
	Pools                    []MKaaSPoolCreateRequest  `json:"pools,omitempty"`
}

type MKaaSClusterUpdateRequest struct {
	Name            string `json:"name"`
	MasterNodeCount int    `json:"master_node_count"`
}

// MKaaSCluster represents an EdgecenterCloud MkaaS Cluster.
type MKaaSCluster struct {
	ID             int          `json:"id"`
	RegionID       int          `json:"region_id"`
	ProjectID      int          `json:"project_id"`
	SSHKeypairName string       `json:"ssh_keypair_name"`
	Name           string       `json:"name"`
	NetworkID      string       `json:"network_id"`
	SubnetID       string       `json:"subnet_id"`
	ControlPlane   ControlPlane `json:"control_plane"`
	Pools          []MKaaSPool  `json:"pools"`
	InternalIP     string       `json:"internal_ip"`
	ExternalIP     string       `json:"external_ip"`
	Existed        string       `json:"existed,omitempty"` // Duration string (e.g., "237h36m46.703341967s")
	Created        string       `json:"created"`
	Processing     bool         `json:"processing"`
	Status         string       `json:"status"`
	Stage          string       `json:"stage"`
	PodSubnet      string       `json:"pod_subnet"`
	ServiceSubnet  string       `json:"service_subnet"`
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

type MKaaSPoolListOptions struct {
	Name   string `url:"name,omitempty"`
	Status string `url:"status,omitempty"`
	Limit  int    `url:"limit,omitempty"`
	Offset int    `url:"offset,omitempty"`
}

type MKaaSPoolCreateRequest struct {
	Name             string            `json:"name,omitempty"`
	Flavor           string            `json:"flavor,omitempty"`
	MaxNodeCount     *int              `json:"max_node_count,omitempty"`
	MinNodeCount     *int              `json:"min_node_count,omitempty"`
	NodeCount        int               `json:"node_count,omitempty"`
	VolumeSize       int               `json:"volume_size,omitempty"`
	SecurityGroupID  *string           `json:"security_group_id,omitempty"`
	VolumeType       VolumeType        `json:"volume_type,omitempty"`
	Labels           map[string]string `json:"labels,omitempty"`
	Taints           []MKaaSTaint      `json:"taints,omitempty"`
	SecurityGroupIds []string          `json:"security_group_ids,omitempty"`
}

type MKaaSPoolUpdateRequest struct {
	Flavor           *string           `json:"flavor,omitempty"`
	MaxNodeCount     *int              `json:"max_node_count,omitempty"`
	MinNodeCount     *int              `json:"min_node_count,omitempty"`
	VolumeSize       *int              `json:"volume_size,omitempty"`
	SecurityGroupID  *string           `json:"security_group_id,omitempty"`
	VolumeType       *VolumeType       `json:"volume_type,omitempty"`
	Labels           map[string]string `json:"labels,omitempty"`
	Taints           []MKaaSTaint      `json:"taints,omitempty"`
	SecurityGroupIds []string          `json:"security_group_ids,omitempty"`
}

type MKaaSPoolUpdateNameRequest struct {
	Name *string `json:"name,omitempty"`
}

type MKaaSPoolUpdateScaleRequest struct {
	NodeCount *int `json:"node_count,omitempty"`
}

type MKaaSPoolUpdateSecurityGroupsRequest struct {
	SecurityGroupIds []string `json:"security_groups_ids"`
}

type MKaaSPool struct {
	ID               int               `json:"id"`
	Name             string            `json:"name"`
	Flavor           string            `json:"flavor"`
	MaxNodeCount     int               `json:"max_node_count"`
	MinNodeCount     int               `json:"min_node_count"`
	NodeCount        int               `json:"node_count"`
	VolumeSize       int               `json:"volume_size"`
	VolumeType       VolumeType        `json:"volume_type"`
	Labels           map[string]string `json:"labels"`
	Taints           []MKaaSTaint      `json:"taints"`
	State            string            `json:"state"`
	Status           string            `json:"status"`
	SecurityGroupIds []string          `json:"security_group_ids"`
}

// MKaaSTaint configuration for nodes.
type MKaaSTaint struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Effect string `json:"effect"`
}

func (m *MKaaSServiceOp) ClusterCreate(ctx context.Context, reqBody MKaaSClusterCreateRequest) (*TaskResponse, *Response, error) {
	if resp, err := m.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := m.client.addProjectRegionPath(MKaaSClustersBasePathV2)

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

func (m *MKaaSServiceOp) ClustersList(ctx context.Context, opts *MKaaSClusterListOptions) ([]MKaaSCluster, *Response, error) {
	if resp, err := m.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := m.client.addProjectRegionPath(MKaaSClustersBasePathV2)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := m.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(MKaaSClustersRoot)
	resp, err := m.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Clusters, resp, err
}

func (m *MKaaSServiceOp) ClusterGet(ctx context.Context, clusterID int) (*MKaaSCluster, *Response, error) {
	if resp, err := m.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%d", m.client.addProjectRegionPath(MKaaSClustersBasePathV2), clusterID)

	req, err := m.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	cluster := new(MKaaSCluster)
	resp, err := m.client.Do(ctx, req, cluster)
	if err != nil {
		return nil, resp, err
	}

	return cluster, resp, err
}

func (m *MKaaSServiceOp) ClusterUpdate(ctx context.Context, clusterID int, reqBody MKaaSClusterUpdateRequest) (*TaskResponse, *Response, error) {
	if resp, err := m.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%d", m.client.addProjectRegionPath(MKaaSClustersBasePathV2), clusterID)

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

func (m *MKaaSServiceOp) ClusterDelete(ctx context.Context, clusterID int) (*TaskResponse, *Response, error) {
	if resp, err := m.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%d", m.client.addProjectRegionPath(MKaaSClustersBasePathV2), clusterID)

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

func (m *MKaaSServiceOp) PoolCreate(ctx context.Context, clusterID int, reqBody MKaaSPoolCreateRequest) (*TaskResponse, *Response, error) {
	if resp, err := m.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%d/pools", m.client.addProjectRegionPath(MKaaSClustersBasePathV2), clusterID)

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

func (m *MKaaSServiceOp) PoolUpdateName(ctx context.Context, clusterID, poolID int, reqBody MKaaSPoolUpdateNameRequest) (*TaskResponse, *Response, error) {
	if resp, err := m.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%d/pools/%d/name", m.client.addProjectRegionPath(MKaaSClustersBasePathV2), clusterID, poolID)

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

func (m *MKaaSServiceOp) PoolUpdateNodeCount(ctx context.Context, clusterID, poolID int, reqBody MKaaSPoolUpdateScaleRequest) (*TaskResponse, *Response, error) {
	if resp, err := m.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%d/pools/%d/scale", m.client.addProjectRegionPath(MKaaSClustersBasePathV2), clusterID, poolID)

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

func (m *MKaaSServiceOp) PoolsList(ctx context.Context, clusterID int, opts *MKaaSPoolListOptions) ([]MKaaSPool, *Response, error) {
	if resp, err := m.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%d/pools", m.client.addProjectRegionPath(MKaaSClustersBasePathV2), clusterID)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := m.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(MKaaSPoolsRoot)
	resp, err := m.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Pools, resp, err
}

func (m *MKaaSServiceOp) PoolGet(ctx context.Context, clusterID, poolID int) (*MKaaSPool, *Response, error) {
	if resp, err := m.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%d/pools/%d", m.client.addProjectRegionPath(MKaaSClustersBasePathV2), clusterID, poolID)

	req, err := m.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	pool := new(MKaaSPool)
	resp, err := m.client.Do(ctx, req, pool)
	if err != nil {
		return nil, resp, err
	}

	return pool, resp, err
}

func (m *MKaaSServiceOp) PoolDelete(ctx context.Context, clusterID, poolID int) (*TaskResponse, *Response, error) {
	if resp, err := m.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%d/pools/%d", m.client.addProjectRegionPath(MKaaSClustersBasePathV2), clusterID, poolID)

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

func (m *MKaaSServiceOp) PoolUpdateSecurityGroups(ctx context.Context, clusterID, poolID int,
	reqBody MKaaSPoolUpdateSecurityGroupsRequest,
) (*TaskResponse, *Response, error) {
	if resp, err := m.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%d/pools/%d/secgroups", m.client.addProjectRegionPath(MKaaSClustersBasePathV2), clusterID, poolID)

	req, err := m.client.NewRequest(ctx, http.MethodPut, path, reqBody)
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
