package edgecloud

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

const (
	loadbalancersBasePathV1      = "/v1/loadbalancers"
	lblistenersBasePathV1        = "/v1/lblisteners"
	lbpoolsBasePathV1            = "/v1/lbpools"
	lbflavorsBasePathV1          = "/v1/lbflavors"
	loadbalancersCheckLimitsPath = "check_limits"
)

// LoadbalancersService is an interface for creating and managing Loadbalancer with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/loadbalancers
type LoadbalancersService interface {
	List(context.Context, *LoadbalancerListOptions) ([]Loadbalancer, *Response, error)
	Get(context.Context, string) (*Loadbalancer, *Response, error)
	Create(context.Context, *LoadbalancerCreateRequest) (*TaskResponse, *Response, error)
	Delete(context.Context, string) (*TaskResponse, *Response, error)
	CheckLimits(context.Context, *LoadbalancerCheckLimitsRequest) (*map[string]int, *Response, error)

	LoadbalancerListeners
	LoadbalancerPools
	LoadbalancerFlavors
	LoadbalancerMetadata
}

type LoadbalancerListeners interface {
	ListenerGet(context.Context, string) (*Listener, *Response, error)
	ListenerCreate(context.Context, *ListenerCreateRequest) (*TaskResponse, *Response, error)
	ListenerDelete(context.Context, string) (*TaskResponse, *Response, error)
}

type LoadbalancerPools interface {
	PoolGet(context.Context, string) (*Pool, *Response, error)
	PoolCreate(context.Context, *PoolCreateRequest) (*TaskResponse, *Response, error)
	PoolDelete(context.Context, string) (*TaskResponse, *Response, error)
	PoolUpdate(context.Context, string, *PoolUpdateRequest) (*TaskResponse, *Response, error)
	PoolList(context.Context, *PoolListOptions) ([]Pool, *Response, error)
}

type LoadbalancerFlavors interface {
	FlavorList(context.Context, *FlavorsOptions) ([]Flavor, *Response, error)
}

type LoadbalancerMetadata interface {
	MetadataList(context.Context, string) ([]MetadataDetailed, *Response, error)
	MetadataCreate(context.Context, string, *MetadataCreateRequest) (*Response, error)
	MetadataUpdate(context.Context, string, *MetadataCreateRequest) (*Response, error)
	MetadataDeleteItem(context.Context, string, *MetadataItemOptions) (*Response, error)
	MetadataGetItem(context.Context, string, *MetadataItemOptions) (*MetadataDetailed, *Response, error)
}

// LoadbalancersServiceOp handles communication with Loadbalancers methods of the EdgecenterCloud API.
type LoadbalancersServiceOp struct {
	client *Client
}

var _ LoadbalancersService = &LoadbalancersServiceOp{}

// Loadbalancer represents an EdgecenterCloud Loadbalancer.
type Loadbalancer struct {
	ID                 string             `json:"id"`
	Name               string             `json:"name"`
	Flavor             Flavor             `json:"flavor"`
	VipAddress         net.IP             `json:"vip_address"`
	VipPortID          string             `json:"vip_port_id"`
	VipNetworkID       string             `json:"vip_network_id"`
	ProvisioningStatus ProvisioningStatus `json:"provisioning_status"`
	OperationStatus    OperatingStatus    `json:"operating_status"`
	CreatedAt          string             `json:"created_at"`
	UpdatedAt          string             `json:"updated_at"`
	CreatorTaskID      string             `json:"creator_task_id"`
	TaskID             string             `json:"task_id"`
	Metadata           Metadata           `json:"metadata,omitempty"`
	Stats              LoadbalancerStats  `json:"stats"`
	Listeners          []Listener         `json:"listeners"`
	FloatingIPs        []FloatingIP       `json:"floating_ips"`
	VrrpIPs            []VrrpIP           `json:"vrrp_ips"`
	ProjectID          int                `json:"project_id"`
	RegionID           int                `json:"region_id"`
	Region             string             `json:"region"`
}

// Listener represents an EdgecenterCloud Loadbalancer Listener.
type Listener struct {
	ID                 string             `json:"id"`
	Name               string             `json:"name"`
	Description        string             `json:"description"`
	Protocol           string             `json:"protocol"`
	ProtocolPort       int                `json:"protocol_port"`
	OperatingStatus    string             `json:"operating_status"`
	ProvisioningStatus ProvisioningStatus `json:"provisioning_status"`
}

// Pool represents an EdgecenterCloud Loadbalancer Pool.
type Pool struct {
	ID                    string                   `json:"id"`
	Name                  string                   `json:"name"`
	LoadBalancerAlgorithm LoadBalancerAlgorithm    `json:"lb_algorithm"`
	Protocol              LoadbalancerPoolProtocol `json:"protocol"`
	LoadBalancers         []LoadbalancerID         `json:"loadbalancers"`
	Listeners             []ListenersID            `json:"listeners"`
	Members               []PoolMember             `json:"members"`
	HealthMonitor         HealthMonitor            `json:"healthmonitor"`
	SessionPersistence    SessionPersistence       `json:"session_persistence"`
	ProvisioningStatus    ProvisioningStatus       `json:"provisioning_status"`
	OperatingStatus       OperatingStatus          `json:"operating_status"`
	CreatorTaskID         string                   `json:"creator_task_id"`
	TaskID                string                   `json:"task_id"`
	TimeoutClientData     int                      `json:"timeout_client_data"`
	TimeoutMemberData     int                      `json:"timeout_member_data"`
	TimeoutMemberConnect  int                      `json:"timeout_member_connect"`
}

// PoolMember represents an EdgecenterCloud Loadbalancer Pool PoolMember.
type PoolMember struct {
	ID              string          `json:"id"`
	OperatingStatus OperatingStatus `json:"operating_status,omitempty"`
	PoolMemberCreateRequest
}

// HealthMonitor represents an EdgecenterCloud Loadbalancer Pool HealthMonitor.
type HealthMonitor struct {
	HealthMonitorCreateRequest
}

// LoadbalancerID represents a loadbalancer ID struct.
type LoadbalancerID struct {
	ID string `json:"id"`
}

// ListenersID represents a listener ID struct.
type ListenersID struct {
	ID string `json:"id"`
}

// LoadbalancerStats represents an EdgecenterCloud Loadbalancer statistic.
type LoadbalancerStats struct {
	ActiveConnections int `json:"active_connections"`
	BytesIn           int `json:"bytes_in"`
	BytesOut          int `json:"bytes_out"`
	RequestErrors     int `json:"request_errors"`
	TotalConnections  int `json:"total_connections"`
}

// VrrpIP represents an EdgecenterCloud Loadbalancer VrrpIP (Virtual Router Redundancy Protocol IP).
type VrrpIP struct {
	VrrpIPAddress string `json:"vrrp_ip"`
}

type OperatingStatus string

const (
	OperatingStatusOnline    OperatingStatus = "ONLINE"
	OperatingStatusDraining  OperatingStatus = "DRAINING"
	OperatingStatusOffline   OperatingStatus = "OFFLINE"
	OperatingStatusDegraded  OperatingStatus = "DEGRADED"
	OperatingStatusError     OperatingStatus = "ERROR"
	OperatingStatusNoMonitor OperatingStatus = "NO_MONITOR"
)

type ProvisioningStatus string

const (
	ProvisioningStatusActive        ProvisioningStatus = "ACTIVE"
	ProvisioningStatusDeleted       ProvisioningStatus = "DELETED"
	ProvisioningStatusError         ProvisioningStatus = "ERROR"
	ProvisioningStatusPendingCreate ProvisioningStatus = "PENDING_CREATE"
	ProvisioningStatusPendingUpdate ProvisioningStatus = "PENDING_UPDATE"
	ProvisioningStatusPendingDelete ProvisioningStatus = "PENDING_DELETE"
)

// LoadbalancerCreateRequest represents a request to create a Loadbalancer.
type LoadbalancerCreateRequest struct {
	Name         string                              `json:"name" required:"true" validate:"required,name"`
	Flavor       string                              `json:"flavor,omitempty"`
	Listeners    []LoadbalancerListenerCreateRequest `json:"listeners,omitempty" validate:"omitempty,dive"`
	VipPortID    string                              `json:"vip_port_id,omitempty"`
	VipNetworkID string                              `json:"vip_network_id,omitempty"`
	VipSubnetID  string                              `json:"vip_subnet_id,omitempty"`
	Metadata     Metadata                            `json:"metadata,omitempty" validate:"omitempty,dive"`
	Tags         []string                            `json:"tag,omitempty"`
	FloatingIP   *InterfaceFloatingIP                `json:"floating_ip,omitempty" validate:"omitempty,dive"`
}

type LoadBalancerAlgorithm string

const (
	LoadBalancerAlgorithmRoundRobin       LoadBalancerAlgorithm = "ROUND_ROBIN"
	LoadBalancerAlgorithmLeastConnections LoadBalancerAlgorithm = "LEAST_CONNECTIONS"
	LoadBalancerAlgorithmSourceIP         LoadBalancerAlgorithm = "SOURCE_IP"
	LoadBalancerAlgorithmSourceIPPort     LoadBalancerAlgorithm = "SOURCE_IP_PORT"
)

type LoadbalancerPoolProtocol string

const (
	LBPoolProtocolHTTP  LoadbalancerPoolProtocol = "HTTP"
	LBPoolProtocolHTTPS LoadbalancerPoolProtocol = "HTTPS"
	LBPoolProtocolTCP   LoadbalancerPoolProtocol = "TCP"
	LBPoolProtocolUDP   LoadbalancerPoolProtocol = "UDP"
	LBPoolProtocolProxy LoadbalancerPoolProtocol = "PROXY"
)

// LoadbalancerPoolCreateRequest represents a request to create a Loadbalancer Pool.
// Used as part of a request to create a Loadbalancer.
type LoadbalancerPoolCreateRequest struct {
	Name                  string                          `json:"name" required:"true" validate:"required,name"`
	LoadBalancerAlgorithm LoadBalancerAlgorithm           `json:"lb_algorithm,omitempty"`
	Protocol              LoadbalancerPoolProtocol        `json:"protocol" required:"true"`
	LoadbalancerID        string                          `json:"loadbalancer_id,omitempty"`
	ListenerID            string                          `json:"listener_id,omitempty"`
	TimeoutClientData     int                             `json:"timeout_client_data,omitempty"`
	TimeoutMemberData     int                             `json:"timeout_member_data,omitempty"`
	TimeoutMemberConnect  int                             `json:"timeout_member_connect,omitempty"`
	Members               []PoolMemberCreateRequest       `json:"members"`
	HealthMonitor         HealthMonitorCreateRequest      `json:"healthmonitor,omitempty"`
	SessionPersistence    *LoadbalancerSessionPersistence `json:"session_persistence,omitempty"`
}

// PoolMemberCreateRequest represents a request to create a Loadbalancer Pool PoolMember.
type PoolMemberCreateRequest struct {
	ProtocolPort int    `json:"protocol_port" required:"true"`
	Address      net.IP `json:"address" required:"true"`
	SubnetID     string `json:"subnet_id,omitempty"`
	InstanceID   string `json:"instance_id,omitempty"`
	Weight       int    `json:"weight,omitempty"`
	AdminStateUP bool   `json:"admin_state_up,omitempty"`
}

type HealthMonitorType string

const (
	HealthMonitorTypeHTTP       HealthMonitorType = "HTTP"
	HealthMonitorTypeHTTPS      HealthMonitorType = "HTTPS"
	HealthMonitorTypePING       HealthMonitorType = "PING"
	HealthMonitorTypeTCP        HealthMonitorType = "TCP"
	HealthMonitorTypeTLSHello   HealthMonitorType = "TLS-HELLO"
	HealthMonitorTypeUDPConnect HealthMonitorType = "UDP-CONNECT"
)

type HTTPMethod string

const (
	HTTPMethodCONNECT HTTPMethod = "CONNECT"
	HTTPMethodDELETE  HTTPMethod = "DELETE"
	HTTPMethodGET     HTTPMethod = "GET"
	HTTPMethodHEAD    HTTPMethod = "HEAD"
	HTTPMethodOPTIONS HTTPMethod = "OPTIONS"
	HTTPMethodPATCH   HTTPMethod = "PATCH"
	HTTPMethodPOST    HTTPMethod = "POST"
	HTTPMethodPUT     HTTPMethod = "PUT"
	HTTPMethodTRACE   HTTPMethod = "TRACE"
)

// HealthMonitorCreateRequest represents a request to create a Loadbalancer Pool Health Monitor.
type HealthMonitorCreateRequest struct {
	ID             string            `json:"id,omitempty"`
	Type           HealthMonitorType `json:"type" required:"true"`
	Delay          int               `json:"delay" required:"true"`
	MaxRetries     int               `json:"max_retries" required:"true"`
	Timeout        int               `json:"timeout" required:"true"`
	MaxRetriesDown int               `json:"max_retries_down,omitempty"`
	URLPath        string            `json:"url_path,omitempty"`
	ExpectedCodes  *string           `json:"expected_codes,omitempty"`
	HTTPMethod     *HTTPMethod       `json:"http_method,omitempty"`
}

// HealthMonitorUpdateRequest represents a request to update a Loadbalancer Pool Health Monitor.
type HealthMonitorUpdateRequest struct {
	ID             string            `json:"id,omitempty"`
	Type           HealthMonitorType `json:"type,omitempty"`
	Delay          int               `json:"delay" required:"true"`
	MaxRetries     int               `json:"max_retries" required:"true"`
	Timeout        int               `json:"timeout" required:"true"`
	MaxRetriesDown int               `json:"max_retries_down,omitempty"`
	HTTPMethod     *HTTPMethod       `json:"http_method,omitempty"`
	URLPath        *string           `json:"url_path,omitempty"`
	ExpectedCodes  *string           `json:"expected_codes,omitempty"`
}

// LoadbalancerListenerCreateRequest represents a request to create a Loadbalancer Listener.
// Used as part of a request to create a Loadbalancer.
type LoadbalancerListenerCreateRequest struct {
	Name             string                          `json:"name" required:"true" validate:"required,name"`
	Protocol         LoadbalancerListenerProtocol    `json:"protocol" required:"true"`
	ProtocolPort     int                             `json:"protocol_port" required:"true"`
	Certificate      string                          `json:"certificate,omitempty"`
	CertificateChain string                          `json:"certificate_chain,omitempty"`
	PrivateKey       string                          `json:"private_key,omitempty"`
	SecretID         string                          `json:"secret_id,omitempty"`
	InsertXForwarded bool                            `json:"insert_x_forwarded"`
	SNISecretID      []string                        `json:"sni_secret_id,omitempty"`
	Pools            []LoadbalancerPoolCreateRequest `json:"pools,omitempty" validate:"omitempty,dive"`
}

type LoadbalancerListenerProtocol string

const (
	ListenerProtocolHTTP            LoadbalancerListenerProtocol = "HTTP"
	ListenerProtocolHTTPS           LoadbalancerListenerProtocol = "HTTPS"
	ListenerProtocolTCP             LoadbalancerListenerProtocol = "TCP"
	ListenerProtocolUDP             LoadbalancerListenerProtocol = "UDP"
	ListenerProtocolTerminatedHTTPS LoadbalancerListenerProtocol = "TERMINATED_HTTPS"
)

type SessionPersistence string

const (
	SessionPersistenceAppCookie  SessionPersistence = "APP_COOKIE"
	SessionPersistenceHTTPCookie SessionPersistence = "HTTP_COOKIE"
	SessionPersistenceSourceIP   SessionPersistence = "SOURCE_IP"
)

// LoadbalancerSessionPersistence represents a request to create a Loadbalancer Pool Persistence Session.
type LoadbalancerSessionPersistence struct {
	Type                   SessionPersistence `json:"type" required:"true"`
	CookieName             string             `json:"cookie_name,omitempty"`
	PersistenceTimeout     int                `json:"persistence_timeout,omitempty"`
	PersistenceGranularity string             `json:"persistence_granularity,omitempty"`
}

// LoadbalancerListOptions specifies the optional query parameters to List method.
type LoadbalancerListOptions struct {
	ShowStats        bool   `url:"show_stats,omitempty"  validate:"omitempty"`
	AssignedFloating bool   `url:"assigned_floating,omitempty"  validate:"omitempty"`
	MetadataKV       string `url:"metadata_kv,omitempty"  validate:"omitempty"`
	MetadataK        string `url:"metadata_k,omitempty"  validate:"omitempty"`
}

// loadbalancersRoot represents a Loadbalancers root.
type loadbalancersRoot struct {
	Count         int
	Loadbalancers []Loadbalancer `json:"results"`
}

// loadbalancerRoot represents a Loadbalancer Pools root.
type loadbalancerPoolsRoot struct {
	Count int
	Pools []Pool `json:"results"`
}

// ListenerCreateRequest represents a request to create a Loadbalancer Listener.
// Used as a separate request to create Listener.
type ListenerCreateRequest struct {
	Name             string                       `json:"name" required:"true" validate:"required,name"`
	Protocol         LoadbalancerListenerProtocol `json:"protocol" required:"true"`
	ProtocolPort     int                          `json:"protocol_port" required:"true"`
	LoadBalancerID   string                       `json:"loadbalancer_id" required:"true"`
	InsertXForwarded bool                         `json:"insert_x_forwarded"`
	SecretID         string                       `json:"secret_id,omitempty"`
	SNISecretID      []string                     `json:"sni_secret_id,omitempty"`
}

// PoolCreateRequest represents a request to create a Loadbalancer Listener Pool.
// Used as a separate request to create Pool.
type PoolCreateRequest struct {
	LoadbalancerPoolCreateRequest
}

// PoolUpdateRequest represents a request to update a Loadbalancer Listener Pool.
type PoolUpdateRequest struct {
	ID                    string                          `json:"id,omitempty"`
	Name                  string                          `json:"name,omitempty"`
	LoadBalancerAlgorithm LoadBalancerAlgorithm           `json:"lb_algorithm,omitempty"`
	SessionPersistence    *LoadbalancerSessionPersistence `json:"session_persistence,omitempty"`
	Members               []PoolMember                    `json:"members,omitempty"`
	HealthMonitor         *HealthMonitorUpdateRequest     `json:"healthmonitor,omitempty"`
	TimeoutClientData     int                             `json:"timeout_client_data,omitempty"`
	TimeoutMemberData     int                             `json:"timeout_member_data,omitempty"`
	TimeoutMemberConnect  int                             `json:"timeout_member_connect,omitempty"`
}

type PoolListOptions struct {
	LoadBalancerID string `url:"loadbalancer_id,omitempty"`
	ListenerID     string `url:"listener_id,omitempty"`
	Details        bool   `url:"details,omitempty"` // if true Details show the member and healthmonitor details
}

// LoadbalancerCheckLimitsRequest represents a request to check the limits of a loadbalancer.
type LoadbalancerCheckLimitsRequest struct {
	FloatingIP InterfaceFloatingIP `json:"floating_ip,omitempty"`
}

// loadbalancerFlavorRoot represents a Loadbalancer Flavor root.
type loadbalancerFlavorRoot struct {
	Count   int
	Flavors []Flavor `json:"results"`
}

// List get load balancers.
func (s *LoadbalancersServiceOp) List(ctx context.Context, opts *LoadbalancerListOptions) ([]Loadbalancer, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(loadbalancersBasePathV1)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(loadbalancersRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Loadbalancers, resp, err
}

// Get individual Loadbalancer.
func (s *LoadbalancersServiceOp) Get(ctx context.Context, loadbalancerID string) (*Loadbalancer, *Response, error) {
	if resp, err := isValidUUID(loadbalancerID, "loadbalancerID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(loadbalancersBasePathV1), loadbalancerID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	loadBalancer := new(Loadbalancer)
	resp, err := s.client.Do(ctx, req, loadBalancer)
	if err != nil {
		return nil, resp, err
	}

	return loadBalancer, resp, err
}

// Create a Loadbalancer.
func (s *LoadbalancersServiceOp) Create(ctx context.Context, createRequest *LoadbalancerCreateRequest) (*TaskResponse, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(loadbalancersBasePathV1)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
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

// Delete the Loadbalancer.
func (s *LoadbalancersServiceOp) Delete(ctx context.Context, loadbalancerID string) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(loadbalancerID, "loadbalancerID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(loadbalancersBasePathV1), loadbalancerID)

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

// ListenerGet a Loadbalancer Listener.
func (s *LoadbalancersServiceOp) ListenerGet(ctx context.Context, listenerID string) (*Listener, *Response, error) {
	if resp, err := isValidUUID(listenerID, "listenerID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(lblistenersBasePathV1), listenerID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	listener := new(Listener)
	resp, err := s.client.Do(ctx, req, listener)
	if err != nil {
		return nil, resp, err
	}

	return listener, resp, err
}

// ListenerCreate a Loadbalancer Listener.
func (s *LoadbalancersServiceOp) ListenerCreate(ctx context.Context, createRequest *ListenerCreateRequest) (*TaskResponse, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(lblistenersBasePathV1)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
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

// ListenerDelete the Loadbalancer Listener.
func (s *LoadbalancersServiceOp) ListenerDelete(ctx context.Context, listenerID string) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(listenerID, "listenerID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(lblistenersBasePathV1), listenerID)

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

// PoolGet a Loadbalancer Pool.
func (s *LoadbalancersServiceOp) PoolGet(ctx context.Context, poolID string) (*Pool, *Response, error) {
	if resp, err := isValidUUID(poolID, "poolID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(lbpoolsBasePathV1), poolID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	pool := new(Pool)
	resp, err := s.client.Do(ctx, req, pool)
	if err != nil {
		return nil, resp, err
	}

	return pool, resp, err
}

// PoolCreate a Loadbalancer Pool.
func (s *LoadbalancersServiceOp) PoolCreate(ctx context.Context, createRequest *PoolCreateRequest) (*TaskResponse, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(lbpoolsBasePathV1)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
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

// PoolDelete the Loadbalancer Pool.
func (s *LoadbalancersServiceOp) PoolDelete(ctx context.Context, poolID string) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(poolID, "poolID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(lbpoolsBasePathV1), poolID)

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

// PoolUpdate a Loadbalancer Pool.
func (s *LoadbalancersServiceOp) PoolUpdate(ctx context.Context, poolID string, updateRequest *PoolUpdateRequest) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(poolID, "poolID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(lbpoolsBasePathV1), poolID)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
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

// PoolList get Loadbalancer Pools.
func (s *LoadbalancersServiceOp) PoolList(ctx context.Context, opts *PoolListOptions) ([]Pool, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(lbpoolsBasePathV1)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(loadbalancerPoolsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Pools, resp, err
}

// CheckLimits check a quota for load balancer creation.
func (s *LoadbalancersServiceOp) CheckLimits(ctx context.Context, checkLimitsRequest *LoadbalancerCheckLimitsRequest) (*map[string]int, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(loadbalancersBasePathV1)
	path = fmt.Sprintf("%s/%s", path, loadbalancersCheckLimitsPath)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, checkLimitsRequest)
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

// FlavorList get load balancer flavors.
func (s *LoadbalancersServiceOp) FlavorList(ctx context.Context, opts *FlavorsOptions) ([]Flavor, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(lbflavorsBasePathV1)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(loadbalancerFlavorRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Flavors, resp, err
}

// MetadataList load balancer detailed metadata items.
func (s *LoadbalancersServiceOp) MetadataList(ctx context.Context, loadbalancerID string) ([]MetadataDetailed, *Response, error) {
	if resp, err := isValidUUID(loadbalancerID, "loadbalancerID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	return metadataList(ctx, s.client, loadbalancerID, loadbalancersBasePathV1)
}

// MetadataCreate or update load balancer metadata.
func (s *LoadbalancersServiceOp) MetadataCreate(ctx context.Context, loadbalancerID string, metadata *MetadataCreateRequest) (*Response, error) {
	if resp, err := isValidUUID(loadbalancerID, "loadbalancerID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	return metadataCreate(ctx, s.client, loadbalancerID, loadbalancersBasePathV1, metadata)
}

// MetadataUpdate load balancer metadata.
func (s *LoadbalancersServiceOp) MetadataUpdate(ctx context.Context, loadbalancerID string, metadata *MetadataCreateRequest) (*Response, error) {
	if resp, err := isValidUUID(loadbalancerID, "loadbalancerID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	return metadataUpdate(ctx, s.client, loadbalancerID, loadbalancersBasePathV1, metadata)
}

// MetadataDeleteItem a load balancer metadata item by key.
func (s *LoadbalancersServiceOp) MetadataDeleteItem(ctx context.Context, loadbalancerID string, opts *MetadataItemOptions) (*Response, error) {
	if resp, err := isValidUUID(loadbalancerID, "loadbalancerID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	return metadataDeleteItem(ctx, s.client, loadbalancerID, loadbalancersBasePathV1, opts)
}

// MetadataGetItem load balancer detailed metadata.
func (s *LoadbalancersServiceOp) MetadataGetItem(ctx context.Context, loadbalancerID string, opts *MetadataItemOptions) (*MetadataDetailed, *Response, error) {
	if resp, err := isValidUUID(loadbalancerID, "loadbalancerID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	return metadataGetItem(ctx, s.client, loadbalancerID, loadbalancersBasePathV1, opts)
}
