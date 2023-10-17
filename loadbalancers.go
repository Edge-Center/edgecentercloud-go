package edgecloud

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/google/uuid"
)

const (
	loadbalancersBasePathV1 = "/v1/loadbalancers"
	lblistenersBasePathV1   = "/v1/lblisteners"
	lbpoolsBasePathV1       = "/v1/lbpools"
)

// LoadbalancersService is an interface for creating and managing Loadbalancer with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/loadbalancers
type LoadbalancersService interface {
	Get(context.Context, string, *ServicePath) (*Loadbalancer, *Response, error)
	Create(context.Context, *LoadbalancerCreateRequest, *ServicePath) (*TaskResponse, *Response, error)
	Delete(context.Context, string, *ServicePath) (*TaskResponse, *Response, error)
	LoadbalancerListeners
	LoadbalancerPools
}

type LoadbalancerListeners interface {
	ListenerGet(context.Context, string, *ServicePath) (*Listener, *Response, error)
	ListenerCreate(context.Context, *ListenerCreateRequest, *ServicePath) (*TaskResponse, *Response, error)
	ListenerDelete(context.Context, string, *ServicePath) (*TaskResponse, *Response, error)
}

type LoadbalancerPools interface {
	PoolGet(context.Context, string, *ServicePath) (*Pool, *Response, error)
	PoolCreate(context.Context, *PoolCreateRequest, *ServicePath) (*TaskResponse, *Response, error)
	PoolDelete(context.Context, string, *ServicePath) (*TaskResponse, *Response, error)
	PoolUpdate(context.Context, string, *PoolUpdateRequest, *ServicePath) (*TaskResponse, *Response, error)
	PoolList(context.Context, *ServicePath, *PoolListOptions) ([]Pool, *Response, error)
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
	Members               []Member                 `json:"members"`
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

// Member represents an EdgecenterCloud Loadbalancer Pool Member.
type Member struct {
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
	FloatingIP   InterfaceFloatingIP                 `json:"floating_ip,omitempty"`
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

// PoolMemberCreateRequest represents a request to create a Loadbalancer Pool Member.
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

// loadbalancerRoot represents a Loadbalancer root.
type loadbalancerRoot struct {
	Loadbalancer *Loadbalancer `json:"loadbalancer"`
	Listener     *Listener     `json:"listener"`
	Pool         *Pool         `json:"pool"`
	Tasks        *TaskResponse `json:"tasks"`
}

// loadbalancerRoot represents a Loadbalancer Pools root.
type loadbalancerPoolsRoot struct {
	Pools []Pool `json:"pools"`
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
	Members               []Member                        `json:"members,omitempty"`
	HealthMonitor         HealthMonitorCreateRequest      `json:"healthmonitor,omitempty"`
	TimeoutClientData     int                             `json:"timeout_client_data,omitempty"`
	TimeoutMemberData     int                             `json:"timeout_member_data,omitempty"`
	TimeoutMemberConnect  int                             `json:"timeout_member_connect,omitempty"`
}

type PoolListOptions struct {
	LoadBalancerID string `url:"loadbalancer_id,omitempty"`
	ListenerID     string `url:"listener_id,omitempty"`
	Details        bool   `url:"details,omitempty"` // if true Details show the member and healthmonitor details
}

// Get individual Loadbalancer.
func (s *LoadbalancersServiceOp) Get(ctx context.Context, loadbalancerID string, p *ServicePath) (*Loadbalancer, *Response, error) {
	if _, err := uuid.Parse(loadbalancerID); err != nil {
		return nil, nil, NewArgError("loadbalancerID", "should be the correct UUID")
	}

	if p == nil {
		return nil, nil, NewArgError("ServicePath", "cannot be nil")
	}

	path := addServicePath(loadbalancersBasePathV1, p)
	path = fmt.Sprintf("%s/%s", path, loadbalancerID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(loadbalancerRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Loadbalancer, resp, err
}

// Create a Loadbalancer.
func (s *LoadbalancersServiceOp) Create(ctx context.Context, createRequest *LoadbalancerCreateRequest, p *ServicePath) (*TaskResponse, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	if p == nil {
		return nil, nil, NewArgError("ServicePath", "cannot be nil")
	}

	path := addServicePath(loadbalancersBasePathV1, p)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(loadbalancerRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Tasks, resp, err
}

// Delete the Loadbalancer.
func (s *LoadbalancersServiceOp) Delete(ctx context.Context, loadbalancerID string, p *ServicePath) (*TaskResponse, *Response, error) {
	if _, err := uuid.Parse(loadbalancerID); err != nil {
		return nil, nil, NewArgError("loadbalancerID", "should be the correct UUID")
	}

	if p == nil {
		return nil, nil, NewArgError("ServicePath", "cannot be nil")
	}

	path := addServicePath(loadbalancersBasePathV1, p)
	path = fmt.Sprintf("%s/%s", path, loadbalancerID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(loadbalancerRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Tasks, resp, err
}

// ListenerGet a Loadbalancer Listener.
func (s *LoadbalancersServiceOp) ListenerGet(ctx context.Context, listenerID string, p *ServicePath) (*Listener, *Response, error) {
	if _, err := uuid.Parse(listenerID); err != nil {
		return nil, nil, NewArgError("listenerID", "should be the correct UUID")
	}

	if p == nil {
		return nil, nil, NewArgError("ServicePath", "cannot be nil")
	}

	path := addServicePath(lblistenersBasePathV1, p)
	path = fmt.Sprintf("%s/%s", path, listenerID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(loadbalancerRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Listener, resp, err
}

// ListenerCreate a Loadbalancer Listener.
func (s *LoadbalancersServiceOp) ListenerCreate(ctx context.Context, createRequest *ListenerCreateRequest, p *ServicePath) (*TaskResponse, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	if p == nil {
		return nil, nil, NewArgError("ServicePath", "cannot be nil")
	}

	path := addServicePath(lblistenersBasePathV1, p)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(loadbalancerRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Tasks, resp, err
}

// ListenerDelete the Loadbalancer Listener.
func (s *LoadbalancersServiceOp) ListenerDelete(ctx context.Context, listenerID string, p *ServicePath) (*TaskResponse, *Response, error) {
	if _, err := uuid.Parse(listenerID); err != nil {
		return nil, nil, NewArgError("listenerID", "should be the correct UUID")
	}

	if p == nil {
		return nil, nil, NewArgError("ServicePath", "cannot be nil")
	}

	path := addServicePath(lblistenersBasePathV1, p)
	path = fmt.Sprintf("%s/%s", path, listenerID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(loadbalancerRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Tasks, resp, err
}

// PoolGet a Loadbalancer Pool.
func (s *LoadbalancersServiceOp) PoolGet(ctx context.Context, poolID string, p *ServicePath) (*Pool, *Response, error) {
	if _, err := uuid.Parse(poolID); err != nil {
		return nil, nil, NewArgError("poolID", "should be the correct UUID")
	}

	if p == nil {
		return nil, nil, NewArgError("ServicePath", "cannot be nil")
	}

	path := addServicePath(lbpoolsBasePathV1, p)
	path = fmt.Sprintf("%s/%s", path, poolID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(loadbalancerRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Pool, resp, err
}

// PoolCreate a Loadbalancer Pool.
func (s *LoadbalancersServiceOp) PoolCreate(ctx context.Context, createRequest *PoolCreateRequest, p *ServicePath) (*TaskResponse, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	if p == nil {
		return nil, nil, NewArgError("ServicePath", "cannot be nil")
	}

	path := addServicePath(lbpoolsBasePathV1, p)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(loadbalancerRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Tasks, resp, err
}

// PoolDelete the Loadbalancer Pool.
func (s *LoadbalancersServiceOp) PoolDelete(ctx context.Context, poolID string, p *ServicePath) (*TaskResponse, *Response, error) {
	if _, err := uuid.Parse(poolID); err != nil {
		return nil, nil, NewArgError("poolID", "should be the correct UUID")
	}

	if p == nil {
		return nil, nil, NewArgError("ServicePath", "cannot be nil")
	}

	path := addServicePath(lbpoolsBasePathV1, p)
	path = fmt.Sprintf("%s/%s", path, poolID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(loadbalancerRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Tasks, resp, err
}

// PoolUpdate a Loadbalancer Pool.
func (s *LoadbalancersServiceOp) PoolUpdate(ctx context.Context, poolID string, updateRequest *PoolUpdateRequest, p *ServicePath) (*TaskResponse, *Response, error) {
	if _, err := uuid.Parse(poolID); err != nil {
		return nil, nil, NewArgError("poolID", "should be the correct UUID")
	}

	if p == nil {
		return nil, nil, NewArgError("ServicePath", "cannot be nil")
	}

	path := addServicePath(lbpoolsBasePathV1, p)
	path = fmt.Sprintf("%s/%s", path, poolID)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(loadbalancerRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Tasks, resp, err
}

// PoolList get Loadbalancer Pools.
func (s *LoadbalancersServiceOp) PoolList(ctx context.Context, p *ServicePath, opts *PoolListOptions) ([]Pool, *Response, error) {
	if p == nil {
		return nil, nil, NewArgError("ServicePath", "cannot be nil")
	}

	path := addServicePath(lbpoolsBasePathV1, p)
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
