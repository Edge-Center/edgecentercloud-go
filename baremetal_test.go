package edgecloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBareMetalServiceOp_CheckQuotasForInstanceCreation(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(bmInstancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), bmCheckLimitsSupPath)
	expResp := Quota(map[string]int{"test": 777})

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, err := json.Marshal(expResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `%s`, string(resp))
	})

	quotaCheckRequest := BareMetalQuotaCheckRequest{}

	respActual, resp, err := client.Instances.BareMetalCheckQuotasForInstanceCreation(ctx, &quotaCheckRequest)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, expResp, respActual)
}

func TestBareMetalServiceOp_CreateInstance(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(bmInstancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))
	expResp := TaskResponse{Tasks: []string{taskID}}

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, err := json.Marshal(expResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `%s`, string(resp))
	})

	bmInstanceCreateRequest := BareMetalServerCreateRequest{}

	respActual, resp, err := client.Instances.BareMetalCreateInstance(ctx, &bmInstanceCreateRequest)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, expResp, *respActual)
}

func TestBareMetalServiceOp_GetCountAvailableNodes(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(bmCapacityBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))
	expResp := BareMetalCapacity{Capacity: map[string]int{"test": 777}}

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `%s`, string(resp))
	})

	respActual, resp, err := client.Instances.BareMetalGetCountAvailableNodes(ctx)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, expResp, *respActual)
}

func TestBareMetalServiceOp_ListFlavors(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(bmInstancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), bmAvailableFlavorsSubPath)
	expResp := []BareMetalFlavor{{FlavorName: "test"}}
	expRespRoot := bareMetalFlavorsRoot{Count: 1, Flavors: expResp}

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, err := json.Marshal(expRespRoot)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `%s`, string(resp))
	})

	opts := BareMetalFlavorsOpts{}
	reqBody := BareMetalFlavorsRequest{}
	respActual, resp, err := client.Instances.BareMetalListFlavors(ctx, &opts, &reqBody)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, expResp, respActual)
}

func TestBareMetalServiceOp_ListInstances(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(bmInstancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))
	expResp := []Instance{{Name: "test"}}
	expRespRoot := instancesRoot{Count: 1, Instances: expResp}

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expRespRoot)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `%s`, string(resp))
	})

	opts := BareMetalInstancesListOpts{}
	respActual, resp, err := client.Instances.BareMetalListInstances(ctx, &opts)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, expResp, respActual)
}

func TestBareMetalServiceOp_RebuildInstance(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(bmInstancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, bmRebuildSubPath)
	expResp := TaskResponse{Tasks: []string{taskID}}
	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, err := json.Marshal(expResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `%s`, string(resp))
	})

	reqBody := BareMetalRebuildRequest{}
	respActual, resp, err := client.Instances.BareMetalRebuildInstance(ctx, testResourceID, &reqBody)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, expResp, *respActual)
}
