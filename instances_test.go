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

func TestInstances_Get(t *testing.T) {
	setup()
	defer teardown()

	const (
		instanceID = "f0d19cec-5c3f-4853-886e-304915960ff6"
		projectID  = "27520"
		regionID   = "8"
	)

	instance := &Instance{ID: instanceID}
	getInstanceURL := fmt.Sprintf("/v1/instances/%s/%s/%s", projectID, regionID, instanceID)

	mux.HandleFunc(getInstanceURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(instance)
		_, _ = fmt.Fprintf(w, `{"instance":%s}`, string(resp))
	})

	opts := ServicePath{Project: projectID, Region: regionID}
	resp, _, err := client.Instances.Get(ctx, instanceID, &opts)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, instance) {
		t.Errorf("Instances.Get\n returned %+v,\n expected %+v", resp, instance)
	}
}

func TestInstances_Create(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID    = "f0d19cec-5c3f-4853-886e-304915960ff6"
		projectID = "27520"
		regionID  = "8"
	)

	instanceCreateRequest := &InstanceCreateRequest{
		Names:          []string{"test-instance"},
		Flavor:         "g1-standard-1-2",
		Interfaces:     []InstanceInterface{{Type: ExternalInterfaceType}},
		SecurityGroups: []InstanceSecurityGroupsCreate{{ID: "f0d19cec-5c3f-4853-886e-304915960ff6"}},
		Volumes: []InstanceVolumeCreate{
			{
				TypeName:  SsdHiIops,
				Size:      5,
				BootIndex: 0,
				Source:    Image,
				ImageID:   "f0d19cec-5c3f-4853-886e-304915960ff6",
			},
		},
	}

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	createInstanceURL := fmt.Sprintf("/v2/instances/%s/%s", projectID, regionID)

	mux.HandleFunc(createInstanceURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(InstanceCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, instanceCreateRequest, reqBody)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprintf(w, `{"tasks":%s}`, string(resp))
	})

	opts := ServicePath{Project: projectID, Region: regionID}
	resp, _, err := client.Instances.Create(ctx, instanceCreateRequest, &opts)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestInstances_Delete(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID     = "f0d19cec-5c3f-4853-886e-304915960ff6"
		instanceID = "f0d19cec-5c3f-4853-886e-304915960ff6"
		projectID  = "27520"
		regionID   = "8"
	)

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	deleteInstanceURL := fmt.Sprintf("/v1/instances/%s/%s/%s", projectID, regionID, instanceID)

	mux.HandleFunc(deleteInstanceURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprintf(w, `{"tasks":%s}`, string(resp))
	})

	opts := ServicePath{Project: projectID, Region: regionID}
	resp, _, err := client.Instances.Delete(ctx, instanceID, &opts, nil)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestInstances_DeleteWithOptions(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID     = "f0d19cec-5c3f-4853-886e-304915960ff6"
		instanceID = "f0d19cec-5c3f-4853-886e-304915960ff6"
		projectID  = "27520"
		regionID   = "8"
	)

	taskResponse := &TaskResponse{Tasks: []string{taskID}}
	instanceDeleteOptions := InstanceDeleteOptions{
		DeleteFloatings: true,
		Volumes:         []string{"f0d19cec-5c3f-4853-886e-304915960ff6"},
	}

	deleteInstanceURL := fmt.Sprintf("/v1/instances/%s/%s/%s", projectID, regionID, instanceID)
	mux.HandleFunc(deleteInstanceURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprintf(w, `{"tasks":%s}`, string(resp))
	})

	opts := ServicePath{Project: projectID, Region: regionID}
	resp, _, err := client.Instances.Delete(ctx, instanceID, &opts, &instanceDeleteOptions)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestInstances_MetadataGet(t *testing.T) {
	setup()
	defer teardown()

	const (
		instanceID = "f0d19cec-5c3f-4853-886e-304915960ff6"
		projectID  = "27520"
		regionID   = "8"
	)

	metadata := &MetadataDetailed{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}
	getInstanceMetadataURL := fmt.Sprintf("/v1/instances/%s/%s/%s/%s", projectID, regionID, instanceID, metadataPath)

	mux.HandleFunc(getInstanceMetadataURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(metadata)
		_, _ = fmt.Fprintf(w, `{"metadata":%s}`, string(resp))
	})

	opts := ServicePath{Project: projectID, Region: regionID}
	resp, _, err := client.Instances.MetadataGet(ctx, instanceID, &opts)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, metadata) {
		t.Errorf("Instances.MetadataGet\n returned %+v,\n expected %+v", resp, metadata)
	}
}

func TestInstances_MetadataCreate(t *testing.T) {
	setup()
	defer teardown()

	const (
		instanceID = "f0d19cec-5c3f-4853-886e-304915960ff6"
		projectID  = "27520"
		regionID   = "8"
	)

	metadataCreateRequest := &MetadataCreateRequest{
		map[string]interface{}{"key": "value"},
	}

	createInstanceMetadataURL := fmt.Sprintf("/v1/instances/%s/%s/%s/%s", projectID, regionID, instanceID, metadataPath)

	mux.HandleFunc(createInstanceMetadataURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
	})

	opts := ServicePath{Project: projectID, Region: regionID}
	_, err := client.Instances.MetadataCreate(ctx, instanceID, metadataCreateRequest, &opts)
	require.NoError(t, err)
}
