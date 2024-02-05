package edgecloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	securitygroupsBasePathV1      = "/v1/securitygroups"
	securitygroupsRulesBasePathV1 = "/v1/securitygrouprules"
)

const (
	securitygroupsCopy  = "copy"
	securitygroupsRules = "rules"
)

var (
	ErrSGEtherTypeNotAllowed     = fmt.Errorf("invalid EtherType, allowed only %s or %s", EtherTypeIPv4, EtherTypeIPv6)
	ErrSGInvalidProtocol         = fmt.Errorf("invalid Protocol")
	ErrSGRuleDirectionNotAllowed = fmt.Errorf("invalid RuleDirection type, allowed only %s or %s", SGRuleDirectionIngress, SGRuleDirectionEgress)
)

// SecurityGroupsService is an interface for creating and managing Security Groups (Firewalls) with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/securitygroups
type SecurityGroupsService interface {
	List(context.Context, *SecurityGroupListOptions) ([]SecurityGroup, *Response, error)
	Get(context.Context, string) (*SecurityGroup, *Response, error)
	Create(context.Context, *SecurityGroupCreateRequest) (*SecurityGroup, *Response, error)
	Delete(context.Context, string) (*Response, error)
	Update(context.Context, string, *SecurityGroupUpdateRequest) (*SecurityGroup, *Response, error)
	DeepCopy(context.Context, string, *Name) (*Response, error)

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
	MetadataCreate(context.Context, string, *Metadata) (*Response, error)
	MetadataUpdate(context.Context, string, *Metadata) (*Response, error)
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
	EtherType       *EtherType                 `json:"ethertype"`
	Protocol        *SecurityGroupRuleProtocol `json:"protocol"`
	PortRangeMax    *int                       `json:"port_range_max"`
	PortRangeMin    *int                       `json:"port_range_min"`
	Description     *string                    `json:"description"`
	RemoteIPPrefix  *string                    `json:"remote_ip_prefix"`
	CreatedAt       string                     `json:"created_at"`
	UpdatedAt       string                     `json:"updated_at"`
	RevisionNumber  int                        `json:"revision_number"`
}

// SecurityGroupCreateRequest represents a request to create a Security Group.
type SecurityGroupCreateRequest struct {
	SecurityGroup SecurityGroupCreateRequestInner `json:"security_group" required:"true"`
	Instances     []ID                            `json:"instances,omitempty"`
}

type SecurityGroupCreateRequestInner struct {
	Name               string              `json:"name" required:"true"`
	Description        *string             `json:"description,omitempty"`
	Metadata           Metadata            `json:"metadata,omitempty"`
	Tags               []string            `json:"tags,omitempty"`
	SecurityGroupRules []RuleCreateRequest `json:"security_group_rules,omitempty"`
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
	SGRuleProtocolDCCP    SecurityGroupRuleProtocol = "dccp"
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
	SGRuleProtocolIPEncap SecurityGroupRuleProtocol = "ipencap"
)

type SecurityGroupRuleDirection string

const (
	SGRuleDirectionEgress  SecurityGroupRuleDirection = "egress"
	SGRuleDirectionIngress SecurityGroupRuleDirection = "ingress"
)

func (rd SecurityGroupRuleDirection) IsValid() error {
	switch rd {
	case SGRuleDirectionEgress,
		SGRuleDirectionIngress:
		return nil
	}

	return ErrSGRuleDirectionNotAllowed
}

func (rd SecurityGroupRuleDirection) ValidOrNil() (*SecurityGroupRuleDirection, error) {
	if rd.String() == "" {
		return nil, nil //nolint:nilnil
	}
	err := rd.IsValid()
	if err != nil {
		return &rd, err
	}

	return &rd, nil
}

func (rd SecurityGroupRuleDirection) String() string {
	return string(rd)
}

func (rd SecurityGroupRuleDirection) List() []SecurityGroupRuleDirection {
	return []SecurityGroupRuleDirection{
		SGRuleDirectionIngress,
		SGRuleDirectionEgress,
	}
}

func (rd SecurityGroupRuleDirection) StringList() []string {
	lst := rd.List()
	strings := make([]string, 0, len(lst))
	for _, x := range lst {
		strings = append(strings, x.String())
	}

	return strings
}

// UnmarshalJSON - implements Unmarshaler interface.
func (rd *SecurityGroupRuleDirection) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := SecurityGroupRuleDirection(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*rd = v

	return nil
}

// MarshalJSON - implements Marshaler interface.
func (rd *SecurityGroupRuleDirection) MarshalJSON() ([]byte, error) {
	return json.Marshal(rd.String())
}

func (et EtherType) IsValid() error {
	switch et {
	case EtherTypeIPv4, EtherTypeIPv6:
		return nil
	}
	return ErrSGEtherTypeNotAllowed
}

func (et EtherType) ValidOrNil() (*EtherType, error) {
	if et.String() == "" {
		return nil, nil //nolint:nilnil
	}
	err := et.IsValid()
	if err != nil {
		return &et, err
	}

	return &et, nil
}

func (et EtherType) String() string {
	return string(et)
}

func (et EtherType) List() []EtherType {
	return []EtherType{
		EtherTypeIPv4,
		EtherTypeIPv6,
	}
}

func (et EtherType) StringList() []string {
	lst := et.List()
	strings := make([]string, 0, len(lst))
	for _, x := range lst {
		strings = append(strings, x.String())
	}

	return strings
}

// UnmarshalJSON - implements Unmarshaler interface.
func (et *EtherType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := EtherType(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*et = v

	return nil
}

// MarshalJSON - implements Marshaler interface.
func (et *EtherType) MarshalJSON() ([]byte, error) {
	return json.Marshal(et.String())
}

func (p SecurityGroupRuleProtocol) IsValid() error {
	switch p {
	case SGRuleProtocolUDP,
		SGRuleProtocolTCP,
		SGRuleProtocolANY,
		SGRuleProtocolICMP,
		SGRuleProtocolAH,
		SGRuleProtocolDCCP,
		SGRuleProtocolEGP,
		SGRuleProtocolESP,
		SGRuleProtocolGRE,
		SGRuleProtocolIGMP,
		SGRuleProtocolOSPF,
		SGRuleProtocolPGM,
		SGRuleProtocolRSVP,
		SGRuleProtocolSCTP,
		SGRuleProtocolUDPLITE,
		SGRuleProtocolVRRP,
		SGRuleProtocolIPIP,
		SGRuleProtocolIPEncap:
		return nil
	}

	return ErrSGInvalidProtocol
}

func (p SecurityGroupRuleProtocol) ValidOrNil() (*SecurityGroupRuleProtocol, error) {
	if p.String() == "" {
		return nil, nil //nolint:nilnil
	}
	err := p.IsValid()
	if err != nil {
		return &p, err
	}

	return &p, nil
}

func (p SecurityGroupRuleProtocol) String() string {
	return string(p)
}

func (p SecurityGroupRuleProtocol) List() []SecurityGroupRuleProtocol {
	return []SecurityGroupRuleProtocol{
		SGRuleProtocolUDP,
		SGRuleProtocolTCP,
		SGRuleProtocolANY,
		SGRuleProtocolICMP,
		SGRuleProtocolAH,
		SGRuleProtocolDCCP,
		SGRuleProtocolEGP,
		SGRuleProtocolESP,
		SGRuleProtocolGRE,
		SGRuleProtocolIGMP,
		SGRuleProtocolOSPF,
		SGRuleProtocolPGM,
		SGRuleProtocolRSVP,
		SGRuleProtocolSCTP,
		SGRuleProtocolUDPLITE,
		SGRuleProtocolVRRP,
		SGRuleProtocolIPIP,
		SGRuleProtocolIPEncap,
	}
}

func (p SecurityGroupRuleProtocol) StringList() []string {
	lst := p.List()
	strings := make([]string, 0, len(lst))
	for _, x := range lst {
		strings = append(strings, x.String())
	}

	return strings
}

// UnmarshalJSON - implements Unmarshaler interface.
func (p *SecurityGroupRuleProtocol) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v := SecurityGroupRuleProtocol(s)
	err := v.IsValid()
	if err != nil {
		return err
	}
	*p = v

	return nil
}

// MarshalJSON - implements Marshaler interface.
func (p *SecurityGroupRuleProtocol) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

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

type RuleCreateRequest struct {
	Description     *string                    `json:"description"`
	RemoteIPPrefix  *string                    `json:"remote_ip_prefix,omitempty"`
	SecurityGroupID *string                    `json:"security_group_id,omitempty"`
	PortRangeMax    *int                       `json:"port_range_max,omitempty"`
	Protocol        SecurityGroupRuleProtocol  `json:"protocol,omitempty"`
	PortRangeMin    *int                       `json:"port_range_min,omitempty"`
	EtherType       EtherType                  `json:"ethertype,omitempty" required:"true"`
	RemoteGroupID   *string                    `json:"remote_group_id,omitempty"`
	Direction       SecurityGroupRuleDirection `json:"direction"`
}

type RuleUpdateRequest struct {
	ID              string                     `json:"id"`
	Description     string                     `json:"description"`
	RemoteIPPrefix  string                     `json:"remote_ip_prefix,omitempty"`
	SecurityGroupID string                     `json:"security_group_id,omitempty"`
	PortRangeMax    int                        `json:"port_range_max,omitempty"`
	Protocol        SecurityGroupRuleProtocol  `json:"protocol,omitempty"`
	PortRangeMin    int                        `json:"port_range_min,omitempty"`
	EtherType       EtherType                  `json:"ethertype,omitempty" required:"true"`
	RemoteGroupID   string                     `json:"remote_group_id,omitempty"`
	Direction       SecurityGroupRuleDirection `json:"direction"`
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
func (s *SecurityGroupsServiceOp) Create(ctx context.Context, reqBody *SecurityGroupCreateRequest) (*SecurityGroup, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(securitygroupsBasePathV1)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
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

// Delete the Security Group.
func (s *SecurityGroupsServiceOp) Delete(ctx context.Context, securityGroupID string) (*Response, error) {
	if resp, err := isValidUUID(securityGroupID, "securityGroupID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(securitygroupsBasePathV1), securityGroupID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Update a Security Group.
func (s *SecurityGroupsServiceOp) Update(ctx context.Context, securityGroupID string, reqBody *SecurityGroupUpdateRequest) (*SecurityGroup, *Response, error) {
	if resp, err := isValidUUID(securityGroupID, "securityGroupID"); err != nil {
		return nil, resp, err
	}

	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(securitygroupsBasePathV1), securityGroupID)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, reqBody)
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
func (s *SecurityGroupsServiceOp) DeepCopy(ctx context.Context, securityGroupID string, reqBody *Name) (*Response, error) {
	if resp, err := isValidUUID(securityGroupID, "securityGroupID"); err != nil {
		return resp, err
	}

	if reqBody == nil {
		return nil, NewArgError("deepCopyRequest", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	path := s.client.addProjectRegionPath(securitygroupsBasePathV1)
	path = fmt.Sprintf("%s/%s/%s", path, securityGroupID, securitygroupsCopy)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// RuleCreate to a security group.
func (s *SecurityGroupsServiceOp) RuleCreate(ctx context.Context, securityGroupID string, reqBody *RuleCreateRequest) (*SecurityGroupRule, *Response, error) {
	if resp, err := isValidUUID(securityGroupID, "securityGroupID"); err != nil {
		return nil, resp, err
	}

	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(securitygroupsBasePathV1)
	path = fmt.Sprintf("%s/%s/%s", path, securityGroupID, securitygroupsRules)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
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
// todo cloud-api deletes rule without tash.
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
func (s *SecurityGroupsServiceOp) RuleUpdate(ctx context.Context, securityGroupID string, reqBody *RuleUpdateRequest) (*SecurityGroupRule, *Response, error) {
	if resp, err := isValidUUID(securityGroupID, "securityGroupID"); err != nil {
		return nil, resp, err
	}

	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}

	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(securitygroupsRulesBasePathV1), securityGroupID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, reqBody)
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
func (s *SecurityGroupsServiceOp) MetadataCreate(ctx context.Context, securityGroupID string, reqBody *Metadata) (*Response, error) {
	if resp, err := isValidUUID(securityGroupID, "securityGroupID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	return metadataCreate(ctx, s.client, securityGroupID, securitygroupsBasePathV1, reqBody)
}

// MetadataUpdate security group metadata.
func (s *SecurityGroupsServiceOp) MetadataUpdate(ctx context.Context, securityGroupID string, reqBody *Metadata) (*Response, error) {
	if resp, err := isValidUUID(securityGroupID, "securityGroupID"); err != nil {
		return resp, err
	}

	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	return metadataUpdate(ctx, s.client, securityGroupID, securitygroupsBasePathV1, reqBody)
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
