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
		projectID      = "27520"
		regionID       = "8"
	)

	loadbalancer := &Loadbalancer{ID: loadbalancerID}
	getLoadbalancerURL := fmt.Sprintf("/v1/loadbalancers/%s/%s/%s", projectID, regionID, loadbalancerID)

	mux.HandleFunc(getLoadbalancerURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(loadbalancer)
		_, _ = fmt.Fprintf(w, `{"loadbalancer":%s}`, string(resp))
	})

	opts := ServicePath{Project: projectID, Region: regionID}
	resp, _, err := client.Loadbalancers.Get(ctx, loadbalancerID, &opts)
	if err != nil {
		t.Errorf("Loadbalancers.Get returned error: %v", err)
	}

	if !reflect.DeepEqual(resp, loadbalancer) {
		t.Errorf("Loadbalancers.Get\n returned %+v,\n expected %+v", resp, loadbalancer)
	}
}

func TestLoadbalancers_Create(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID    = "f0d19cec-5c3f-4853-886e-304915960ff6"
		projectID = "27520"
		regionID  = "8"
	)

	loadbalancerCreateRequest := &LoadbalancerCreateRequest{
		Name:   "test-loadbalancer",
		Flavor: "g1-standard-1-2",
	}

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	createLoadbalancerURL := fmt.Sprintf("/v1/loadbalancers/%s/%s", projectID, regionID)

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

	opts := ServicePath{Project: projectID, Region: regionID}
	resp, _, err := client.Loadbalancers.Create(ctx, loadbalancerCreateRequest, &opts)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestLoadbalancers_Delete(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID         = "f0d19cec-5c3f-4853-886e-304915960ff6"
		loadbalancerID = "f0d19cec-5c3f-4853-886e-304915960ff6"
		projectID      = "27520"
		regionID       = "8"
	)

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	deleteLoadbalancerURL := fmt.Sprintf("/v1/loadbalancers/%s/%s/%s", projectID, regionID, loadbalancerID)

	mux.HandleFunc(deleteLoadbalancerURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprintf(w, `{"tasks":%s}`, string(resp))
	})

	opts := ServicePath{Project: projectID, Region: regionID}
	resp, _, err := client.Loadbalancers.Delete(ctx, loadbalancerID, &opts)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}
