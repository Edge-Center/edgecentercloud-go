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

	instances := []Instance{{ID: testResourceID}}
	URL := fmt.Sprintf("/v1/instances/%d/%d", projectID, regionID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
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

	instance := &Instance{ID: testResourceID}
	URL := fmt.Sprintf("/v1/instances/%d/%d/%s", projectID, regionID, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(instance)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Instances.Get(ctx, testResourceID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, instance) {
		t.Errorf("Instances.Get\n returned %+v,\n expected %+v", resp, instance)
	}
}

func TestInstances_Create(t *testing.T) {
	setup()
	defer teardown()

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
	URL := fmt.Sprintf("/v2/instances/%d/%d", projectID, regionID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
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

	taskResponse := &TaskResponse{Tasks: []string{taskID}}
	URL := fmt.Sprintf("/v1/instances/%d/%d/%s", projectID, regionID, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Instances.Delete(ctx, testResourceID, nil)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestInstances_DeleteWithOptions(t *testing.T) {
	setup()
	defer teardown()

	taskResponse := &TaskResponse{Tasks: []string{taskID}}
	instanceDeleteOptions := InstanceDeleteOptions{
		DeleteFloatings: true,
		Volumes:         []string{"f0d19cec-5c3f-4853-886e-304915960ff6"},
	}
	URL := fmt.Sprintf("/v1/instances/%d/%d/%s", projectID, regionID, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Instances.Delete(ctx, testResourceID, &instanceDeleteOptions)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestInstances_MetadataGet(t *testing.T) {
	setup()
	defer teardown()

	metadata := &MetadataDetailed{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}
	URL := fmt.Sprintf("/v1/instances/%d/%d/%s/%s", projectID, regionID, testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(metadata)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Instances.MetadataGet(ctx, testResourceID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, metadata) {
		t.Errorf("Instances.MetadataGet\n returned %+v,\n expected %+v", resp, metadata)
	}
}

func TestInstances_MetadataCreate(t *testing.T) {
	setup()
	defer teardown()

	metadataCreateRequest := &MetadataCreateRequest{
		map[string]interface{}{"key": "value"},
	}
	URL := fmt.Sprintf("/v1/instances/%d/%d/%s/%s", projectID, regionID, testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
	})

	_, err := client.Instances.MetadataCreate(ctx, testResourceID, metadataCreateRequest)
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
	URL := fmt.Sprintf("/v2/instances/%d/%d/%s", projectID, regionID, instancesCheckLimitsPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
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

	instanceFlavorUpdateRequest := &InstanceFlavorUpdateRequest{
		FlavorID: "g1-standard-1-2",
	}
	taskResponse := &TaskResponse{Tasks: []string{taskID}}
	URL := fmt.Sprintf("/v1/instances/%d/%d/%s/%s", projectID, regionID, testResourceID, instancesChangeFlavorPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Instances.UpdateFlavor(ctx, testResourceID, instanceFlavorUpdateRequest)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestInstances_AvailableFlavors(t *testing.T) {
	setup()
	defer teardown()

	const (
		flavor = "g1-standard-2-8"
	)

	instanceCheckFlavorVolumeRequest := &InstanceCheckFlavorVolumeRequest{
		Volumes: []InstanceVolumeCreate{{Source: ExistingVolume}},
	}
	flavors := []Flavor{{FlavorID: flavor}}
	URL := fmt.Sprintf("/v1/instances/%d/%d/%s", projectID, regionID, instancesAvailableFlavorsPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, _ := json.Marshal(flavors)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.Instances.AvailableFlavors(ctx, instanceCheckFlavorVolumeRequest, nil)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, flavors) {
		t.Errorf("Instances.AvailableFlavors\n returned %+v,\n expected %+v", resp, flavors)
	}
}
