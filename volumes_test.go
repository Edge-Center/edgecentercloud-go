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

func TestVolumes_ChangeType(t *testing.T) {
	setup()
	defer teardown()

	const (
		volumeID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	changeTypeRequest := &VolumeChangeTypeRequest{
		VolumeType: SsdHiIops,
	}

	volumeResponse := &Volume{ID: volumeID, VolumeType: SsdHiIops}

	changeTypeVolumeURL := fmt.Sprintf("/v1/volumes/%d/%d/%s/%s", projectID, regionID, volumeID, volumesRetypePath)

	mux.HandleFunc(changeTypeVolumeURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(VolumeChangeTypeRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, changeTypeRequest, reqBody)
		resp, _ := json.Marshal(volumeResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Volumes.ChangeType(ctx, volumeID, changeTypeRequest)
	require.NoError(t, err)

	assert.Equal(t, volumeResponse, resp)
}

func TestVolumes_ExtendSize(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID   = "f0d19cec-5c3f-4853-886e-304915960ff6"
		volumeID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	extendSizeRequest := &VolumeExtendSizeRequest{
		Size: 20,
	}

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	extendSizeVolumeURL := fmt.Sprintf("/v1/volumes/%d/%d/%s/%s", projectID, regionID, volumeID, volumesExtendPath)

	mux.HandleFunc(extendSizeVolumeURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(VolumeExtendSizeRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, extendSizeRequest, reqBody)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Volumes.Extend(ctx, volumeID, extendSizeRequest)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestVolumes_Rename(t *testing.T) {
	setup()
	defer teardown()

	const (
		volumeID = "f0d19cec-5c3f-4853-886e-304915960ff6"
		name     = "new-name"
	)

	volumeRenameRequest := &VolumeRenameRequest{
		Name: name,
	}

	volumeResponse := &Volume{ID: volumeID, Name: name}

	renameVolumeURL := fmt.Sprintf("/v1/volumes/%d/%d/%s", projectID, regionID, volumeID)

	mux.HandleFunc(renameVolumeURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		reqBody := new(VolumeRenameRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, volumeRenameRequest, reqBody)
		resp, _ := json.Marshal(volumeResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Volumes.Rename(ctx, volumeID, volumeRenameRequest)
	require.NoError(t, err)

	assert.Equal(t, volumeResponse, resp)
}

func TestVolumes_Attach(t *testing.T) {
	setup()
	defer teardown()

	const (
		volumeID   = "f0d19cec-5c3f-4853-886e-304915960ff6"
		instanceID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	volumeAttachRequest := &VolumeAttachRequest{
		InstanceID: instanceID,
	}

	volumeResponse := &Volume{ID: volumeID}

	attachVolumeURL := fmt.Sprintf("/v1/volumes/%d/%d/%s/%s", projectID, regionID, volumeID, volumesAttachPath)

	mux.HandleFunc(attachVolumeURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(VolumeAttachRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, volumeAttachRequest, reqBody)
		resp, _ := json.Marshal(volumeResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Volumes.Attach(ctx, volumeID, volumeAttachRequest)
	require.NoError(t, err)

	assert.Equal(t, volumeResponse, resp)
}

func TestVolumes_Detach(t *testing.T) {
	setup()
	defer teardown()

	const (
		volumeID   = "f0d19cec-5c3f-4853-886e-304915960ff6"
		instanceID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	volumeDetachRequest := &VolumeDetachRequest{
		InstanceID: instanceID,
	}

	volumeResponse := &Volume{ID: volumeID}

	detachVolumeURL := fmt.Sprintf("/v1/volumes/%d/%d/%s/%s", projectID, regionID, volumeID, volumesDetachPath)

	mux.HandleFunc(detachVolumeURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(VolumeDetachRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, volumeDetachRequest, reqBody)
		resp, _ := json.Marshal(volumeResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Volumes.Detach(ctx, volumeID, volumeDetachRequest)
	require.NoError(t, err)

	assert.Equal(t, volumeResponse, resp)
}

func TestVolumes_Revert(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID   = "f0d19cec-5c3f-4853-886e-304915960ff6"
		volumeID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	revertVolumeURL := fmt.Sprintf("/v1/volumes/%d/%d/%s/%s", projectID, regionID, volumeID, volumesRevertPath)

	mux.HandleFunc(revertVolumeURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Volumes.Revert(ctx, volumeID)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestVolumes_MetadataList(t *testing.T) {
	setup()
	defer teardown()

	const (
		volumeID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	metadataList := []MetadataDetailed{{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}}
	getVolumesMetadataListURL := fmt.Sprintf("/v1/volumes/%d/%d/%s/%s", projectID, regionID, volumeID, metadataPath)

	mux.HandleFunc(getVolumesMetadataListURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(metadataList)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.Volumes.MetadataList(ctx, volumeID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, metadataList) {
		t.Errorf("Volumes.MetadataList\n returned %+v,\n expected %+v", resp, metadataList)
	}
}

func TestVolumes_MetadataCreate(t *testing.T) {
	setup()
	defer teardown()

	const (
		volumeID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	metadataCreateRequest := &MetadataCreateRequest{
		map[string]interface{}{"key": "value"},
	}

	createVolumesMetadataURL := fmt.Sprintf("/v1/volumes/%d/%d/%s/%s", projectID, regionID, volumeID, metadataPath)

	mux.HandleFunc(createVolumesMetadataURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
	})

	_, err := client.Volumes.MetadataCreate(ctx, volumeID, metadataCreateRequest)
	require.NoError(t, err)
}

func TestVolumes_MetadataUpdate(t *testing.T) {
	setup()
	defer teardown()

	const (
		volumeID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	metadataCreateRequest := &MetadataCreateRequest{
		map[string]interface{}{"key": "value"},
	}

	createVolumesMetadataURL := fmt.Sprintf("/v1/volumes/%d/%d/%s/%s", projectID, regionID, volumeID, metadataPath)

	mux.HandleFunc(createVolumesMetadataURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
	})

	_, err := client.Volumes.MetadataUpdate(ctx, volumeID, metadataCreateRequest)
	require.NoError(t, err)
}

func TestVolumes_MetadataDeleteItem(t *testing.T) {
	setup()
	defer teardown()

	const (
		volumeID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	deleteVolumesMetadataItemURL := fmt.Sprintf("/v1/volumes/%d/%d/%s/%s", projectID, regionID, volumeID, metadataItemPath)

	mux.HandleFunc(deleteVolumesMetadataItemURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.Volumes.MetadataDeleteItem(ctx, volumeID, nil)
	require.NoError(t, err)
}

func TestVolumes_MetadataGetItem(t *testing.T) {
	setup()
	defer teardown()

	const (
		volumeID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	metadata := &MetadataDetailed{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}
	getVolumesMetadataItemURL := fmt.Sprintf("/v1/volumes/%d/%d/%s/%s", projectID, regionID, volumeID, metadataItemPath)

	mux.HandleFunc(getVolumesMetadataItemURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(metadata)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Volumes.MetadataGetItem(ctx, volumeID, nil)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, metadata) {
		t.Errorf("Volumes.MetadataGetItem\n returned %+v,\n expected %+v", resp, metadata)
	}
}
