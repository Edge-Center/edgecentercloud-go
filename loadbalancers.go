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
)

// LoadbalancersService is an interface for creating and managing Loadbalancer with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/loadbalancers
type LoadbalancersService interface {
	Get(context.Context, string, *ServicePath) (*Loadbalancer, *Response, error)
	Create(context.Context, *LoadbalancerCreateRequest, *ServicePath) (*TaskResponse, *Response, error)
	Delete(context.Context, string, *ServicePath) (*TaskResponse, *Response, error)
}

// LoadbalancersServiceOp handles communication with Loadbalancers methods of the EdgecenterCloud API.
type LoadbalancersServiceOp struct {
	client *Client
}

var _ LoadbalancersService = &LoadbalancersServiceOp{}

// Loadbalancer represents an EdgecenterCloud Loadbalancer.
type Loadbalancer struct {
	ID                 string                 `json:"id"`
	Name               string                 `json:"name"`
	Flavor             Flavor                 `json:"flavor"`
	VipAddress         net.IP                 `json:"vip_address"`
	VipPortID          string                 `json:"vip_port_id"`
	VipNetworkID       string                 `json:"vip_network_id"`
	ProvisioningStatus ProvisioningStatus     `json:"provisioning_status"`
	OperationStatus    OperatingStatus        `json:"operating_status"`
	CreatedAt          string                 `json:"created_at"`
	UpdatedAt          string                 `json:"updated_at"`
	CreatorTaskID      string                 `json:"creator_task_id"`
	TaskID             string                 `json:"task_id"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
	Stats              LoadbalancerStats      `json:"stats"`
	Listeners          []Listener             `json:"listeners"`
	FloatingIPs        []FloatingIP           `json:"floating_ips"`
	VrrpIPs            []VrrpIP               `json:"vrrp_ips"`
	ProjectID          int                    `json:"project_id"`
	RegionID           int                    `json:"region_id"`
	Region             string                 `json:"region"`
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
	Metadata     map[string]interface{}              `json:"metadata,omitempty" validate:"omitempty,dive"`
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

	HTTPMethod *HTTPMethod `json:"http_method,omitempty"`
}

// LoadbalancerListenerCreateRequest represents a request to create a Loadbalancer Listener.
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
	Tasks        *TaskResponse `json:"tasks"`
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
