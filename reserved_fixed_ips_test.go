package edgecloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReservedFixedIPs_List(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []ReservedFixedIP{{Name: "test-ReservedFixedIP"}}
	URL := path.Join(reservedFixedIPsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.ReservedFixedIP.List(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestReservedFixedIPs_Create(t *testing.T) {
	setup()
	defer teardown()

	request := &ReservedFixedIPCreateRequest{}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(reservedFixedIPsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(ReservedFixedIPCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, reqBody)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.ReservedFixedIP.Create(ctx, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestReservedFixedIPs_Get(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &ReservedFixedIP{Name: "test-ReservedFixedIP"}
	URL := path.Join(reservedFixedIPsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.ReservedFixedIP.Get(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestReservedFixedIPs_Delete(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(reservedFixedIPsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.ReservedFixedIP.Delete(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestReservedFixedIPs_SwitchVIPStatus(t *testing.T) {
	setup()
	defer teardown()

	request := &SwitchVIPStatusRequest{}
	expectedResp := &ReservedFixedIP{Name: "test-ReservedFixedIP"}
	URL := path.Join(reservedFixedIPsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		reqBody := new(SwitchVIPStatusRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, reqBody)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.ReservedFixedIP.SwitchVIPStatus(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestReservedFixedIPs_ListInstancePorts(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []InstancePort{{InstanceID: testResourceID}}
	URL := path.Join(reservedFixedIPsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, reservedFixedIPsConnectedDevices)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.ReservedFixedIP.ListInstancePorts(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestReservedFixedIPs_AddInstancePorts(t *testing.T) {
	setup()
	defer teardown()

	request := &AddInstancePortsRequest{}
	expectedResp := []InstancePort{{InstanceID: testResourceID}}
	URL := path.Join(reservedFixedIPsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, reservedFixedIPsConnectedDevices)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		reqBody := new(AddInstancePortsRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, reqBody)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.ReservedFixedIP.AddInstancePorts(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestReservedFixedIPs_ReplaceInstancePorts(t *testing.T) {
	setup()
	defer teardown()

	request := &AddInstancePortsRequest{}
	expectedResp := []InstancePort{{InstanceID: testResourceID}}
	URL := path.Join(reservedFixedIPsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, reservedFixedIPsConnectedDevices)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		reqBody := new(AddInstancePortsRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, reqBody)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.ReservedFixedIP.ReplaceInstancePorts(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestReservedFixedIPs_ListInstancePortsAvailable(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []InstancePort{{InstanceID: testResourceID}}
	URL := path.Join(reservedFixedIPsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, reservedFixedIPsAvailableDevices)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.ReservedFixedIP.ListInstancePortsAvailable(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}
