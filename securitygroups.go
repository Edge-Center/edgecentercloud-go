package edgecloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

const (
	securitygroupsBasePathV1 = "/v1/securitygroups"
)

// SecurityGroupsService is an interface for creating and managing Security Groups (Firewalls) with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/securitygroups
type SecurityGroupsService interface {
	Get(context.Context, string, *ServicePath) (*SecurityGroup, *Response, error)
	Create(context.Context, *SecurityGroupCreateRequest, *ServicePath) (*TaskResponse, *Response, error)
	Delete(context.Context, string, *ServicePath) (*TaskResponse, *Response, error)
}

// SecurityGroupsServiceOp handles communication with Security Groups (Firewalls) methods of the EdgecenterCloud API.
type SecurityGroupsServiceOp struct {
	client *Client
}

var _ SecurityGroupsService = &SecurityGroupsServiceOp{}

// SecurityGroup represents a EdgecenterCloud Security Group.
type SecurityGroup struct {
	ID                 string              `json:"id"`
	CreatedAt          string              `json:"created_at"`
	UpdatedAt          string              `json:"updated_at"`
	RevisionNumber     int                 `json:"revision_number"`
	Name               string              `json:"name"`
	Description        string              `json:"description"`
	SecurityGroupRules []SecurityGroupRule `json:"security_group_rules"`
	Metadata           []MetadataDetailed  `json:"metadata"`
	ProjectID          int                 `json:"project_id"`
	RegionID           int                 `json:"region_id"`
	Region             string              `json:"region"`
	Tags               []string            `json:"tags"`
}

// SecurityGroupRule represents a EdgecenterCloud Security Group Rule.
type SecurityGroupRule struct {
	ID              string                     `json:"id"`
	SecurityGroupID string                     `json:"security_group_id"`
	RemoteGroupID   string                     `json:"remote_group_id"`
	Direction       SecurityGroupRuleDirection `json:"direction"`
	EtherType       EtherType                  `json:"ethertype"`
	Protocol        SecurityGroupRuleProtocol  `json:"protocol"`
	PortRangeMax    int                        `json:"port_range_max"`
	PortRangeMin    int                        `json:"port_range_min"`
	Description     string                     `json:"description"`
	RemoteIPPrefix  string                     `json:"remote_ip_prefix"`
	CreatedAt       string                     `json:"created_at"`
	UpdatedAt       string                     `json:"updated_at"`
	RevisionNumber  int                        `json:"revision_number"`
}

// SecurityGroupCreateRequest represents a request to create a Security Group.
type SecurityGroupCreateRequest struct {
	SecurityGroup SecurityGroup                  `json:"security_group" required:"true"`
	Instances     []InstanceSecurityGroupsCreate `json:"instances"`
}

// SecurityGroupRuleCreateRequest represents a request to create a Security Group Rule.
type SecurityGroupRuleCreateRequest struct {
	EtherType       EtherType                  `json:"ethertype,omitempty" required:"true"`
	Description     string                     `json:"description,omitempty"`
	RemoteGroupID   string                     `json:"remote_group_id,omitempty"`
	PortRangeMin    int                        `json:"port_range_min,omitempty"`
	PortRangeMax    int                        `json:"port_range_max,omitempty"`
	RemoteIPPrefix  string                     `json:"remote_ip_prefix,omitempty"`
	Protocol        SecurityGroupRuleProtocol  `json:"protocol,omitempty" required:"true"`
	Direction       SecurityGroupRuleDirection `json:"direction" required:"true"`
	SecurityGroupID *string                    `json:"security_group_id,omitempty"`
}

type EtherType string

const (
	EtherTypeIPv4 EtherType = "IPv4"
	EtherTypeIPv6 EtherType = "IPv6"
)

type SecurityGroupRuleProtocol string

const (
	SGRuleProtocolANY     SecurityGroupRuleProtocol = "any"
	SGRuleProtocolAH      SecurityGroupRuleProtocol = "ah"
	SGRuleProtocolACCP    SecurityGroupRuleProtocol = "dccp"
	SGRuleProtocolEGP     SecurityGroupRuleProtocol = "egp"
	SGRuleProtocolESP     SecurityGroupRuleProtocol = "esp"
	SGRuleProtocolGRE     SecurityGroupRuleProtocol = "gre"
	SGRuleProtocolICMP    SecurityGroupRuleProtocol = "icmp"
	SGRuleProtocolIGMP    SecurityGroupRuleProtocol = "igmp"
	SGRuleProtocolIPIP    SecurityGroupRuleProtocol = "ipip"
	SGRuleProtocolOSPF    SecurityGroupRuleProtocol = "ospf"
	SGRuleProtocolPGM     SecurityGroupRuleProtocol = "pgm"
	SGRuleProtocolRSVP    SecurityGroupRuleProtocol = "rsvp"
	SGRuleProtocolSCTP    SecurityGroupRuleProtocol = "sctp"
	SGRuleProtocolTCP     SecurityGroupRuleProtocol = "tcp"
	SGRuleProtocolUDP     SecurityGroupRuleProtocol = "udp"
	SGRuleProtocolUDPLITE SecurityGroupRuleProtocol = "udplite"
	SGRuleProtocolVRRP    SecurityGroupRuleProtocol = "vrrp"
)

type SecurityGroupRuleDirection string

const (
	SGRuleDirectionEgress  SecurityGroupRuleDirection = "egress"
	SGRuleDirectionIngress SecurityGroupRuleDirection = "ingress"
)

// securityGroupRoot represents a Security Group root.
type securityGroupRoot struct {
	SecurityGroup *SecurityGroup `json:"security_group"`
	Tasks         *TaskResponse  `json:"tasks"`
}

// Get individual Security Group.
func (s *SecurityGroupsServiceOp) Get(ctx context.Context, sgID string, p *ServicePath) (*SecurityGroup, *Response, error) {
	if _, err := uuid.Parse(sgID); err != nil {
		return nil, nil, NewArgError("sgID", "should be the correct UUID")
	}

	if p == nil {
		return nil, nil, NewArgError("ServicePath", "cannot be nil")
	}

	path := addServicePath(securitygroupsBasePathV1, p)
	path = fmt.Sprintf("%s/%s", path, sgID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(securityGroupRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.SecurityGroup, resp, err
}

// Create a Security Group.
func (s *SecurityGroupsServiceOp) Create(ctx context.Context, createRequest *SecurityGroupCreateRequest, p *ServicePath) (*TaskResponse, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	if p == nil {
		return nil, nil, NewArgError("ServicePath", "cannot be nil")
	}

	path := addServicePath(securitygroupsBasePathV1, p)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(securityGroupRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Tasks, resp, err
}

// Delete the Security Group.
func (s *SecurityGroupsServiceOp) Delete(ctx context.Context, sgID string, p *ServicePath) (*TaskResponse, *Response, error) {
	if _, err := uuid.Parse(sgID); err != nil {
		return nil, nil, NewArgError("sgID", "should be the correct UUID")
	}

	if p == nil {
		return nil, nil, NewArgError("ServicePath", "cannot be nil")
	}

	path := addServicePath(securitygroupsBasePathV1, p)
	path = fmt.Sprintf("%s/%s", path, sgID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(securityGroupRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Tasks, resp, err
}
