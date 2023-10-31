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

	networks := []Network{{ID: testResourceID}}
	URL := fmt.Sprintf("/v1/networks/%d/%d", projectID, regionID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
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

	network := &Network{ID: testResourceID}
	URL := fmt.Sprintf("/v1/networks/%d/%d/%s", projectID, regionID, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(network)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Networks.Get(ctx, testResourceID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, network) {
		t.Errorf("Networks.Get\n returned %+v,\n expected %+v", resp, network)
	}
}

func TestNetworks_Create(t *testing.T) {
	setup()
	defer teardown()

	networkCreateRequest := &NetworkCreateRequest{
		Name:         "test-instance",
		CreateRouter: false,
	}
	taskResponse := &TaskResponse{Tasks: []string{taskID}}
	URL := fmt.Sprintf("/v1/networks/%d/%d", projectID, regionID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
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

	taskResponse := &TaskResponse{Tasks: []string{taskID}}
	URL := fmt.Sprintf("/v1/networks/%d/%d/%s", projectID, regionID, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Networks.Delete(ctx, testResourceID)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestNetworks_UpdateName(t *testing.T) {
	setup()
	defer teardown()

	const (
		newName = "new-network-name"
	)

	networkUpdateNameRequest := &NetworkUpdateNameRequest{
		Name: newName,
	}
	networkResponse := &Network{Name: newName}
	URL := fmt.Sprintf("/v1/networks/%d/%d/%s", projectID, regionID, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		reqBody := new(NetworkUpdateNameRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, networkUpdateNameRequest, reqBody)
		resp, _ := json.Marshal(networkResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Networks.UpdateName(ctx, testResourceID, networkUpdateNameRequest)
	require.NoError(t, err)

	assert.Equal(t, networkResponse, resp)
}

func Test_ListNetworksWithSubnets(t *testing.T) {
	setup()
	defer teardown()

	networksSubnetworks := []NetworkSubnetwork{
		{ID: testResourceID, Subnets: []Subnetwork{{ID: testResourceID}}},
	}
	URL := fmt.Sprintf("/v1/availablenetworks/%d/%d", projectID, regionID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
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

func TestNetworks_PortList(t *testing.T) {
	setup()
	defer teardown()

	portsInstances := []PortsInstance{{ID: testResourceID}}
	URL := fmt.Sprintf("/v1/networks/%d/%d/%s/%s", projectID, regionID, testResourceID, networksPortsPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(portsInstances)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.Networks.PortList(ctx, testResourceID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, portsInstances) {
		t.Errorf("Networks.PortList\n returned %+v,\n expected %+v", resp, portsInstances)
	}
}

func TestNetworks_MetadataList(t *testing.T) {
	setup()
	defer teardown()

	metadataList := []MetadataDetailed{{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}}
	URL := fmt.Sprintf("/v1/networks/%d/%d/%s/%s", projectID, regionID, testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(metadataList)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.Networks.MetadataList(ctx, testResourceID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, metadataList) {
		t.Errorf("Networks.MetadataList\n returned %+v,\n expected %+v", resp, metadataList)
	}
}

func TestNetworks_MetadataCreate(t *testing.T) {
	setup()
	defer teardown()

	metadataCreateRequest := &MetadataCreateRequest{
		map[string]interface{}{"key": "value"},
	}
	URL := fmt.Sprintf("/v1/networks/%d/%d/%s/%s", projectID, regionID, testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
	})

	_, err := client.Networks.MetadataCreate(ctx, testResourceID, metadataCreateRequest)
	require.NoError(t, err)
}

func TestNetworks_MetadataUpdate(t *testing.T) {
	setup()
	defer teardown()

	metadataCreateRequest := &MetadataCreateRequest{
		map[string]interface{}{"key": "value"},
	}
	URL := fmt.Sprintf("/v1/networks/%d/%d/%s/%s", projectID, regionID, testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
	})

	_, err := client.Networks.MetadataUpdate(ctx, testResourceID, metadataCreateRequest)
	require.NoError(t, err)
}

func TestNetworks_MetadataDeleteItem(t *testing.T) {
	setup()
	defer teardown()

	URL := fmt.Sprintf("/v1/networks/%d/%d/%s/%s", projectID, regionID, testResourceID, metadataItemPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.Networks.MetadataDeleteItem(ctx, testResourceID, nil)
	require.NoError(t, err)
}

func TestNetworks_MetadataGetItem(t *testing.T) {
	setup()
	defer teardown()

	metadata := &MetadataDetailed{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}
	URL := fmt.Sprintf("/v1/networks/%d/%d/%s/%s", projectID, regionID, testResourceID, metadataItemPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(metadata)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Networks.MetadataGetItem(ctx, testResourceID, nil)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, metadata) {
		t.Errorf("Networks.MetadataGetItem\n returned %+v,\n expected %+v", resp, metadata)
	}
}
