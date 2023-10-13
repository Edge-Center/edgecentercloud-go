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

func TestSubnetworks_Get(t *testing.T) {
	setup()
	defer teardown()

	const (
		subnetworkID = "f0d19cec-5c3f-4853-886e-304915960ff6"
		projectID    = "27520"
		regionID     = "8"
	)

	subnetwork := &Subnetwork{ID: subnetworkID}
	getSubnetworkURL := fmt.Sprintf("/v1/subnets/%s/%s/%s", projectID, regionID, subnetworkID)

	mux.HandleFunc(getSubnetworkURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(subnetwork)
		_, _ = fmt.Fprintf(w, `{"subnetwork":%s}`, string(resp))
	})

	opts := ServicePath{Project: projectID, Region: regionID}
	resp, _, err := client.Subnetworks.Get(ctx, subnetworkID, &opts)
	if err != nil {
		t.Errorf("Subnetworks.Get returned error: %v", err)
	}

	if !reflect.DeepEqual(resp, subnetwork) {
		t.Errorf("Subnetworks.Get\n returned %+v,\n expected %+v", resp, subnetwork)
	}
}

func TestSubnetworks_Create(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID    = "f0d19cec-5c3f-4853-886e-304915960ff6"
		networkID = "f0d19cec-5c3f-4853-886e-304915960ff6"
		projectID = "27520"
		regionID  = "8"
	)

	subnetworkCreateRequest := &SubnetworkCreateRequest{
		Name:      "test-subnet",
		NetworkID: networkID,
	}

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	createSubnetworkURL := fmt.Sprintf("/v1/subnets/%s/%s", projectID, regionID)

	mux.HandleFunc(createSubnetworkURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(SubnetworkCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, subnetworkCreateRequest, reqBody)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprintf(w, `{"tasks":%s}`, string(resp))
	})

	opts := ServicePath{Project: projectID, Region: regionID}
	resp, _, err := client.Subnetworks.Create(ctx, subnetworkCreateRequest, &opts)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestSubnetworks_Delete(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID       = "f0d19cec-5c3f-4853-886e-304915960ff6"
		subnetworkID = "f0d19cec-5c3f-4853-886e-304915960ff6"
		projectID    = "27520"
		regionID     = "8"
	)

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	deleteNetworkURL := fmt.Sprintf("/v1/subnets/%s/%s/%s", projectID, regionID, subnetworkID)

	mux.HandleFunc(deleteNetworkURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprintf(w, `{"tasks":%s}`, string(resp))
	})

	opts := ServicePath{Project: projectID, Region: regionID}
	resp, _, err := client.Subnetworks.Delete(ctx, subnetworkID, &opts)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}
