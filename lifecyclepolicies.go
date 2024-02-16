package edgecloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ladydascalie/currency"
	"github.com/shopspring/decimal"
)

type (
	LifeCyclePolicyScheduleType string
	LifeCyclePolicyStatus       string
	LifeCyclePolicyAction       string
)

const (
	lifecyclePoliciesBasePathV1                                     = "/v1/lifecycle_policies/"
	addSchedulesSubPath                                             = "add_schedules"
	addVolumesSubPath                                               = "add_volumes_to_policy"
	removeSchedulesSubPath                                          = "remove_schedules"
	removeVolumesSubPath                                            = "remove_volumes_from_policy"
	estimateMaxPolicyUsageSubPath                                   = "estimate_max_policy_usage"
	LifeCyclePolicyScheduleTypeCron     LifeCyclePolicyScheduleType = "cron"
	LifeCyclePolicyScheduleTypeInterval LifeCyclePolicyScheduleType = "interval"
	LifeCyclePolicyStatusActive         LifeCyclePolicyStatus       = "active"
	LifeCyclePolicyStatusPaused         LifeCyclePolicyStatus       = "paused"
	LifeCyclePolicyActionVolumeSnapshot LifeCyclePolicyAction       = "volume_snapshot"
)

var (
	ErrLifeCyclePolicyInvalidScheduleType = fmt.Errorf("invalid schedule type")
	ErrLifeCyclePolicyInvalidStatus       = fmt.Errorf("invalid lifecycle policy status")
	ErrLifeCyclePolicyInvalidAction       = fmt.Errorf("invalid lifecycle policy action")
)

func (t LifeCyclePolicyScheduleType) List() []LifeCyclePolicyScheduleType {
	return []LifeCyclePolicyScheduleType{LifeCyclePolicyScheduleTypeInterval, LifeCyclePolicyScheduleTypeCron}
}

func (t LifeCyclePolicyScheduleType) String() string {
	return string(t)
}

func (t LifeCyclePolicyScheduleType) StringList() []string {
	lst := t.List()
	strings := make([]string, 0, len(lst))
	for _, x := range lst {
		strings = append(strings, x.String())
	}

	return strings
}

func (t LifeCyclePolicyScheduleType) IsValid() error {
	for _, x := range t.List() {
		if t == x {
			return nil
		}
	}

	return fmt.Errorf("%w: %v", ErrLifeCyclePolicyInvalidScheduleType, t)
}

func (s LifeCyclePolicyStatus) List() []LifeCyclePolicyStatus {
	return []LifeCyclePolicyStatus{LifeCyclePolicyStatusPaused, LifeCyclePolicyStatusActive}
}

func (s LifeCyclePolicyStatus) String() string {
	return string(s)
}

func (s LifeCyclePolicyStatus) StringList() []string {
	lst := s.List()
	strings := make([]string, 0, len(lst))
	for _, x := range lst {
		strings = append(strings, x.String())
	}

	return strings
}

func (s LifeCyclePolicyStatus) IsValid() error {
	for _, x := range s.List() {
		if s == x {
			return nil
		}
	}

	return fmt.Errorf("%w: %v", ErrLifeCyclePolicyInvalidStatus, s)
}

func (a LifeCyclePolicyAction) List() []LifeCyclePolicyAction {
	return []LifeCyclePolicyAction{LifeCyclePolicyActionVolumeSnapshot}
}

func (a LifeCyclePolicyAction) String() string {
	return string(a)
}

func (a LifeCyclePolicyAction) StringList() []string {
	lst := a.List()
	s := make([]string, 0, len(lst))
	for _, x := range lst {
		s = append(s, x.String())
	}

	return s
}

func (a LifeCyclePolicyAction) IsValid() error {
	for _, x := range a.List() {
		if a == x {
			return nil
		}
	}

	return fmt.Errorf("%w: %v", ErrLifeCyclePolicyInvalidAction, a)
}

type LifeCyclePolicyRetentionTimer struct {
	Weeks   int `json:"weeks,omitempty"`
	Days    int `json:"days,omitempty"`
	Hours   int `json:"hours,omitempty"`
	Minutes int `json:"minutes,omitempty"`
}

// LifeCyclePolicySchedule represents a schedule resource.
type LifeCyclePolicySchedule interface {
	GetCommonSchedule() LifeCyclePolicyCommonSchedule
}

type LifeCyclePolicyCommonSchedule struct {
	Type                 LifeCyclePolicyScheduleType    `json:"type"`
	ID                   string                         `json:"id"`
	Owner                string                         `json:"owner"`
	OwnerID              int                            `json:"owner_id"`
	MaxQuantity          int                            `json:"max_quantity"`
	UserID               int                            `json:"user_id"`
	ResourceNameTemplate string                         `json:"resource_name_template"`
	RetentionTime        *LifeCyclePolicyRetentionTimer `json:"retention_time"`
}

type LifeCyclePolicyIntervalSchedule struct {
	LifeCyclePolicyCommonSchedule
	Weeks   int `json:"weeks"`
	Days    int `json:"days"`
	Hours   int `json:"hours"`
	Minutes int `json:"minutes"`
}

type LifeCyclePolicyCronSchedule struct {
	LifeCyclePolicyCommonSchedule
	Timezone  string `json:"timezone"`
	Week      string `json:"week"`
	DayOfWeek string `json:"day_of_week"`
	Month     string `json:"month"`
	Day       string `json:"day"`
	Hour      string `json:"hour"`
	Minute    string `json:"minute"`
}

func (s LifeCyclePolicyCronSchedule) GetCommonSchedule() LifeCyclePolicyCommonSchedule {
	return s.LifeCyclePolicyCommonSchedule
}

func (s LifeCyclePolicyIntervalSchedule) GetCommonSchedule() LifeCyclePolicyCommonSchedule {
	return s.LifeCyclePolicyCommonSchedule
}

// LifeCyclePolicyRawSchedule is internal struct for unmarshalling into LifeCyclePolicySchedule.
type LifeCyclePolicyRawSchedule struct {
	json.RawMessage
}

// Cook is method for unmarshalling LifeCyclePolicyRawSchedule into LifeCyclePolicySchedule.
func (r LifeCyclePolicyRawSchedule) Cook() (LifeCyclePolicySchedule, error) {
	var typeStruct struct {
		LifeCyclePolicyScheduleType `json:"type"`
	}
	//nolint:staticcheck
	if err := json.Unmarshal(r.RawMessage, &typeStruct); err != nil {
		return nil, err
	}
	switch typeStruct.LifeCyclePolicyScheduleType {
	default:
		return nil, fmt.Errorf("%w: %s", ErrLifeCyclePolicyInvalidScheduleType, typeStruct.LifeCyclePolicyScheduleType)
	case LifeCyclePolicyScheduleTypeCron:
		var cronSchedule LifeCyclePolicyCronSchedule
		if err := json.Unmarshal(r.RawMessage, &cronSchedule); err != nil {
			return nil, err
		}
		return cronSchedule, nil
	case LifeCyclePolicyScheduleTypeInterval:
		var intervalSchedule LifeCyclePolicyIntervalSchedule
		if err := json.Unmarshal(r.RawMessage, &intervalSchedule); err != nil {
			return nil, err
		}
		return intervalSchedule, nil
	}
}

// LifeCyclePolicyVolume represents a volume resource.
type LifeCyclePolicyVolume struct {
	ID   string `json:"volume_id"`
	Name string `json:"volume_name"`
}

// LifeCyclePolicy represents a lifecycle policy resource.
type LifeCyclePolicy struct {
	Name      string                    `json:"name"`
	ID        int                       `json:"id"`
	Action    LifeCyclePolicyAction     `json:"action"`
	ProjectID int                       `json:"project_id"`
	Status    LifeCyclePolicyStatus     `json:"status"`
	UserID    int                       `json:"user_id"`
	RegionID  int                       `json:"region_id"`
	Volumes   []LifeCyclePolicyVolume   `json:"volumes"`
	Schedules []LifeCyclePolicySchedule `json:"schedules"`
}

type rawLifeCyclePolicyRoot struct {
	Count        int
	LifePolicies []rawLifeCyclePolicy `json:"results"`
}

// rawLifeCyclePolicy is internal struct for unmarshalling into LifeCyclePolicy.
type rawLifeCyclePolicy struct {
	Name      string                       `json:"name"`
	ID        int                          `json:"id"`
	Action    LifeCyclePolicyAction        `json:"action"`
	ProjectID int                          `json:"project_id"`
	Status    LifeCyclePolicyStatus        `json:"status"`
	UserID    int                          `json:"user_id"`
	RegionID  int                          `json:"region_id"`
	Volumes   []LifeCyclePolicyVolume      `json:"volumes"`
	Schedules []LifeCyclePolicyRawSchedule `json:"schedules"`
}

//
// type MaxPolicyUsage struct {
// 	CountUsage     int             `json:"max_volume_snapshot_count_usage"`
// 	SizeUsage      int             `json:"max_volume_snapshot_size_usage"`
// 	SequenceLength int             `json:"max_volume_snapshot_sequence_length"`
// 	MaxCost        PolicyUsageCost `json:"max_cost"`
// }
//
// type PolicyUsageCost struct {
// 	CurrencyCode  edgecloud.LifeCyclePolicyCurrency `json:"currency_code"`
// 	PricePerHour  decimal.Decimal    `json:"price_per_hour"`
// 	PricePerMonth decimal.Decimal    `json:"price_per_month"`
// 	PriceStatus   string             `json:"price_status"`
// }

// cook is internal method for unmarshalling rawLifeCyclePolicy into LifeCyclePolicy.
func (rawPolicy rawLifeCyclePolicy) cook() (*LifeCyclePolicy, error) {
	schedules := make([]LifeCyclePolicySchedule, len(rawPolicy.Schedules))
	for i, b := range rawPolicy.Schedules {
		s, err := b.Cook()
		if err != nil {
			return nil, err
		}
		schedules[i] = s
	}
	rawPolicy.Schedules = nil
	b, err := json.Marshal(rawPolicy)
	if err != nil {
		return nil, err
	}
	var policy LifeCyclePolicy
	if err := json.Unmarshal(b, &policy); err != nil {
		return nil, err
	}
	policy.Schedules = schedules

	return &policy, nil
}

type LifeCyclePolicyGetOptions struct {
	NeedVolumes bool `url:"need_volumes, omitempty" validate:"omitempty"`
}

type LifeCyclePolicyListOptions LifeCyclePolicyGetOptions

type LifeCyclePolicyCreateRequest struct {
	Name      string                                 `json:"name" validate:"required,name"`
	Status    LifeCyclePolicyStatus                  `json:"status,omitempty" validate:"omitempty,enum"`
	Action    LifeCyclePolicyAction                  `json:"action" validate:"required,enum"`
	Schedules []LifeCyclePolicyCreateScheduleRequest `json:"schedules,omitempty" validate:"dive"`
	VolumeIds []string                               `json:"volume_ids,omitempty"`
}

// LifeCyclePolicyAddSchedulesRequest represents options for AddSchedules.
type LifeCyclePolicyAddSchedulesRequest struct {
	Schedules []LifeCyclePolicyCreateScheduleRequest `json:"schedules" validate:"required,dive"`
}

// LifeCyclePolicyRemoveSchedulesRequest represents options for RemoveSchedules.
type LifeCyclePolicyRemoveSchedulesRequest struct {
	ScheduleIDs []string `json:"schedule_ids" validate:"required"`
}

// LifeCyclePolicyAddVolumesRequest represents options for AddVolumes.
// Volumes already managed by policy are ignored.
type LifeCyclePolicyAddVolumesRequest struct {
	VolumeIds []string `json:"volume_ids" validate:"required"`
}

// LifeCyclePolicyRemoveVolumesRequest represents options for AddVolumes.
// Volumes already managed by policy are ignored.
type LifeCyclePolicyRemoveVolumesRequest struct {
	VolumeIds []string `json:"volume_ids" validate:"required"`
}

type LifeCyclePolicyUpdateRequest struct {
	Name   string                `json:"name" validate:"required,name"`
	Status LifeCyclePolicyStatus `json:"status,omitempty" validate:"omitempty,enum"`
}

type LifeCyclePolicyEstimateOpts struct {
	Name      string                `json:"name" required:"true" validate:"required"`
	VolumeIds []string              `json:"volume_ids"`
	Status    LifeCyclePolicyStatus `json:"status,omitempty" validate:"omitempty,enum"`
	Action    LifeCyclePolicyAction `json:"action" validate:"required,enum"`
}

// LifeCyclePolicyEstimateCronRequest represent options for EstimateCronMaxPolicyUsage.
type LifeCyclePolicyEstimateCronRequest struct {
	LifeCyclePolicyEstimateOpts
	Schedules []LifeCyclePolicyCreateScheduleRequest `json:"schedules"`
}

// LifeCyclePolicyEstimateIntervalRequest represent options for EstimateIntervalMaxPolicyUsage.
type LifeCyclePolicyEstimateIntervalRequest struct {
	LifeCyclePolicyEstimateOpts
	Schedules []LifeCyclePolicyCreateIntervalScheduleRequest `json:"schedules"`
}

// LifeCyclePolicyCreateScheduleRequest represents options used to create a single schedule.
type LifeCyclePolicyCreateScheduleRequest interface {
	SetCommonCreateScheduleOpts(opts LifeCyclePolicyCommonCreateScheduleRequest)
}

type LifeCyclePolicyCurrency struct {
	*currency.Currency
}

func ParseCurrency(s string) (*LifeCyclePolicyCurrency, error) {
	c, err := currency.Get(s)
	if err != nil {
		return nil, err
	}
	return &LifeCyclePolicyCurrency{Currency: c}, nil
}

// String - implements Stringer.
func (c LifeCyclePolicyCurrency) String() string {
	return c.Currency.Code()
}

// UnmarshalJSON - implements Unmarshaler interface for LifeCyclePolicyCurrency.
func (c *LifeCyclePolicyCurrency) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v, err := ParseCurrency(s)
	if err != nil {
		return err
	}
	*c = *v

	return nil
}

// MarshalJSON - implements Marshaler interface for LifeCyclePolicyCurrency.
func (c LifeCyclePolicyCurrency) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

type LifeCyclePolicyCommonCreateScheduleRequest struct {
	Type                 LifeCyclePolicyScheduleType    `json:"type" validate:"required,enum"`
	ResourceNameTemplate string                         `json:"resource_name_template,omitempty"`
	MaxQuantity          int                            `json:"max_quantity" validate:"required,gt=0,lt=10001"`
	RetentionTime        *LifeCyclePolicyRetentionTimer `json:"retention_time,omitempty"`
}

// LifeCyclePolicyCreateCronScheduleRequest represents options used to create a single cron schedule.
type LifeCyclePolicyCreateCronScheduleRequest struct { // TODO: validate?
	LifeCyclePolicyCommonCreateScheduleRequest
	Timezone  string `json:"timezone,omitempty"`
	Week      string `json:"week,omitempty"`
	DayOfWeek string `json:"day_of_week,omitempty"`
	Month     string `json:"month,omitempty"`
	Day       string `json:"day,omitempty"`
	Hour      string `json:"hour,omitempty"`
	Minute    string `json:"minute,omitempty"`
}

// LifeCyclePolicyCreateIntervalScheduleRequest represents options used to create a single interval schedule.
type LifeCyclePolicyCreateIntervalScheduleRequest struct {
	LifeCyclePolicyCommonCreateScheduleRequest
	Weeks   int `json:"weeks,omitempty" validate:"required_without_all=Days Hours Minutes"`
	Days    int `json:"days,omitempty" validate:"required_without_all=Weeks Hours Minutes"`
	Hours   int `json:"hours,omitempty" validate:"required_without_all=Weeks Days Minutes"`
	Minutes int `json:"minutes,omitempty" validate:"required_without_all=Weeks Days Hours"`
}

type LifeCyclePolicyMaxPolicyUsage struct {
	CountUsage     int                            `json:"max_volume_snapshot_count_usage"`
	SizeUsage      int                            `json:"max_volume_snapshot_size_usage"`
	SequenceLength int                            `json:"max_volume_snapshot_sequence_length"`
	MaxCost        LifeCyclePolicyPolicyUsageCost `json:"max_cost"`
}

type LifeCyclePolicyPolicyUsageCost struct {
	CurrencyCode  LifeCyclePolicyCurrency `json:"currency_code"`
	PricePerHour  decimal.Decimal         `json:"price_per_hour"`
	PricePerMonth decimal.Decimal         `json:"price_per_month"`
	PriceStatus   string                  `json:"price_status"`
}

func (opts *LifeCyclePolicyCreateCronScheduleRequest) SetCommonCreateScheduleOpts(common LifeCyclePolicyCommonCreateScheduleRequest) {
	opts.LifeCyclePolicyCommonCreateScheduleRequest = common
}

func (opts *LifeCyclePolicyCreateIntervalScheduleRequest) SetCommonCreateScheduleOpts(common LifeCyclePolicyCommonCreateScheduleRequest) {
	opts.LifeCyclePolicyCommonCreateScheduleRequest = common
}

// LifeCyclePoliciesService is an interface for creating and managing lifecycle policies with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/Lifecycle-policy
type LifeCyclePoliciesService interface {
	List(context.Context, *LifeCyclePolicyListOptions) ([]LifeCyclePolicy, *Response, error)
	Get(context.Context, int, *LifeCyclePolicyGetOptions) (*LifeCyclePolicy, *Response, error)
	Create(context.Context, *LifeCyclePolicyCreateRequest) (*LifeCyclePolicy, *Response, error)
	Delete(context.Context, int) (*Response, error)
	Update(context.Context, int, *LifeCyclePolicyUpdateRequest) (*LifeCyclePolicy, *Response, error)
	AddSchedules(context.Context, int, *LifeCyclePolicyAddSchedulesRequest) (*LifeCyclePolicy, *Response, error)
	RemoveSchedules(context.Context, int, *LifeCyclePolicyRemoveSchedulesRequest) (*LifeCyclePolicy, *Response, error)
	AddVolumes(context.Context, int, *LifeCyclePolicyAddVolumesRequest) (*LifeCyclePolicy, *Response, error)
	RemoveVolumes(context.Context, int, *LifeCyclePolicyRemoveVolumesRequest) (*LifeCyclePolicy, *Response, error)
	EstimateCronMaxPolicyUsage(context.Context, *LifeCyclePolicyEstimateCronRequest) (*LifeCyclePolicyMaxPolicyUsage, *Response, error)
	EstimateIntervalMaxPolicyUsage(context.Context, *LifeCyclePolicyEstimateIntervalRequest) (*LifeCyclePolicyMaxPolicyUsage, *Response, error)

	// Share(context.Context, string, *KeyPairShareRequest) (*KeyPair, *Response, error)
}

// LifeCyclePoliciesServiceOp handles communication with lifecycle policies methods of the EdgecenterCloud API.
type LifeCyclePoliciesServiceOp struct {
	client *Client
}

// List returns a list of lifecycle policies.
func (s LifeCyclePoliciesServiceOp) List(ctx context.Context, listOpts *LifeCyclePolicyListOptions) ([]LifeCyclePolicy, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(lifecyclePoliciesBasePathV1)
	path, err := addOptions(path, listOpts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	rawRoot := new(rawLifeCyclePolicyRoot)
	resp, err := s.client.Do(ctx, req, rawRoot)
	if err != nil {
		return nil, resp, err
	}
	lifeCeyclePolicies := make([]LifeCyclePolicy, 0, len(rawRoot.LifePolicies))
	for _, rawPolicy := range rawRoot.LifePolicies {
		lifeCeyclePolicy, err := rawPolicy.cook()
		if err != nil {
			return nil, resp, err
		}
		lifeCeyclePolicies = append(lifeCeyclePolicies, *lifeCeyclePolicy)
	}

	return lifeCeyclePolicies, resp, err
}

// Get returns a lifecycle policy with specified unique id.
func (s LifeCyclePoliciesServiceOp) Get(ctx context.Context, lifecyclePolicyID int, getOpts *LifeCyclePolicyGetOptions) (*LifeCyclePolicy, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%d", s.client.addProjectRegionPath(lifecyclePoliciesBasePathV1), lifecyclePolicyID)
	path, err := addOptions(path, getOpts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	rawLifeCyclePolicy := new(rawLifeCyclePolicy)

	resp, err := s.client.Do(ctx, req, rawLifeCyclePolicy)
	if err != nil {
		return nil, resp, err
	}
	lifeCyclePolicy, err := rawLifeCyclePolicy.cook()
	if err != nil {
		return nil, nil, err
	}

	return lifeCyclePolicy, resp, err
}

// Create is create new lifecycle policy.
func (s LifeCyclePoliciesServiceOp) Create(ctx context.Context, reqBody *LifeCyclePolicyCreateRequest) (*LifeCyclePolicy, *Response, error) {
	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(lifecyclePoliciesBasePathV1)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	rawLifeCyclePolicy := new(rawLifeCyclePolicy)

	resp, err := s.client.Do(ctx, req, rawLifeCyclePolicy)
	if err != nil {
		return nil, resp, err
	}
	lifeCyclePolicy, err := rawLifeCyclePolicy.cook()
	if err != nil {
		return nil, nil, err
	}

	return lifeCyclePolicy, resp, err
}

// Update updates a lifecycle policy with specified unique id.
// reqBody are used to construct request body.
func (s LifeCyclePoliciesServiceOp) Update(ctx context.Context, lifeCyclePolicyID int, reqBody *LifeCyclePolicyUpdateRequest) (*LifeCyclePolicy, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}
	path := s.client.addProjectRegionPath(lifecyclePoliciesBasePathV1)
	path = fmt.Sprintf("%s/%d", path, lifeCyclePolicyID)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	rawLifeCyclePolicy := new(rawLifeCyclePolicy)

	resp, err := s.client.Do(ctx, req, rawLifeCyclePolicy)
	if err != nil {
		return nil, resp, err
	}
	lifeCyclePolicy, err := rawLifeCyclePolicy.cook()
	if err != nil {
		return nil, nil, err
	}

	return lifeCyclePolicy, resp, err
}

// Delete deletes a lifecycle policy with specified unique id.
func (s LifeCyclePoliciesServiceOp) Delete(ctx context.Context, lifeCyclePolicyID int) (*Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	path := fmt.Sprintf("%s/%d", s.client.addProjectRegionPath(lifecyclePoliciesBasePathV1), lifeCyclePolicyID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// AddSchedules adds a schedules to lifecycle policy with specified unique id.
func (s LifeCyclePoliciesServiceOp) AddSchedules(ctx context.Context, lifeCyclePolicyID int, reqBody *LifeCyclePolicyAddSchedulesRequest) (*LifeCyclePolicy, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}
	path := s.client.addProjectRegionPath(lifecyclePoliciesBasePathV1)
	path = fmt.Sprintf("%s/%d/%s", path, lifeCyclePolicyID, addSchedulesSubPath)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	rawLifeCyclePolicy := new(rawLifeCyclePolicy)

	resp, err := s.client.Do(ctx, req, rawLifeCyclePolicy)
	if err != nil {
		return nil, resp, err
	}
	lifeCyclePolicy, err := rawLifeCyclePolicy.cook()
	if err != nil {
		return nil, nil, err
	}

	return lifeCyclePolicy, resp, err
}

// RemoveSchedules removes a schedules from lifecycle policy with specified unique id.
func (s LifeCyclePoliciesServiceOp) RemoveSchedules(ctx context.Context, lifeCyclePolicyID int, reqBody *LifeCyclePolicyRemoveSchedulesRequest) (*LifeCyclePolicy, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}
	path := s.client.addProjectRegionPath(lifecyclePoliciesBasePathV1)
	path = fmt.Sprintf("%s/%d/%s", path, lifeCyclePolicyID, removeSchedulesSubPath)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	rawLifeCyclePolicy := new(rawLifeCyclePolicy)

	resp, err := s.client.Do(ctx, req, rawLifeCyclePolicy)
	if err != nil {
		return nil, resp, err
	}
	lifeCyclePolicy, err := rawLifeCyclePolicy.cook()
	if err != nil {
		return nil, nil, err
	}

	return lifeCyclePolicy, resp, err
}

// AddVolumes adds a volumes to lifecycle policy with specified unique id.
func (s LifeCyclePoliciesServiceOp) AddVolumes(ctx context.Context, lifeCyclePolicyID int, reqBody *LifeCyclePolicyAddVolumesRequest) (*LifeCyclePolicy, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}
	path := s.client.addProjectRegionPath(lifecyclePoliciesBasePathV1)
	path = fmt.Sprintf("%s/%d/%s", path, lifeCyclePolicyID, addVolumesSubPath)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	rawLifeCyclePolicy := new(rawLifeCyclePolicy)

	resp, err := s.client.Do(ctx, req, rawLifeCyclePolicy)
	if err != nil {
		return nil, resp, err
	}
	lifeCyclePolicy, err := rawLifeCyclePolicy.cook()
	if err != nil {
		return nil, nil, err
	}

	return lifeCyclePolicy, resp, err
}

// RemoveVolumes removes a volumes from lifecycle policy with specified unique id.
func (s LifeCyclePoliciesServiceOp) RemoveVolumes(ctx context.Context, lifeCyclePolicyID int, reqBody *LifeCyclePolicyRemoveVolumesRequest) (*LifeCyclePolicy, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}
	path := s.client.addProjectRegionPath(lifecyclePoliciesBasePathV1)
	path = fmt.Sprintf("%s/%d/%s", path, lifeCyclePolicyID, removeVolumesSubPath)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	rawLifeCyclePolicy := new(rawLifeCyclePolicy)

	resp, err := s.client.Do(ctx, req, rawLifeCyclePolicy)
	if err != nil {
		return nil, resp, err
	}
	lifeCyclePolicy, err := rawLifeCyclePolicy.cook()
	if err != nil {
		return nil, nil, err
	}

	return lifeCyclePolicy, resp, err
}

// EstimateCronMaxPolicyUsage estimates usage of resources and costs for CRON lifecycle policy.
func (s LifeCyclePoliciesServiceOp) EstimateCronMaxPolicyUsage(ctx context.Context, reqBody *LifeCyclePolicyEstimateCronRequest) (*LifeCyclePolicyMaxPolicyUsage, *Response, error) {
	return s.estimateMaxPolicyUsage(ctx, reqBody)
}

// EstimateIntervalMaxPolicyUsage estimates usage of resources and costs for Interval lifecycle policy.
func (s LifeCyclePoliciesServiceOp) EstimateIntervalMaxPolicyUsage(ctx context.Context, reqBody *LifeCyclePolicyEstimateIntervalRequest) (*LifeCyclePolicyMaxPolicyUsage, *Response, error) {
	return s.estimateMaxPolicyUsage(ctx, reqBody)
}

func (s LifeCyclePoliciesServiceOp) estimateMaxPolicyUsage(ctx context.Context, reqBody interface{}) (*LifeCyclePolicyMaxPolicyUsage, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	if reqBody == nil {
		return nil, nil, NewArgError("reqBody", "cannot be nil")
	}
	path := s.client.addProjectRegionPath(lifecyclePoliciesBasePathV1)
	path = fmt.Sprintf("%s/%s", path, estimateMaxPolicyUsageSubPath)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	lifePolicyMaxUsage := new(LifeCyclePolicyMaxPolicyUsage)

	resp, err := s.client.Do(ctx, req, lifePolicyMaxUsage)
	if err != nil {
		return nil, resp, err
	}

	return lifePolicyMaxUsage, resp, err
}
