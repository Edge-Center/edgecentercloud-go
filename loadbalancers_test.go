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

func TestLoadbalancers_List(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []Loadbalancer{{ID: testResourceID}}
	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.List(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_Get(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &Loadbalancer{ID: testResourceID}
	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.Get(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_Create(t *testing.T) {
	setup()
	defer teardown()

	request := &LoadbalancerCreateRequest{
		Name:   "test-loadbalancer",
		Flavor: "lb1-1-2",
	}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(LoadbalancerCreateRequest)
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

	respActual, resp, err := client.Loadbalancers.Create(ctx, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_Delete(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.Delete(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_ListenerGet(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &Listener{ID: testResourceID}
	URL := path.Join(lblistenersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.ListenerGet(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_ListenerCreate(t *testing.T) {
	setup()
	defer teardown()

	request := &ListenerCreateRequest{
		Name:           "test-loadbalancer",
		Protocol:       ListenerProtocolTCP,
		ProtocolPort:   80,
		LoadBalancerID: testResourceID,
	}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(lblistenersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(ListenerCreateRequest)
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

	respActual, resp, err := client.Loadbalancers.ListenerCreate(ctx, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_ListenerDelete(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(lblistenersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.ListenerDelete(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_PoolGet(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &Pool{ID: testResourceID}
	URL := path.Join(lbpoolsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.PoolGet(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_PoolCreate(t *testing.T) {
	setup()
	defer teardown()

	request := &PoolCreateRequest{
		LoadbalancerPoolCreateRequest: LoadbalancerPoolCreateRequest{
			Name:       "test-loadbalancer",
			ListenerID: testResourceID,
		},
	}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(lbpoolsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(PoolCreateRequest)
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

	respActual, resp, err := client.Loadbalancers.PoolCreate(ctx, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_PoolDelete(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(lbpoolsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.PoolDelete(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_PoolUpdate(t *testing.T) {
	setup()
	defer teardown()

	request := &PoolUpdateRequest{
		ID:   testResourceID,
		Name: "test-lbpool",
	}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(lbpoolsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		reqBody := new(PoolUpdateRequest)
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

	respActual, resp, err := client.Loadbalancers.PoolUpdate(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_PoolList(t *testing.T) {
	setup()
	defer teardown()

	options := PoolListOptions{LoadBalancerID: testResourceID}
	expectedResp := []Pool{{ID: testResourceID}}
	URL := path.Join(lbpoolsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.PoolList(ctx, &options)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_CheckLimits(t *testing.T) {
	setup()
	defer teardown()

	request := &LoadbalancerCheckLimitsRequest{}
	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), loadbalancersCheckLimits)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, err := json.Marshal(map[string]int{})
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	_, resp, err := client.Loadbalancers.CheckLimits(ctx, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestLoadbalancers_FlavorList(t *testing.T) {
	setup()
	defer teardown()

	options := FlavorsOptions{IncludePrices: true}
	expectedResp := []Flavor{{FlavorID: "g1-gpu-1-2-1"}}
	URL := path.Join(lbflavorsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.FlavorList(ctx, &options)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_MetadataList(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []MetadataDetailed{{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}}
	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.MetadataList(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_MetadataCreate(t *testing.T) {
	setup()
	defer teardown()

	request := &MetadataCreateRequest{Metadata: map[string]interface{}{"key": "value"}}
	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
	})

	resp, err := client.Loadbalancers.MetadataCreate(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestLoadbalancers_MetadataUpdate(t *testing.T) {
	setup()
	defer teardown()

	request := &MetadataCreateRequest{Metadata: map[string]interface{}{"key": "value"}}
	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
	})

	resp, err := client.Loadbalancers.MetadataUpdate(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestLoadbalancers_MetadataDeleteItem(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataItemPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	resp, err := client.Loadbalancers.MetadataDeleteItem(ctx, testResourceID, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestLoadbalancers_MetadataGetItem(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &MetadataDetailed{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}
	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataItemPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.MetadataGetItem(ctx, testResourceID, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}
