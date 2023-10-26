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

func TestNetworks_List(t *testing.T) {
	setup()
	defer teardown()

	const (
		networkID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	networks := []Network{{ID: networkID}}
	getNetworksURL := fmt.Sprintf("/v1/networks/%d/%d", projectID, regionID)

	mux.HandleFunc(getNetworksURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(networks)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.Networks.List(ctx, nil)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, networks) {
		t.Errorf("Networks.List\n returned %+v,\n expected %+v", resp, networks)
	}
}

func TestNetworks_Get(t *testing.T) {
	setup()
	defer teardown()

	const (
		networkID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	network := &Network{ID: networkID}
	getNetworkURL := fmt.Sprintf("/v1/networks/%d/%d/%s", projectID, regionID, networkID)

	mux.HandleFunc(getNetworkURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(network)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Networks.Get(ctx, networkID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, network) {
		t.Errorf("Networks.Get\n returned %+v,\n expected %+v", resp, network)
	}
}

func TestNetworks_Create(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	networkCreateRequest := &NetworkCreateRequest{
		Name:         "test-instance",
		CreateRouter: false,
	}

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	createNetworkURL := fmt.Sprintf("/v1/networks/%d/%d", projectID, regionID)

	mux.HandleFunc(createNetworkURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(NetworkCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, networkCreateRequest, reqBody)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Networks.Create(ctx, networkCreateRequest)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestNetworks_Delete(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID    = "f0d19cec-5c3f-4853-886e-304915960ff6"
		networkID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	deleteNetworkURL := fmt.Sprintf("/v1/networks/%d/%d/%s", projectID, regionID, networkID)

	mux.HandleFunc(deleteNetworkURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Networks.Delete(ctx, networkID)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestNetworks_UpdateName(t *testing.T) {
	setup()
	defer teardown()

	const (
		networkID = "f0d19cec-5c3f-4853-886e-304915960ff6"
		newName   = "new-network-name"
	)

	networkUpdateNameRequest := &NetworkUpdateNameRequest{
		Name: newName,
	}

	networkResponse := &Network{Name: newName}

	updateNameNetworkURL := fmt.Sprintf("/v1/networks/%d/%d/%s", projectID, regionID, networkID)

	mux.HandleFunc(updateNameNetworkURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		reqBody := new(NetworkUpdateNameRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, networkUpdateNameRequest, reqBody)
		resp, _ := json.Marshal(networkResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Networks.UpdateName(ctx, networkID, networkUpdateNameRequest)
	require.NoError(t, err)

	assert.Equal(t, networkResponse, resp)
}

func Test_ListNetworksWithSubnets(t *testing.T) {
	setup()
	defer teardown()

	const (
		networkID    = "f0d19cec-5c3f-4853-886e-304915960ff6"
		subnetworkID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	networksSubnetworks := []NetworkSubnetwork{{ID: networkID, Subnets: []Subnetwork{{ID: subnetworkID}}}}
	availableNetworksURL := fmt.Sprintf("/v1/availablenetworks/%d/%d", projectID, regionID)

	mux.HandleFunc(availableNetworksURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(networksSubnetworks)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.Networks.ListNetworksWithSubnets(ctx, nil)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, networksSubnetworks) {
		t.Errorf("Networks.ListNetworksWithSubnets\n returned %+v,\n expected %+v", resp, networksSubnetworks)
	}
}
