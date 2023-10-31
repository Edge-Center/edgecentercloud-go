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

func TestSubnetworks_List(t *testing.T) {
	setup()
	defer teardown()

	const (
		subnetworkID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	subnetworks := []Subnetwork{{ID: subnetworkID}}
	getSubnetworkURL := fmt.Sprintf("/v1/subnets/%d/%d", projectID, regionID)

	mux.HandleFunc(getSubnetworkURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(subnetworks)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.Subnetworks.List(ctx, nil)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, subnetworks) {
		t.Errorf("Subnetworks.List\n returned %+v,\n expected %+v", resp, subnetworks)
	}
}

func TestSubnetworks_Get(t *testing.T) {
	setup()
	defer teardown()

	const (
		subnetworkID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	subnetwork := &Subnetwork{ID: subnetworkID}
	getSubnetworkURL := fmt.Sprintf("/v1/subnets/%d/%d/%s", projectID, regionID, subnetworkID)

	mux.HandleFunc(getSubnetworkURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(subnetwork)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Subnetworks.Get(ctx, subnetworkID)
	require.NoError(t, err)

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
	)

	subnetworkCreateRequest := &SubnetworkCreateRequest{
		Name:      "test-subnet",
		NetworkID: networkID,
	}

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	createSubnetworkURL := fmt.Sprintf("/v1/subnets/%d/%d", projectID, regionID)

	mux.HandleFunc(createSubnetworkURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(SubnetworkCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, subnetworkCreateRequest, reqBody)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Subnetworks.Create(ctx, subnetworkCreateRequest)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestSubnetworks_Delete(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID       = "f0d19cec-5c3f-4853-886e-304915960ff6"
		subnetworkID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	deleteSubnetworkURL := fmt.Sprintf("/v1/subnets/%d/%d/%s", projectID, regionID, subnetworkID)

	mux.HandleFunc(deleteSubnetworkURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Subnetworks.Delete(ctx, subnetworkID)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestSubnetworks_Update(t *testing.T) {
	setup()
	defer teardown()

	const (
		subnetID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	subnetworkUpdateRequest := &SubnetworkUpdateRequest{
		Name: "test-subnet",
	}

	subnetworkResponse := &Subnetwork{ID: subnetID}

	updateSubnetworkURL := fmt.Sprintf("/v1/subnets/%d/%d/%s", projectID, regionID, subnetID)

	mux.HandleFunc(updateSubnetworkURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		reqBody := new(SubnetworkUpdateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, subnetworkUpdateRequest, reqBody)
		resp, _ := json.Marshal(subnetworkResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Subnetworks.Update(ctx, subnetID, subnetworkUpdateRequest)
	require.NoError(t, err)

	assert.Equal(t, subnetworkResponse, resp)
}

func TestSubnetworks_MetadataList(t *testing.T) {
	setup()
	defer teardown()

	const (
		subnetworkID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	metadataList := []MetadataDetailed{{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}}
	getSubnetworksMetadataListURL := fmt.Sprintf("/v1/subnets/%d/%d/%s/%s", projectID, regionID, subnetworkID, metadataPath)

	mux.HandleFunc(getSubnetworksMetadataListURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(metadataList)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.Subnetworks.MetadataList(ctx, subnetworkID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, metadataList) {
		t.Errorf("Subnetworks.MetadataList\n returned %+v,\n expected %+v", resp, metadataList)
	}
}

func TestSubnetworks_MetadataCreate(t *testing.T) {
	setup()
	defer teardown()

	const (
		subnetworkID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	metadataCreateRequest := &MetadataCreateRequest{
		map[string]interface{}{"key": "value"},
	}

	createSubnetworksMetadataURL := fmt.Sprintf("/v1/subnets/%d/%d/%s/%s", projectID, regionID, subnetworkID, metadataPath)

	mux.HandleFunc(createSubnetworksMetadataURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
	})

	_, err := client.Subnetworks.MetadataCreate(ctx, subnetworkID, metadataCreateRequest)
	require.NoError(t, err)
}

func TestSubnetworks_MetadataUpdate(t *testing.T) {
	setup()
	defer teardown()

	const (
		subnetworkID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	metadataCreateRequest := &MetadataCreateRequest{
		map[string]interface{}{"key": "value"},
	}

	createSubnetworksMetadataURL := fmt.Sprintf("/v1/subnets/%d/%d/%s/%s", projectID, regionID, subnetworkID, metadataPath)

	mux.HandleFunc(createSubnetworksMetadataURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
	})

	_, err := client.Subnetworks.MetadataUpdate(ctx, subnetworkID, metadataCreateRequest)
	require.NoError(t, err)
}

func TestSubnetworks_MetadataDeleteItem(t *testing.T) {
	setup()
	defer teardown()

	const (
		subnetworkID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	deleteSubnetworksMetadataItemURL := fmt.Sprintf("/v1/subnets/%d/%d/%s/%s", projectID, regionID, subnetworkID, metadataItemPath)

	mux.HandleFunc(deleteSubnetworksMetadataItemURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.Subnetworks.MetadataDeleteItem(ctx, subnetworkID, nil)
	require.NoError(t, err)
}

func TestSubnetworks_MetadataGetItem(t *testing.T) {
	setup()
	defer teardown()

	const (
		subnetworkID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	metadata := &MetadataDetailed{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}
	getSubnetworksMetadataItemURL := fmt.Sprintf("/v1/subnets/%d/%d/%s/%s", projectID, regionID, subnetworkID, metadataItemPath)

	mux.HandleFunc(getSubnetworksMetadataItemURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(metadata)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Subnetworks.MetadataGetItem(ctx, subnetworkID, nil)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, metadata) {
		t.Errorf("Subnetworks.MetadataGetItem\n returned %+v,\n expected %+v", resp, metadata)
	}
}
