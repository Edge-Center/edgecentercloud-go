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

func TestNetworks_Get(t *testing.T) {
	setup()
	defer teardown()

	const (
		networkID = "f0d19cec-5c3f-4853-886e-304915960ff6"
		projectID = "27520"
		regionID  = "8"
	)

	network := &Network{ID: networkID}
	getNetworkURL := fmt.Sprintf("/v1/networks/%s/%s/%s", projectID, regionID, networkID)

	mux.HandleFunc(getNetworkURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(network)
		_, _ = fmt.Fprintf(w, `{"network":%s}`, string(resp))
	})

	opts := ServicePath{Project: projectID, Region: regionID}
	resp, _, err := client.Networks.Get(ctx, networkID, &opts)
	if err != nil {
		t.Errorf("Instances.Get returned error: %v", err)
	}

	if !reflect.DeepEqual(resp, network) {
		t.Errorf("Networks.Get\n returned %+v,\n expected %+v", resp, network)
	}
}

func TestNetworks_Create(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID    = "f0d19cec-5c3f-4853-886e-304915960ff6"
		projectID = "27520"
		regionID  = "8"
	)

	networkCreateRequest := &NetworkCreateRequest{
		Name:         "test-instance",
		CreateRouter: false,
	}

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	createNetworkURL := fmt.Sprintf("/v1/networks/%s/%s", projectID, regionID)

	mux.HandleFunc(createNetworkURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(NetworkCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, networkCreateRequest, reqBody)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprintf(w, `{"tasks":%s}`, string(resp))
	})

	opts := ServicePath{Project: projectID, Region: regionID}
	resp, _, err := client.Networks.Create(ctx, networkCreateRequest, &opts)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestNetworks_Delete(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID    = "f0d19cec-5c3f-4853-886e-304915960ff6"
		networkID = "f0d19cec-5c3f-4853-886e-304915960ff6"
		projectID = "27520"
		regionID  = "8"
	)

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	deleteNetworkURL := fmt.Sprintf("/v1/networks/%s/%s/%s", projectID, regionID, networkID)

	mux.HandleFunc(deleteNetworkURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprintf(w, `{"tasks":%s}`, string(resp))
	})

	opts := ServicePath{Project: projectID, Region: regionID}
	resp, _, err := client.Networks.Delete(ctx, networkID, &opts)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}
