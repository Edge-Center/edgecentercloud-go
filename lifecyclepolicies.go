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
	ScheduleType string
	PolicyStatus string
	PolicyAction string
)

const (
	lifecyclePoliciesBasePathV1                = "/v1/lifecycle_policies/"
	addSchedulesSubPath                        = "add_schedules"
	addVolumesSubPath                          = "add_volumes_to_policy"
	removeSchedulesSubPath                     = "remove_schedules"
	removeVolumesSubPath                       = "remove_volumes_from_policy"
	estimateMaxPolicyUsageSubPath              = "estimate_max_policy_usage"
	ScheduleTypeCron              ScheduleType = "cron"
	ScheduleTypeInterval          ScheduleType = "interval"
	PolicyStatusActive            PolicyStatus = "active"
	PolicyStatusPaused            PolicyStatus = "paused"
	PolicyActionVolumeSnapshot    PolicyAction = "volume_snapshot"
)

var (
	ErrLifeCyclePolicyInvalidScheduleType = fmt.Errorf("invalid schedule type")
	ErrLifeCyclePolicyInvalidStatus       = fmt.Errorf("invalid lifecycle policy status")
	ErrLifeCyclePolicyInvalidAction       = fmt.Errorf("invalid lifecycle policy action")
)

func (t ScheduleType) List() []ScheduleType {
	return []ScheduleType{ScheduleTypeInterval, ScheduleTypeCron}
}

func (t ScheduleType) String() string {
	return string(t)
}

func (t ScheduleType) StringList() []string {
	lst := t.List()
	strings := make([]string, 0, len(lst))
	for _, x := range lst {
		strings = append(strings, x.String())
	}

	return strings
}

func (t ScheduleType) IsValid() error {
	for _, x := range t.List() {
		if t == x {
			return nil
		}
	}

	return fmt.Errorf("%w: %v", ErrLifeCyclePolicyInvalidScheduleType, t)
}

func (s PolicyStatus) List() []PolicyStatus {
	return []PolicyStatus{PolicyStatusPaused, PolicyStatusActive}
}

func (s PolicyStatus) String() string {
	return string(s)
}

func (s PolicyStatus) StringList() []string {
	lst := s.List()
	strings := make([]string, 0, len(lst))
	for _, x := range lst {
		strings = append(strings, x.String())
	}

	return strings
}

func (s PolicyStatus) IsValid() error {
	for _, x := range s.List() {
		if s == x {
			return nil
		}
	}

	return fmt.Errorf("%w: %v", ErrLifeCyclePolicyInvalidStatus, s)
}

func (a PolicyAction) List() []PolicyAction {
	return []PolicyAction{PolicyActionVolumeSnapshot}
}

func (a PolicyAction) String() string {
	return string(a)
}

func (a PolicyAction) StringList() []string {
	lst := a.List()
	s := make([]string, 0, len(lst))
	for _, x := range lst {
		s = append(s, x.String())
	}

	return s
}

func (a PolicyAction) IsValid() error {
	for _, x := range a.List() {
		if a == x {
			return nil
		}
	}

	return fmt.Errorf("%w: %v", ErrLifeCyclePolicyInvalidAction, a)
}

type RetentionTimer struct {
	Weeks   int `json:"weeks,omitempty"`
	Days    int `json:"days,omitempty"`
	Hours   int `json:"hours,omitempty"`
	Minutes int `json:"minutes,omitempty"`
}

// Schedule represents a schedule resource.
type Schedule interface {
	GetCommonSchedule() CommonSchedule
}

type CommonSchedule struct {
	Type                 ScheduleType    `json:"type"`
	ID                   string          `json:"id"`
	Owner                string          `json:"owner"`
	OwnerID              int             `json:"owner_id"`
	MaxQuantity          int             `json:"max_quantity"`
	UserID               int             `json:"user_id"`
	ResourceNameTemplate string          `json:"resource_name_template"`
	RetentionTime        *RetentionTimer `json:"retention_time"`
}

type IntervalSchedule struct {
	CommonSchedule
	Weeks   int `json:"weeks"`
	Days    int `json:"days"`
	Hours   int `json:"hours"`
	Minutes int `json:"minutes"`
}

type CronSchedule struct {
	CommonSchedule
	Timezone  string `json:"timezone"`
	Week      string `json:"week"`
	DayOfWeek string `json:"day_of_week"`
	Month     string `json:"month"`
	Day       string `json:"day"`
	Hour      string `json:"hour"`
	Minute    string `json:"minute"`
}

func (s CronSchedule) GetCommonSchedule() CommonSchedule {
	return s.CommonSchedule
}

func (s IntervalSchedule) GetCommonSchedule() CommonSchedule {
	return s.CommonSchedule
}

// RawSchedule is internal struct for unmarshalling into Schedule.
type RawSchedule struct {
	json.RawMessage
}

// Cook is method for unmarshalling RawSchedule into Schedule.
func (r RawSchedule) Cook() (Schedule, error) {
	var typeStruct struct {
		ScheduleType `json:"type"`
	}
	//nolint:staticcheck
	if err := json.Unmarshal(r.RawMessage, &typeStruct); err != nil {
		return nil, err
	}
	switch typeStruct.ScheduleType {
	default:
		return nil, fmt.Errorf("%w: %s", ErrLifeCyclePolicyInvalidScheduleType, typeStruct.ScheduleType)
	case ScheduleTypeCron:
		var cronSchedule CronSchedule
		if err := json.Unmarshal(r.RawMessage, &cronSchedule); err != nil {
			return nil, err
		}
		return cronSchedule, nil
	case ScheduleTypeInterval:
		var intervalSchedule IntervalSchedule
		if err := json.Unmarshal(r.RawMessage, &intervalSchedule); err != nil {
			return nil, err
		}
		return intervalSchedule, nil
	}
}

// LifecyclePolicyVolume represents a volume resource.
type LifecyclePolicyVolume struct {
	ID   string `json:"volume_id"`
	Name string `json:"volume_name"`
}

// LifecyclePolicy represents a lifecycle policy resource.
type LifecyclePolicy struct {
	Name      string                  `json:"name"`
	ID        int                     `json:"id"`
	Action    PolicyAction            `json:"action"`
	ProjectID int                     `json:"project_id"`
	Status    PolicyStatus            `json:"status"`
	UserID    int                     `json:"user_id"`
	RegionID  int                     `json:"region_id"`
	Volumes   []LifecyclePolicyVolume `json:"volumes"`
	Schedules []Schedule              `json:"schedules"`
}

type rawLifeCyclePolicyRoot struct {
	Count        int
	LifePolicies []rawLifecyclePolicy `json:"results"`
}

// rawLifecyclePolicy is internal struct for unmarshalling into LifecyclePolicy.
type rawLifecyclePolicy struct {
	Name      string                  `json:"name"`
	ID        int                     `json:"id"`
	Action    PolicyAction            `json:"action"`
	ProjectID int                     `json:"project_id"`
	Status    PolicyStatus            `json:"status"`
	UserID    int                     `json:"user_id"`
	RegionID  int                     `json:"region_id"`
	Volumes   []LifecyclePolicyVolume `json:"volumes"`
	Schedules []RawSchedule           `json:"schedules"`
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
// 	CurrencyCode  edgecloud.Currency `json:"currency_code"`
// 	PricePerHour  decimal.Decimal    `json:"price_per_hour"`
// 	PricePerMonth decimal.Decimal    `json:"price_per_month"`
// 	PriceStatus   string             `json:"price_status"`
// }

// cook is internal method for unmarshalling rawLifecyclePolicy into LifecyclePolicy.
func (rawPolicy rawLifecyclePolicy) cook() (*LifecyclePolicy, error) {
	schedules := make([]Schedule, len(rawPolicy.Schedules))
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
	var policy LifecyclePolicy
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
	Status    PolicyStatus                           `json:"status,omitempty" validate:"omitempty,enum"`
	Action    PolicyAction                           `json:"action" validate:"required,enum"`
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
	Name   string       `json:"name" validate:"required,name"`
	Status PolicyStatus `json:"status,omitempty" validate:"omitempty,enum"`
}

type LifeCyclePolicyEstimateOpts struct {
	Name      string       `json:"name" required:"true" validate:"required"`
	VolumeIds []string     `json:"volume_ids"`
	Status    PolicyStatus `json:"status,omitempty" validate:"omitempty,enum"`
	Action    PolicyAction `json:"action" validate:"required,enum"`
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
	SetCommonCreateScheduleOpts(opts CommonCreateScheduleOpts)
}

type Currency struct {
	*currency.Currency
}

func ParseCurrency(s string) (*Currency, error) {
	c, err := currency.Get(s)
	if err != nil {
		return nil, err
	}
	return &Currency{Currency: c}, nil
}

// String - implements Stringer.
func (c Currency) String() string {
	return c.Currency.Code()
}

// UnmarshalJSON - implements Unmarshaler interface for Currency.
func (c *Currency) UnmarshalJSON(data []byte) error {
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

// MarshalJSON - implements Marshaler interface for Currency.
func (c Currency) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

type CommonCreateScheduleOpts struct {
	Type                 ScheduleType    `json:"type" validate:"required,enum"`
	ResourceNameTemplate string          `json:"resource_name_template,omitempty"`
	MaxQuantity          int             `json:"max_quantity" validate:"required,gt=0,lt=10001"`
	RetentionTime        *RetentionTimer `json:"retention_time,omitempty"`
}

// LifeCyclePolicyCreateCronScheduleOpts represents options used to create a single cron schedule.
type LifeCyclePolicyCreateCronScheduleOpts struct { // TODO: validate?
	CommonCreateScheduleOpts
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
	CommonCreateScheduleOpts
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
	CurrencyCode  Currency        `json:"currency_code"`
	PricePerHour  decimal.Decimal `json:"price_per_hour"`
	PricePerMonth decimal.Decimal `json:"price_per_month"`
	PriceStatus   string          `json:"price_status"`
}

func (opts *LifeCyclePolicyCreateCronScheduleOpts) SetCommonCreateScheduleOpts(common CommonCreateScheduleOpts) {
	opts.CommonCreateScheduleOpts = common
}

func (opts *LifeCyclePolicyCreateIntervalScheduleRequest) SetCommonCreateScheduleOpts(common CommonCreateScheduleOpts) {
	opts.CommonCreateScheduleOpts = common
}

// LifeCyclePoliciesService is an interface for creating and managing lifecycle policies with the EdgecenterCloud API.
// See: https://apidocs.edgecenter.ru/cloud#tag/Lifecycle-policy
type LifeCyclePoliciesService interface {
	List(context.Context, *LifeCyclePolicyListOptions) ([]LifecyclePolicy, *Response, error)
	Get(context.Context, int, *LifeCyclePolicyGetOptions) (*LifecyclePolicy, *Response, error)
	Create(context.Context, *LifeCyclePolicyCreateRequest) (*LifecyclePolicy, *Response, error)
	Delete(context.Context, int) (*Response, error)
	Update(context.Context, int, *LifeCyclePolicyUpdateRequest) (*LifecyclePolicy, *Response, error)
	AddSchedules(context.Context, int, *LifeCyclePolicyAddSchedulesRequest) (*LifecyclePolicy, *Response, error)
	RemoveSchedules(context.Context, int, *LifeCyclePolicyRemoveSchedulesRequest) (*LifecyclePolicy, *Response, error)
	AddVolumes(context.Context, int, *LifeCyclePolicyAddVolumesRequest) (*LifecyclePolicy, *Response, error)
	RemoveVolumes(context.Context, int, *LifeCyclePolicyRemoveVolumesRequest) (*LifecyclePolicy, *Response, error)
	EstimateCronMaxPolicyUsage(context.Context, *LifeCyclePolicyEstimateCronRequest) (*LifeCyclePolicyMaxPolicyUsage, *Response, error)
	EstimateIntervalMaxPolicyUsage(context.Context, *LifeCyclePolicyEstimateIntervalRequest) (*LifeCyclePolicyMaxPolicyUsage, *Response, error)

	// Share(context.Context, string, *KeyPairShareRequest) (*KeyPair, *Response, error)
}

// LifeCyclePoliciesServiceOp handles communication with lifecycle policies methods of the EdgecenterCloud API.
type LifeCyclePoliciesServiceOp struct {
	client *Client
}

// List returns a list of lifecycle policies.
func (s LifeCyclePoliciesServiceOp) List(ctx context.Context, listOpts *LifeCyclePolicyListOptions) ([]LifecyclePolicy, *Response, error) {
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
	lifeCeyclePolicies := make([]LifecyclePolicy, 0, len(rawRoot.LifePolicies))
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
func (s LifeCyclePoliciesServiceOp) Get(ctx context.Context, lifecyclePolicyID int, getOpts *LifeCyclePolicyGetOptions) (*LifecyclePolicy, *Response, error) {
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

	rawLifeCyclePolicy := new(rawLifecyclePolicy)

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
func (s LifeCyclePoliciesServiceOp) Create(ctx context.Context, reqBody *LifeCyclePolicyCreateRequest) (*LifecyclePolicy, *Response, error) {
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

	rawLifeCyclePolicy := new(rawLifecyclePolicy)

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
func (s LifeCyclePoliciesServiceOp) Update(ctx context.Context, lifeCyclePolicyID int, reqBody *LifeCyclePolicyUpdateRequest) (*LifecyclePolicy, *Response, error) {
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

	rawLifeCyclePolicy := new(rawLifecyclePolicy)

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
func (s LifeCyclePoliciesServiceOp) AddSchedules(ctx context.Context, lifeCyclePolicyID int, reqBody *LifeCyclePolicyAddSchedulesRequest) (*LifecyclePolicy, *Response, error) {
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

	rawLifeCyclePolicy := new(rawLifecyclePolicy)

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
func (s LifeCyclePoliciesServiceOp) RemoveSchedules(ctx context.Context, lifeCyclePolicyID int, reqBody *LifeCyclePolicyRemoveSchedulesRequest) (*LifecyclePolicy, *Response, error) {
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

	rawLifeCyclePolicy := new(rawLifecyclePolicy)

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
func (s LifeCyclePoliciesServiceOp) AddVolumes(ctx context.Context, lifeCyclePolicyID int, reqBody *LifeCyclePolicyAddVolumesRequest) (*LifecyclePolicy, *Response, error) {
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

	rawLifeCyclePolicy := new(rawLifecyclePolicy)

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
func (s LifeCyclePoliciesServiceOp) RemoveVolumes(ctx context.Context, lifeCyclePolicyID int, reqBody *LifeCyclePolicyRemoveVolumesRequest) (*LifecyclePolicy, *Response, error) {
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

	rawLifeCyclePolicy := new(rawLifecyclePolicy)

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
