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

func TestLoadbalancers_Get(t *testing.T) {
	setup()
	defer teardown()

	const (
		loadbalancerID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	loadbalancer := &Loadbalancer{ID: loadbalancerID}
	getLoadbalancerURL := fmt.Sprintf("/v1/loadbalancers/%d/%d/%s", projectID, regionID, loadbalancerID)

	mux.HandleFunc(getLoadbalancerURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(loadbalancer)
		_, _ = fmt.Fprintf(w, `{"loadbalancer":%s}`, string(resp))
	})

	resp, _, err := client.Loadbalancers.Get(ctx, loadbalancerID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, loadbalancer) {
		t.Errorf("Loadbalancers.Get\n returned %+v,\n expected %+v", resp, loadbalancer)
	}
}

func TestLoadbalancers_Create(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	loadbalancerCreateRequest := &LoadbalancerCreateRequest{
		Name:   "test-loadbalancer",
		Flavor: "lb1-1-2",
	}

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	createLoadbalancerURL := fmt.Sprintf("/v1/loadbalancers/%d/%d", projectID, regionID)

	mux.HandleFunc(createLoadbalancerURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(LoadbalancerCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, loadbalancerCreateRequest, reqBody)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprintf(w, `{"tasks":%s}`, string(resp))
	})

	resp, _, err := client.Loadbalancers.Create(ctx, loadbalancerCreateRequest)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestLoadbalancers_Delete(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID         = "f0d19cec-5c3f-4853-886e-304915960ff6"
		loadbalancerID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	deleteLoadbalancerURL := fmt.Sprintf("/v1/loadbalancers/%d/%d/%s", projectID, regionID, loadbalancerID)

	mux.HandleFunc(deleteLoadbalancerURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprintf(w, `{"tasks":%s}`, string(resp))
	})

	resp, _, err := client.Loadbalancers.Delete(ctx, loadbalancerID)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestLoadbalancers_ListenerGet(t *testing.T) {
	setup()
	defer teardown()

	const (
		listenerID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	listener := &Listener{ID: listenerID}
	getLBListenersURL := fmt.Sprintf("/v1/lblisteners/%d/%d/%s", projectID, regionID, listenerID)

	mux.HandleFunc(getLBListenersURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(listener)
		_, _ = fmt.Fprintf(w, `{"listener":%s}`, string(resp))
	})

	resp, _, err := client.Loadbalancers.ListenerGet(ctx, listenerID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, listener) {
		t.Errorf("Loadbalancers.ListenerGet\n returned %+v,\n expected %+v", resp, listener)
	}
}

func TestLoadbalancers_ListenerCreate(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID         = "f0d19cec-5c3f-4853-886e-304915960ff6"
		loadbalancerID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	listenerCreateRequest := &ListenerCreateRequest{
		Name:           "test-loadbalancer",
		Protocol:       ListenerProtocolTCP,
		ProtocolPort:   80,
		LoadBalancerID: loadbalancerID,
	}

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	createLBListenersURL := fmt.Sprintf("/v1/lblisteners/%d/%d", projectID, regionID)

	mux.HandleFunc(createLBListenersURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(ListenerCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, listenerCreateRequest, reqBody)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprintf(w, `{"tasks":%s}`, string(resp))
	})

	resp, _, err := client.Loadbalancers.ListenerCreate(ctx, listenerCreateRequest)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestLoadbalancers_ListenerDelete(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID     = "f0d19cec-5c3f-4853-886e-304915960ff6"
		listenerID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	deleteLBListenersURL := fmt.Sprintf("/v1/lblisteners/%d/%d/%s", projectID, regionID, listenerID)

	mux.HandleFunc(deleteLBListenersURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprintf(w, `{"tasks":%s}`, string(resp))
	})

	resp, _, err := client.Loadbalancers.ListenerDelete(ctx, listenerID)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestLoadbalancers_PoolGet(t *testing.T) {
	setup()
	defer teardown()

	const (
		poolID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	pool := &Pool{ID: poolID}
	getLBPoolsURL := fmt.Sprintf("/v1/lbpools/%d/%d/%s", projectID, regionID, poolID)

	mux.HandleFunc(getLBPoolsURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(pool)
		_, _ = fmt.Fprintf(w, `{"pool":%s}`, string(resp))
	})

	resp, _, err := client.Loadbalancers.PoolGet(ctx, poolID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, pool) {
		t.Errorf("Loadbalancers.PoolGet\n returned %+v,\n expected %+v", resp, pool)
	}
}

func TestLoadbalancers_PoolCreate(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID     = "f0d19cec-5c3f-4853-886e-304915960ff6"
		listenerID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	poolCreateRequest := &PoolCreateRequest{
		LoadbalancerPoolCreateRequest: LoadbalancerPoolCreateRequest{
			Name:       "test-loadbalancer",
			ListenerID: listenerID,
		},
	}

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	createLBPoolsURL := fmt.Sprintf("/v1/lbpools/%d/%d", projectID, regionID)

	mux.HandleFunc(createLBPoolsURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(PoolCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, poolCreateRequest, reqBody)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprintf(w, `{"tasks":%s}`, string(resp))
	})

	resp, _, err := client.Loadbalancers.PoolCreate(ctx, poolCreateRequest)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestLoadbalancers_PoolDelete(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID = "f0d19cec-5c3f-4853-886e-304915960ff6"
		poolID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	deleteLBPoolsURL := fmt.Sprintf("/v1/lbpools/%d/%d/%s", projectID, regionID, poolID)

	mux.HandleFunc(deleteLBPoolsURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprintf(w, `{"tasks":%s}`, string(resp))
	})

	resp, _, err := client.Loadbalancers.PoolDelete(ctx, poolID)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestLoadbalancers_PoolUpdate(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID = "f0d19cec-5c3f-4853-886e-304915960ff6"
		poolID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	poolUpdateRequest := &PoolUpdateRequest{
		ID:   poolID,
		Name: "test-lbpool",
	}

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	createLBPoolsURL := fmt.Sprintf("/v1/lbpools/%d/%d/%s", projectID, regionID, poolID)

	mux.HandleFunc(createLBPoolsURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		reqBody := new(PoolUpdateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, poolUpdateRequest, reqBody)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprintf(w, `{"tasks":%s}`, string(resp))
	})

	resp, _, err := client.Loadbalancers.PoolUpdate(ctx, poolID, poolUpdateRequest)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestLoadbalancers_PoolList(t *testing.T) {
	setup()
	defer teardown()

	const (
		poolID         = "f0d19cec-5c3f-4853-886e-304915960ff6"
		loadbalancerID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	poolListOptions := PoolListOptions{
		LoadBalancerID: loadbalancerID,
	}

	pools := []Pool{{ID: poolID}}

	poolsListURL := fmt.Sprintf("/v1/lbpools/%d/%d", projectID, regionID)
	mux.HandleFunc(poolsListURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(pools)
		_, _ = fmt.Fprintf(w, `{"pools":%s}`, string(resp))
	})

	resp, _, err := client.Loadbalancers.PoolList(ctx, &poolListOptions)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, pools) {
		t.Errorf("Loadbalancers.PoolList\n returned %+v,\n expected %+v", resp, pools)
	}
}
