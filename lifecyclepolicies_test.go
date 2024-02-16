package edgecloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"testing"

	"github.com/ladydascalie/currency"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

const (
	lifePolicyName       = "new-lp"
	actionVolumeSnapshot = "volume_snapshot"
	statusPaused         = "paused"
)

func TestLifeCyclePoliciesServiceOp_List(t *testing.T) {
	setup()
	defer teardown()

	expectedSchedules := make([]LifeCyclePolicySchedule, 0, 1)
	expectedResp := []LifeCyclePolicy{{ID: testResourceIntID, Schedules: expectedSchedules}}
	URL := path.Join(lifecyclePoliciesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})
	respActual, resp, err := client.LifeCyclePolicies.List(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLifeCyclePoliciesServiceOp_Get(t *testing.T) {
	setup()
	defer teardown()

	expectedSchedules := make([]LifeCyclePolicySchedule, 0, 1)
	expectedResp := LifeCyclePolicy{ID: testResourceIntID, Schedules: expectedSchedules}
	rawResp := rawLifeCyclePolicy{ID: testResourceIntID}
	URL := path.Join(lifecyclePoliciesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), strconv.Itoa(testResourceIntID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(rawResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `%s`, string(resp))
	})
	respActual, resp, err := client.LifeCyclePolicies.Get(ctx, testResourceIntID, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, *respActual, expectedResp)
}

func TestLifeCyclePoliciesServiceOp_Create(t *testing.T) {
	setup()
	defer teardown()

	expectedSchedules := make([]LifeCyclePolicySchedule, 0, 1)
	expectedResp := LifeCyclePolicy{ID: testResourceIntID, Name: lifePolicyName, Schedules: expectedSchedules, Action: actionVolumeSnapshot}
	rawResp := rawLifeCyclePolicy{ID: testResourceIntID, Name: lifePolicyName, Action: actionVolumeSnapshot}
	URL := path.Join(lifecyclePoliciesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, err := json.Marshal(rawResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `%s`, string(resp))
	})
	createRequest := LifeCyclePolicyCreateRequest{
		Name:      "new-lp",
		Status:    "",
		Action:    "volume_snapshot",
		Schedules: nil,
		VolumeIds: nil,
	}
	respActual, resp, err := client.LifeCyclePolicies.Create(ctx, &createRequest)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, *respActual, expectedResp)
}

func TestLifeCyclePoliciesServiceOp_CreateWithSched(t *testing.T) {
	setup()
	defer teardown()

	expectedSchedules := make([]LifeCyclePolicySchedule, 0, 1)
	sched := LifeCyclePolicyIntervalSchedule{
		LifeCyclePolicyCommonSchedule: LifeCyclePolicyCommonSchedule{OwnerID: 928, Type: LifeCyclePolicyScheduleTypeInterval},
	}
	expectedSchedules = append(expectedSchedules, sched)
	expectedRawSchedules := make([]LifeCyclePolicyRawSchedule, 0, 1)
	// rawSched := LifeCyclePolicyRawSchedule{json.RawMessage(`{"owner_id": 928, "user_id": 342026, "resource_name_template": "reserve snap of the volume {volume_id}", "max_quantity": 5, "retention_time": {"weeks": 1, "hours": 1, "days": 1, "minutes": 1}, "owner": "lifecycle_policy", "id": "2a4dbd33-9d89-466d-9ce0-08beaf9d70ea", "type": "interval", "hours": 1, "minutes": 1, "weeks": 1, "days": 1}`)}
	rawSched := LifeCyclePolicyRawSchedule{json.RawMessage(`{"owner_id": 928, "type": "interval"}`)}
	expectedRawSchedules = append(expectedRawSchedules, rawSched)
	expectedResp := LifeCyclePolicy{ID: testResourceIntID, Name: lifePolicyName, Schedules: expectedSchedules, Action: actionVolumeSnapshot}
	rawResp := rawLifeCyclePolicy{ID: testResourceIntID, Name: lifePolicyName, Action: actionVolumeSnapshot, Schedules: expectedRawSchedules}
	URL := path.Join(lifecyclePoliciesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, err := json.Marshal(rawResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `%s`, string(resp))
	})
	createRequest := LifeCyclePolicyCreateRequest{
		Name:      "new-lp",
		Status:    "",
		Action:    "volume_snapshot",
		Schedules: nil,
		VolumeIds: nil,
	}
	respActual, resp, err := client.LifeCyclePolicies.Create(ctx, &createRequest)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, expectedResp, *respActual)
}

func TestLifeCyclePoliciesServiceOp_Update(t *testing.T) {
	setup()
	defer teardown()

	expectedSchedules := make([]LifeCyclePolicySchedule, 0, 1)
	expectedResp := LifeCyclePolicy{ID: testResourceIntID, Name: lifePolicyName, Schedules: expectedSchedules, Action: actionVolumeSnapshot, Status: statusPaused}
	rawResp := rawLifeCyclePolicy{ID: testResourceIntID, Name: lifePolicyName, Action: actionVolumeSnapshot, Status: statusPaused}
	URL := path.Join(lifecyclePoliciesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), strconv.Itoa(testResourceIntID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		resp, err := json.Marshal(rawResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `%s`, string(resp))
	})

	updateRequest := LifeCyclePolicyUpdateRequest{
		Name:   lifePolicyName,
		Status: statusPaused,
	}
	respActual, resp, err := client.LifeCyclePolicies.Update(ctx, testResourceIntID, &updateRequest)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, *respActual, expectedResp)
}

func TestLifeCyclePoliciesServiceOp_Delete(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(lifecyclePoliciesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), strconv.Itoa(testResourceIntID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
		resp, err := json.Marshal(nil)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `%s`, string(resp))
	})

	resp, err := client.LifeCyclePolicies.Delete(ctx, testResourceIntID)
	require.NoError(t, err)
	require.Equal(t, 204, resp.StatusCode)
}

func TestLifeCyclePoliciesServiceOp_AddSchedules(t *testing.T) {
	setup()
	defer teardown()

	expectedSchedules := make([]LifeCyclePolicySchedule, 0, 1)
	expectedResp := LifeCyclePolicy{ID: testResourceIntID, Name: lifePolicyName, Schedules: expectedSchedules, Action: actionVolumeSnapshot, Status: statusPaused}
	rawResp := rawLifeCyclePolicy{ID: testResourceIntID, Name: lifePolicyName, Action: actionVolumeSnapshot, Status: statusPaused}
	URL := path.Join(lifecyclePoliciesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), strconv.Itoa(testResourceIntID), addSchedulesSubPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, err := json.Marshal(rawResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `%s`, string(resp))
	})

	scheduleIntervalReq := LifeCyclePolicyCreateIntervalScheduleRequest{}
	addSchedules := make([]LifeCyclePolicyCreateScheduleRequest, 0, 1)
	addSchedules = append(addSchedules, &scheduleIntervalReq)

	addSchedulesRequest := LifeCyclePolicyAddSchedulesRequest{Schedules: addSchedules}

	respActual, resp, err := client.LifeCyclePolicies.AddSchedules(ctx, testResourceIntID, &addSchedulesRequest)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, *respActual, expectedResp)
}

func TestLifeCyclePoliciesServiceOp_RemoveSchedules(t *testing.T) {
	setup()
	defer teardown()

	expectedSchedules := make([]LifeCyclePolicySchedule, 0, 1)
	expectedResp := LifeCyclePolicy{ID: testResourceIntID, Name: lifePolicyName, Schedules: expectedSchedules, Action: actionVolumeSnapshot, Status: statusPaused}
	rawResp := rawLifeCyclePolicy{ID: testResourceIntID, Name: lifePolicyName, Action: actionVolumeSnapshot, Status: statusPaused}
	URL := path.Join(lifecyclePoliciesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), strconv.Itoa(testResourceIntID), removeSchedulesSubPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, err := json.Marshal(rawResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `%s`, string(resp))
	})

	removeSchedulesRequest := LifeCyclePolicyRemoveSchedulesRequest{ScheduleIDs: []string{""}}

	respActual, resp, err := client.LifeCyclePolicies.RemoveSchedules(ctx, testResourceIntID, &removeSchedulesRequest)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, *respActual, expectedResp)
}

func TestLifeCyclePoliciesServiceOp_AddVolumes(t *testing.T) {
	setup()
	defer teardown()

	expectedSchedules := make([]LifeCyclePolicySchedule, 0, 1)
	expectedResp := LifeCyclePolicy{ID: testResourceIntID, Name: lifePolicyName, Schedules: expectedSchedules, Action: actionVolumeSnapshot, Status: statusPaused}
	rawResp := rawLifeCyclePolicy{ID: testResourceIntID, Name: lifePolicyName, Action: actionVolumeSnapshot, Status: statusPaused}
	URL := path.Join(lifecyclePoliciesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), strconv.Itoa(testResourceIntID), addVolumesSubPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		resp, err := json.Marshal(rawResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `%s`, string(resp))
	})

	addVolumesRequest := LifeCyclePolicyAddVolumesRequest{VolumeIds: []string{""}}

	respActual, resp, err := client.LifeCyclePolicies.AddVolumes(ctx, testResourceIntID, &addVolumesRequest)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, *respActual, expectedResp)
}

func TestLifeCyclePoliciesServiceOp_RemoveVolumes(t *testing.T) {
	setup()
	defer teardown()

	expectedSchedules := make([]LifeCyclePolicySchedule, 0, 1)
	expectedResp := LifeCyclePolicy{ID: testResourceIntID, Name: lifePolicyName, Schedules: expectedSchedules, Action: actionVolumeSnapshot, Status: statusPaused}
	rawResp := rawLifeCyclePolicy{ID: testResourceIntID, Name: lifePolicyName, Action: actionVolumeSnapshot, Status: statusPaused}
	URL := path.Join(lifecyclePoliciesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), strconv.Itoa(testResourceIntID), removeVolumesSubPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		resp, err := json.Marshal(rawResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `%s`, string(resp))
	})

	removeVolumesRequest := LifeCyclePolicyRemoveVolumesRequest{VolumeIds: []string{""}}

	respActual, resp, err := client.LifeCyclePolicies.RemoveVolumes(ctx, testResourceIntID, &removeVolumesRequest)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, *respActual, expectedResp)
}

func TestLifeCyclePoliciesServiceOp_EstimateCronMaxPolicyUsage(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(lifecyclePoliciesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), estimateMaxPolicyUsageSubPath)
	expResp := LifeCyclePolicyMaxPolicyUsage{
		CountUsage:     2,
		SizeUsage:      2,
		SequenceLength: 2,
		MaxCost: LifeCyclePolicyPolicyUsageCost{
			CurrencyCode:  LifeCyclePolicyCurrency{Currency: &currency.RUB},
			PricePerHour:  decimal.New(0, 0),
			PricePerMonth: decimal.New(0, 0),
		},
	}

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, err := json.Marshal(expResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `%s`, string(resp))
	})

	ectimateCronRequest := LifeCyclePolicyEstimateCronRequest{}

	respActual, resp, err := client.LifeCyclePolicies.EstimateCronMaxPolicyUsage(ctx, &ectimateCronRequest)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, expResp, *respActual)
}

func TestLifeCyclePoliciesServiceOp_EstimateIntervalMaxPolicyUsage(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(lifecyclePoliciesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), estimateMaxPolicyUsageSubPath)
	expResp := LifeCyclePolicyMaxPolicyUsage{
		CountUsage:     2,
		SizeUsage:      2,
		SequenceLength: 2,
		MaxCost: LifeCyclePolicyPolicyUsageCost{
			CurrencyCode:  LifeCyclePolicyCurrency{Currency: &currency.RUB},
			PricePerHour:  decimal.New(0, 0),
			PricePerMonth: decimal.New(0, 0),
		},
	}

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, err := json.Marshal(expResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `%s`, string(resp))
	})

	estimateIntervalRequest := LifeCyclePolicyEstimateIntervalRequest{}

	respActual, resp, err := client.LifeCyclePolicies.EstimateIntervalMaxPolicyUsage(ctx, &estimateIntervalRequest)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, expResp, *respActual)
}
