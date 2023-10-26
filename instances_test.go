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

func TestInstances_List(t *testing.T) {
	setup()
	defer teardown()

	const (
		instanceID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	instances := []Instance{{ID: instanceID}}
	getInstancesURL := fmt.Sprintf("/v1/instances/%d/%d", projectID, regionID)

	mux.HandleFunc(getInstancesURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(instances)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.Instances.List(ctx, nil)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, instances) {
		t.Errorf("Instances.List\n returned %+v,\n expected %+v", resp, instances)
	}
}

func TestInstances_Get(t *testing.T) {
	setup()
	defer teardown()

	const (
		instanceID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	instance := &Instance{ID: instanceID}
	getInstanceURL := fmt.Sprintf("/v1/instances/%d/%d/%s", projectID, regionID, instanceID)

	mux.HandleFunc(getInstanceURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(instance)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Instances.Get(ctx, instanceID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, instance) {
		t.Errorf("Instances.Get\n returned %+v,\n expected %+v", resp, instance)
	}
}

func TestInstances_Create(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID = "f0d19cec-5c3f-4853-886e-304915960ff6"
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

	createInstanceURL := fmt.Sprintf("/v2/instances/%d/%d", projectID, regionID)

	mux.HandleFunc(createInstanceURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(InstanceCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, instanceCreateRequest, reqBody)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Instances.Create(ctx, instanceCreateRequest)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestInstances_Delete(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID     = "f0d19cec-5c3f-4853-886e-304915960ff6"
		instanceID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	deleteInstanceURL := fmt.Sprintf("/v1/instances/%d/%d/%s", projectID, regionID, instanceID)

	mux.HandleFunc(deleteInstanceURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Instances.Delete(ctx, instanceID, nil)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestInstances_DeleteWithOptions(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID     = "f0d19cec-5c3f-4853-886e-304915960ff6"
		instanceID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	taskResponse := &TaskResponse{Tasks: []string{taskID}}
	instanceDeleteOptions := InstanceDeleteOptions{
		DeleteFloatings: true,
		Volumes:         []string{"f0d19cec-5c3f-4853-886e-304915960ff6"},
	}

	deleteInstanceURL := fmt.Sprintf("/v1/instances/%d/%d/%s", projectID, regionID, instanceID)
	mux.HandleFunc(deleteInstanceURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Instances.Delete(ctx, instanceID, &instanceDeleteOptions)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestInstances_MetadataGet(t *testing.T) {
	setup()
	defer teardown()

	const (
		instanceID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	metadata := &MetadataDetailed{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}
	getInstanceMetadataURL := fmt.Sprintf("/v1/instances/%d/%d/%s/%s", projectID, regionID, instanceID, instanceMetadataPath)

	mux.HandleFunc(getInstanceMetadataURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(metadata)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Instances.MetadataGet(ctx, instanceID)
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
	)

	metadataCreateRequest := &MetadataCreateRequest{
		map[string]interface{}{"key": "value"},
	}

	createInstanceMetadataURL := fmt.Sprintf("/v1/instances/%d/%d/%s/%s", projectID, regionID, instanceID, instanceMetadataPath)

	mux.HandleFunc(createInstanceMetadataURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
	})

	_, err := client.Instances.MetadataCreate(ctx, instanceID, metadataCreateRequest)
	require.NoError(t, err)
}

func TestInstances_CheckLimits(t *testing.T) {
	setup()
	defer teardown()

	checkLimitsRequest := &InstanceCheckLimitsRequest{
		Names:      []string{""},
		Interfaces: []InstanceInterface{},
		Volumes:    []InstanceCheckLimitsVolume{},
	}

	checkLimitsInstanceURL := fmt.Sprintf("/v2/instances/%d/%d/%s", projectID, regionID, instancesCheckLimitsPath)

	mux.HandleFunc(checkLimitsInstanceURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, _ := json.Marshal(map[string]int{})
		_, _ = fmt.Fprint(w, string(resp))
	})

	_, _, err := client.Instances.CheckLimits(ctx, checkLimitsRequest)
	require.NoError(t, err)
}

func TestInstances_UpdateFlavor(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID     = "f0d19cec-5c3f-4853-886e-304915960ff6"
		instanceID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	instanceFlavorUpdateRequest := &InstanceFlavorUpdateRequest{
		FlavorID: "g1-standard-1-2",
	}

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	instancesChangeFlavorURL := fmt.Sprintf("/v1/instances/%d/%d/%s/%s", projectID, regionID, instanceID, instancesChangeFlavorPath)

	mux.HandleFunc(instancesChangeFlavorURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Instances.UpdateFlavor(ctx, instanceID, instanceFlavorUpdateRequest)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}
