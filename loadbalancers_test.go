package edgecloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadbalancers_List(t *testing.T) {
	setup()
	defer teardown()

	loadbalancers := []Loadbalancer{{ID: testResourceID}}
	URL := fmt.Sprintf("/v1/loadbalancers/%d/%d", projectID, regionID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(loadbalancers)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.Loadbalancers.List(ctx, nil)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, loadbalancers) {
		t.Errorf("Loadbalancers.List\n returned %+v,\n expected %+v", resp, loadbalancers)
	}
}

func TestLoadbalancers_Get(t *testing.T) {
	setup()
	defer teardown()

	loadbalancer := &Loadbalancer{ID: testResourceID}
	URL := fmt.Sprintf("/v1/loadbalancers/%d/%d/%s", projectID, regionID, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(loadbalancer)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Loadbalancers.Get(ctx, testResourceID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, loadbalancer) {
		t.Errorf("Loadbalancers.Get\n returned %+v,\n expected %+v", resp, loadbalancer)
	}
}

func TestLoadbalancers_Create(t *testing.T) {
	setup()
	defer teardown()

	loadbalancerCreateRequest := &LoadbalancerCreateRequest{
		Name:   "test-loadbalancer",
		Flavor: "lb1-1-2",
	}
	taskResponse := &TaskResponse{Tasks: []string{taskID}}
	URL := fmt.Sprintf("/v1/loadbalancers/%d/%d", projectID, regionID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(LoadbalancerCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, loadbalancerCreateRequest, reqBody)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Loadbalancers.Create(ctx, loadbalancerCreateRequest)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestLoadbalancers_Delete(t *testing.T) {
	setup()
	defer teardown()

	taskResponse := &TaskResponse{Tasks: []string{taskID}}
	URL := fmt.Sprintf("/v1/loadbalancers/%d/%d/%s", projectID, regionID, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Loadbalancers.Delete(ctx, testResourceID)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestLoadbalancers_ListenerGet(t *testing.T) {
	setup()
	defer teardown()

	listener := &Listener{ID: testResourceID}
	URL := fmt.Sprintf("/v1/lblisteners/%d/%d/%s", projectID, regionID, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(listener)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Loadbalancers.ListenerGet(ctx, testResourceID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, listener) {
		t.Errorf("Loadbalancers.ListenerGet\n returned %+v,\n expected %+v", resp, listener)
	}
}

func TestLoadbalancers_ListenerCreate(t *testing.T) {
	setup()
	defer teardown()

	listenerCreateRequest := &ListenerCreateRequest{
		Name:           "test-loadbalancer",
		Protocol:       ListenerProtocolTCP,
		ProtocolPort:   80,
		LoadBalancerID: testResourceID,
	}
	taskResponse := &TaskResponse{Tasks: []string{taskID}}
	URL := fmt.Sprintf("/v1/lblisteners/%d/%d", projectID, regionID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(ListenerCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, listenerCreateRequest, reqBody)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Loadbalancers.ListenerCreate(ctx, listenerCreateRequest)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestLoadbalancers_ListenerDelete(t *testing.T) {
	setup()
	defer teardown()

	taskResponse := &TaskResponse{Tasks: []string{taskID}}
	URL := fmt.Sprintf("/v1/lblisteners/%d/%d/%s", projectID, regionID, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Loadbalancers.ListenerDelete(ctx, testResourceID)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestLoadbalancers_PoolGet(t *testing.T) {
	setup()
	defer teardown()

	pool := &Pool{ID: testResourceID}
	URL := fmt.Sprintf("/v1/lbpools/%d/%d/%s", projectID, regionID, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(pool)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Loadbalancers.PoolGet(ctx, testResourceID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, pool) {
		t.Errorf("Loadbalancers.PoolGet\n returned %+v,\n expected %+v", resp, pool)
	}
}

func TestLoadbalancers_PoolCreate(t *testing.T) {
	setup()
	defer teardown()

	poolCreateRequest := &PoolCreateRequest{
		LoadbalancerPoolCreateRequest: LoadbalancerPoolCreateRequest{
			Name:       "test-loadbalancer",
			ListenerID: testResourceID,
		},
	}
	taskResponse := &TaskResponse{Tasks: []string{taskID}}
	URL := fmt.Sprintf("/v1/lbpools/%d/%d", projectID, regionID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(PoolCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, poolCreateRequest, reqBody)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Loadbalancers.PoolCreate(ctx, poolCreateRequest)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestLoadbalancers_PoolDelete(t *testing.T) {
	setup()
	defer teardown()

	taskResponse := &TaskResponse{Tasks: []string{taskID}}
	URL := fmt.Sprintf("/v1/lbpools/%d/%d/%s", projectID, regionID, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Loadbalancers.PoolDelete(ctx, testResourceID)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestLoadbalancers_PoolUpdate(t *testing.T) {
	setup()
	defer teardown()

	poolUpdateRequest := &PoolUpdateRequest{
		ID:   testResourceID,
		Name: "test-lbpool",
	}
	taskResponse := &TaskResponse{Tasks: []string{taskID}}
	URL := fmt.Sprintf("/v1/lbpools/%d/%d/%s", projectID, regionID, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		reqBody := new(PoolUpdateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, poolUpdateRequest, reqBody)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Loadbalancers.PoolUpdate(ctx, testResourceID, poolUpdateRequest)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestLoadbalancers_PoolList(t *testing.T) {
	setup()
	defer teardown()

	poolListOptions := PoolListOptions{
		LoadBalancerID: testResourceID,
	}
	pools := []Pool{{ID: testResourceID}}
	URL := fmt.Sprintf("/v1/lbpools/%d/%d", projectID, regionID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(pools)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.Loadbalancers.PoolList(ctx, &poolListOptions)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, pools) {
		t.Errorf("Loadbalancers.PoolList\n returned %+v,\n expected %+v", resp, pools)
	}
}

func TestLoadbalancers_CheckLimits(t *testing.T) {
	setup()
	defer teardown()

	checkLimitsRequest := &LoadbalancerCheckLimitsRequest{}
	URL := fmt.Sprintf("/v1/loadbalancers/%d/%d/%s", projectID, regionID, loadbalancersCheckLimitsPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, _ := json.Marshal(map[string]int{})
		_, _ = fmt.Fprint(w, string(resp))
	})

	_, _, err := client.Loadbalancers.CheckLimits(ctx, checkLimitsRequest)
	require.NoError(t, err)
}

func TestLoadbalancers_FlavorList(t *testing.T) {
	setup()
	defer teardown()

	loadbalancerFlavorsOptions := FlavorsOptions{
		IncludePrices: true,
	}
	flavors := []Flavor{{FlavorID: "g1-gpu-1-2-1"}}
	URL := fmt.Sprintf("/v1/lbflavors/%d/%d", projectID, regionID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(flavors)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.Loadbalancers.FlavorList(ctx, &loadbalancerFlavorsOptions)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, flavors) {
		t.Errorf("Loadbalancers.FlavorList\n returned %+v,\n expected %+v", resp, flavors)
	}
}

func TestLoadbalancers_MetadataList(t *testing.T) {
	setup()
	defer teardown()

	metadataList := []MetadataDetailed{{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}}
	URL := fmt.Sprintf("/v1/loadbalancers/%d/%d/%s/%s", projectID, regionID, testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(metadataList)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.Loadbalancers.MetadataList(ctx, testResourceID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, metadataList) {
		t.Errorf("Loadbalancers.MetadataList\n returned %+v,\n expected %+v", resp, metadataList)
	}
}

func TestLoadbalancers_MetadataCreate(t *testing.T) {
	setup()
	defer teardown()

	metadataCreateRequest := &MetadataCreateRequest{
		map[string]interface{}{"key": "value"},
	}
	URL := fmt.Sprintf("/v1/loadbalancers/%d/%d/%s/%s", projectID, regionID, testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
	})

	_, err := client.Loadbalancers.MetadataCreate(ctx, testResourceID, metadataCreateRequest)
	require.NoError(t, err)
}

func TestLoadbalancers_MetadataUpdate(t *testing.T) {
	setup()
	defer teardown()

	metadataCreateRequest := &MetadataCreateRequest{
		map[string]interface{}{"key": "value"},
	}
	URL := fmt.Sprintf("/v1/loadbalancers/%d/%d/%s/%s", projectID, regionID, testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
	})

	_, err := client.Loadbalancers.MetadataUpdate(ctx, testResourceID, metadataCreateRequest)
	require.NoError(t, err)
}

func TestLoadbalancers_MetadataDeleteItem(t *testing.T) {
	setup()
	defer teardown()

	URL := fmt.Sprintf("/v1/loadbalancers/%d/%d/%s/%s", projectID, regionID, testResourceID, metadataItemPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.Loadbalancers.MetadataDeleteItem(ctx, testResourceID, nil)
	require.NoError(t, err)
}

func TestLoadbalancers_MetadataGetItem(t *testing.T) {
	setup()
	defer teardown()

	metadata := &MetadataDetailed{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}
	URL := fmt.Sprintf("/v1/loadbalancers/%d/%d/%s/%s", projectID, regionID, testResourceID, metadataItemPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(metadata)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Loadbalancers.MetadataGetItem(ctx, testResourceID, nil)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, metadata) {
		t.Errorf("Loadbalancers.MetadataGetItem\n returned %+v,\n expected %+v", resp, metadata)
	}
}
