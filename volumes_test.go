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

func TestVolumes_List(t *testing.T) {
	setup()
	defer teardown()

	const (
		volumeID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	volumes := []Volume{{ID: volumeID}}
	getVolumesURL := fmt.Sprintf("/v1/volumes/%d/%d", projectID, regionID)

	mux.HandleFunc(getVolumesURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(volumes)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.Volumes.List(ctx, nil)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, volumes) {
		t.Errorf("Volumes.List\n returned %+v,\n expected %+v", resp, volumes)
	}
}

func TestVolumes_Get(t *testing.T) {
	setup()
	defer teardown()

	const (
		volumeID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	volume := &Volume{ID: volumeID}
	getVolumesURL := fmt.Sprintf("/v1/volumes/%d/%d/%s", projectID, regionID, volumeID)

	mux.HandleFunc(getVolumesURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(volume)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Volumes.Get(ctx, volumeID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, volume) {
		t.Errorf("Volumes.Get\n returned %+v,\n expected %+v", resp, volume)
	}
}

func TestVolumes_Create(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	volumeCreateRequest := &VolumeCreateRequest{
		Name:     "test-volume",
		Size:     20,
		TypeName: Standard,
		Source:   NewVolume,
	}

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	createVolumeURL := fmt.Sprintf("/v1/volumes/%d/%d", projectID, regionID)

	mux.HandleFunc(createVolumeURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(VolumeCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, volumeCreateRequest, reqBody)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Volumes.Create(ctx, volumeCreateRequest)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestVolumes_Delete(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID   = "f0d19cec-5c3f-4853-886e-304915960ff6"
		volumeID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	deleteVolumeURL := fmt.Sprintf("/v1/volumes/%d/%d/%s", projectID, regionID, volumeID)

	mux.HandleFunc(deleteVolumeURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Volumes.Delete(ctx, volumeID)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}
