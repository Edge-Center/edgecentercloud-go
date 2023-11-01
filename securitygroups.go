package edgecloud

import (
	"context"
	"fmt"
	"net/http"
)

const (
	securitygroupsBasePathV1      = "/v1/securitygroups"
	securitygroupsRulesBasePathV1 = "/v1/securitygrouprules"
	securitygroupsCopyPath        = "copy"
	securitygroupsRulesPath       = "rules"
)

// SecurityGroupsService is an interface for creating and managing Security Groups (Firewalls) with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/securitygroups
type SecurityGroupsService interface {
	List(context.Context, *SecurityGroupListOptions) ([]SecurityGroup, *Response, error)
	Get(context.Context, string) (*SecurityGroup, *Response, error)
	Create(context.Context, *SecurityGroupCreateRequest) (*TaskResponse, *Response, error)
	Delete(context.Context, string) (*TaskResponse, *Response, error)
	Update(context.Context, string, *SecurityGroupUpdateRequest) (*SecurityGroup, *Response, error)
	DeepCopy(context.Context, string, *SecurityGroupDeepCopyRequest) (*Response, error)

	SecurityGroupsRules
	SecurityGroupsMetadata
}

type SecurityGroupsRules interface {
	RuleCreate(context.Context, string, *RuleCreateRequest) (*SecurityGroupRule, *Response, error)
	RuleDelete(context.Context, string) (*TaskResponse, *Response, error)
	RuleUpdate(context.Context, string, *RuleUpdateRequest) (*SecurityGroupRule, *Response, error)
}

type SecurityGroupsMetadata interface {
	MetadataList(context.Context, string) ([]MetadataDetailed, *Response, error)
	MetadataCreate(context.Context, string, *MetadataCreateRequest) (*Response, error)
	MetadataUpdate(context.Context, string, *MetadataCreateRequest) (*Response, error)
	MetadataDeleteItem(context.Context, string, *MetadataItemOptions) (*Response, error)
	MetadataGetItem(context.Context, string, *MetadataItemOptions) (*MetadataDetailed, *Response, error)
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

// SecurityGroupListOptions specifies the optional query parameters to List method.
type SecurityGroupListOptions struct {
	MetadataKV string `url:"metadata_kv,omitempty"  validate:"omitempty"`
	MetadataK  string `url:"metadata_k,omitempty"  validate:"omitempty"`
}

// securityGroupsRoot represents a SecurityGroup root.
type securityGroupsRoot struct {
	Count          int
	SecurityGroups []SecurityGroup `json:"results"`
}

type SecurityGroupUpdateRequest struct {
	Name         string         `json:"name"`
	ChangedRules []ChangedRules `json:"changed_rules"`
}

type ChangedRuleAction string

const (
	ChangedRuleCreate ChangedRuleAction = "create"
	ChangedRuleDelete ChangedRuleAction = "delete"
)

type ChangedRules struct {
	Description         string                     `json:"description"`
	RemoteIPPrefix      string                     `json:"remote_ip_prefix,omitempty"`
	SecurityGroupRuleID string                     `json:"security_group_rule_id,omitempty"`
	PortRangeMax        int                        `json:"port_range_max,omitempty"`
	Protocol            SecurityGroupRuleProtocol  `json:"protocol,omitempty"`
	PortRangeMin        int                        `json:"port_range_min,omitempty"`
	EtherType           EtherType                  `json:"ethertype,omitempty" required:"true"`
	RemoteGroupID       string                     `json:"remote_group_id,omitempty"`
	Direction           SecurityGroupRuleDirection `json:"direction"`
	Action              ChangedRuleAction          `json:"action"`
}

type SecurityGroupDeepCopyRequest struct {
	Name string `json:"name"`
}

type RuleCreateRequest struct {
	Description         string                     `json:"description"`
	RemoteIPPrefix      string                     `json:"remote_ip_prefix,omitempty"`
	SecurityGroupRuleID string                     `json:"security_group_rule_id,omitempty"`
	PortRangeMax        int                        `json:"port_range_max,omitempty"`
	Protocol            SecurityGroupRuleProtocol  `json:"protocol,omitempty"`
	PortRangeMin        int                        `json:"port_range_min,omitempty"`
	EtherType           EtherType                  `json:"ethertype,omitempty" required:"true"`
	RemoteGroupID       string                     `json:"remote_group_id,omitempty"`
	Direction           SecurityGroupRuleDirection `json:"direction"`
}

type RuleUpdateRequest struct {
	ID                  string                     `json:"id"`
	Description         string                     `json:"description"`
	RemoteIPPrefix      string                     `json:"remote_ip_prefix,omitempty"`
	SecurityGroupRuleID string                     `json:"security_group_rule_id,omitempty"`
	PortRangeMax        int                        `json:"port_range_max,omitempty"`
	Protocol            SecurityGroupRuleProtocol  `json:"protocol,omitempty"`
	PortRangeMin        int                        `json:"port_range_min,omitempty"`
	EtherType           EtherType                  `json:"ethertype,omitempty" required:"true"`
	RemoteGroupID       string                     `json:"remote_group_id,omitempty"`
	Direction           SecurityGroupRuleDirection `json:"direction"`
}

// List get security groups.
func (s *SecurityGroupsServiceOp) List(ctx context.Context, opts *SecurityGroupListOptions) ([]SecurityGroup, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(securitygroupsBasePathV1)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(securityGroupsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.SecurityGroups, resp, err
}

// Get individual Security Group.
func (s *SecurityGroupsServiceOp) Get(ctx context.Context, securityGroupID string) (*SecurityGroup, *Response, error) {
	if resp, err := isValidUUID(securityGroupID, "securityGroupID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(securitygroupsBasePathV1), securityGroupID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	securityGroup := new(SecurityGroup)
	resp, err := s.client.Do(ctx, req, securityGroup)
	if err != nil {
		return nil, resp, err
	}

	return securityGroup, resp, err
}

// Create a Security Group.
func (s *SecurityGroupsServiceOp) Create(ctx context.Context, createRequest *SecurityGroupCreateRequest) (*TaskResponse, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(securitygroupsBasePathV1)

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

// Delete the Security Group.
func (s *SecurityGroupsServiceOp) Delete(ctx context.Context, securityGroupID string) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(securityGroupID, "securityGroupID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(securitygroupsBasePathV1), securityGroupID)

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

// Update a Security Group.
func (s *SecurityGroupsServiceOp) Update(ctx context.Context, securityGroupID string, updateRequest *SecurityGroupUpdateRequest) (*SecurityGroup, *Response, error) {
	if resp, err := isValidUUID(securityGroupID, "securityGroupID"); err != nil {
		return nil, resp, err
	}

	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(securitygroupsBasePathV1), securityGroupID)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	securityGroup := new(SecurityGroup)
	resp, err := s.client.Do(ctx, req, securityGroup)
	if err != nil {
		return nil, resp, err
	}

	return securityGroup, resp, err
}

// DeepCopy creates a deep copy of a security group.
func (s *SecurityGroupsServiceOp) DeepCopy(ctx context.Context, securityGroupID string, deepCopyRequest *SecurityGroupDeepCopyRequest) (*Response, error) {
	if resp, err := isValidUUID(securityGroupID, "securityGroupID"); err != nil {
		return resp, err
	}

	if deepCopyRequest == nil {
		return nil, NewArgError("deepCopyRequest", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	path := s.client.addProjectRegionPath(securitygroupsBasePathV1)
	path = fmt.Sprintf("%s/%s/%s", path, securityGroupID, securitygroupsCopyPath)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, deepCopyRequest)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// RuleCreate to a security group.
func (s *SecurityGroupsServiceOp) RuleCreate(ctx context.Context, securityGroupID string, ruleCreateRequest *RuleCreateRequest) (*SecurityGroupRule, *Response, error) {
	if resp, err := isValidUUID(securityGroupID, "securityGroupID"); err != nil {
		return nil, resp, err
	}

	if ruleCreateRequest == nil {
		return nil, nil, NewArgError("ruleCreateRequest", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(securitygroupsBasePathV1)
	path = fmt.Sprintf("%s/%s/%s", path, securityGroupID, securitygroupsRulesPath)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, ruleCreateRequest)
	if err != nil {
		return nil, nil, err
	}

	securityGroupRule := new(SecurityGroupRule)
	resp, err := s.client.Do(ctx, req, securityGroupRule)
	if err != nil {
		return nil, resp, err
	}

	return securityGroupRule, resp, err
}

// RuleDelete a security group rule.
func (s *SecurityGroupsServiceOp) RuleDelete(ctx context.Context, securityGroupID string) (*TaskResponse, *Response, error) {
	if resp, err := isValidUUID(securityGroupID, "securityGroupID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(securitygroupsRulesBasePathV1), securityGroupID)

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

// RuleUpdate a security group rule.
func (s *SecurityGroupsServiceOp) RuleUpdate(ctx context.Context, securityGroupID string, ruleUpdateRequest *RuleUpdateRequest) (*SecurityGroupRule, *Response, error) {
	if resp, err := isValidUUID(securityGroupID, "securityGroupID"); err != nil {
		return nil, resp, err
	}

	if ruleUpdateRequest == nil {
		return nil, nil, NewArgError("ruleUpdateRequest", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(securitygroupsRulesBasePathV1), securityGroupID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, ruleUpdateRequest)
	if err != nil {
		return nil, nil, err
	}

	securityGroupRule := new(SecurityGroupRule)
	resp, err := s.client.Do(ctx, req, securityGroupRule)
	if err != nil {
		return nil, resp, err
	}

	return securityGroupRule, resp, err
}

// MetadataList security group detailed metadata items.
func (s *SecurityGroupsServiceOp) MetadataList(ctx context.Context, securityGroupID string) ([]MetadataDetailed, *Response, error) {
	if resp, err := isValidUUID(securityGroupID, "securityGroupID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	return metadataList(ctx, s.client, securityGroupID, securitygroupsBasePathV1)
}

// MetadataCreate or update security group metadata.
func (s *SecurityGroupsServiceOp) MetadataCreate(ctx context.Context, securityGroupID string, metadata *MetadataCreateRequest) (*Response, error) {
	if resp, err := isValidUUID(securityGroupID, "securityGroupID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	return metadataCreate(ctx, s.client, securityGroupID, securitygroupsBasePathV1, metadata)
}

// MetadataUpdate security group metadata.
func (s *SecurityGroupsServiceOp) MetadataUpdate(ctx context.Context, securityGroupID string, metadata *MetadataCreateRequest) (*Response, error) {
	if resp, err := isValidUUID(securityGroupID, "securityGroupID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	return metadataUpdate(ctx, s.client, securityGroupID, securitygroupsBasePathV1, metadata)
}

// MetadataDeleteItem a security group metadata item by key.
func (s *SecurityGroupsServiceOp) MetadataDeleteItem(ctx context.Context, securityGroupID string, opts *MetadataItemOptions) (*Response, error) {
	if resp, err := isValidUUID(securityGroupID, "securityGroupID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	return metadataDeleteItem(ctx, s.client, securityGroupID, securitygroupsBasePathV1, opts)
}

// MetadataGetItem security group detailed metadata.
func (s *SecurityGroupsServiceOp) MetadataGetItem(ctx context.Context, securityGroupID string, opts *MetadataItemOptions) (*MetadataDetailed, *Response, error) {
	if resp, err := isValidUUID(securityGroupID, "securityGroupID"); err != nil {
		return nil, resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	return metadataGetItem(ctx, s.client, securityGroupID, securitygroupsBasePathV1, opts)
}
