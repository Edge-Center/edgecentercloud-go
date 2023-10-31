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

	volumes := []Volume{{ID: testResourceID}}
	URL := fmt.Sprintf("/v1/volumes/%d/%d", projectID, regionID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
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

	volume := &Volume{ID: testResourceID}
	URL := fmt.Sprintf("/v1/volumes/%d/%d/%s", projectID, regionID, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(volume)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Volumes.Get(ctx, testResourceID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, volume) {
		t.Errorf("Volumes.Get\n returned %+v,\n expected %+v", resp, volume)
	}
}

func TestVolumes_Create(t *testing.T) {
	setup()
	defer teardown()

	volumeCreateRequest := &VolumeCreateRequest{
		Name:     "test-volume",
		Size:     20,
		TypeName: Standard,
		Source:   NewVolume,
	}
	taskResponse := &TaskResponse{Tasks: []string{testResourceID}}
	URL := fmt.Sprintf("/v1/volumes/%d/%d", projectID, regionID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
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

	taskResponse := &TaskResponse{Tasks: []string{taskID}}
	URL := fmt.Sprintf("/v1/volumes/%d/%d/%s", projectID, regionID, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Volumes.Delete(ctx, testResourceID)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestVolumes_ChangeType(t *testing.T) {
	setup()
	defer teardown()

	changeTypeRequest := &VolumeChangeTypeRequest{
		VolumeType: SsdHiIops,
	}
	volumeResponse := &Volume{ID: testResourceID, VolumeType: SsdHiIops}
	URL := fmt.Sprintf("/v1/volumes/%d/%d/%s/%s", projectID, regionID, testResourceID, volumesRetypePath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(VolumeChangeTypeRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, changeTypeRequest, reqBody)
		resp, _ := json.Marshal(volumeResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Volumes.ChangeType(ctx, testResourceID, changeTypeRequest)
	require.NoError(t, err)

	assert.Equal(t, volumeResponse, resp)
}

func TestVolumes_ExtendSize(t *testing.T) {
	setup()
	defer teardown()

	extendSizeRequest := &VolumeExtendSizeRequest{
		Size: 20,
	}
	taskResponse := &TaskResponse{Tasks: []string{taskID}}
	URL := fmt.Sprintf("/v1/volumes/%d/%d/%s/%s", projectID, regionID, testResourceID, volumesExtendPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(VolumeExtendSizeRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, extendSizeRequest, reqBody)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Volumes.Extend(ctx, testResourceID, extendSizeRequest)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestVolumes_Rename(t *testing.T) {
	setup()
	defer teardown()

	const (
		name = "new-name"
	)

	volumeRenameRequest := &VolumeRenameRequest{
		Name: name,
	}
	volumeResponse := &Volume{ID: testResourceID, Name: name}
	URL := fmt.Sprintf("/v1/volumes/%d/%d/%s", projectID, regionID, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		reqBody := new(VolumeRenameRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, volumeRenameRequest, reqBody)
		resp, _ := json.Marshal(volumeResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Volumes.Rename(ctx, testResourceID, volumeRenameRequest)
	require.NoError(t, err)

	assert.Equal(t, volumeResponse, resp)
}

func TestVolumes_Attach(t *testing.T) {
	setup()
	defer teardown()

	volumeAttachRequest := &VolumeAttachRequest{
		InstanceID: testResourceID,
	}
	volumeResponse := &Volume{ID: testResourceID}
	URL := fmt.Sprintf("/v1/volumes/%d/%d/%s/%s", projectID, regionID, testResourceID, volumesAttachPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(VolumeAttachRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, volumeAttachRequest, reqBody)
		resp, _ := json.Marshal(volumeResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Volumes.Attach(ctx, testResourceID, volumeAttachRequest)
	require.NoError(t, err)

	assert.Equal(t, volumeResponse, resp)
}

func TestVolumes_Detach(t *testing.T) {
	setup()
	defer teardown()

	volumeDetachRequest := &VolumeDetachRequest{
		InstanceID: testResourceID,
	}
	volumeResponse := &Volume{ID: testResourceID}
	URL := fmt.Sprintf("/v1/volumes/%d/%d/%s/%s", projectID, regionID, testResourceID, volumesDetachPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(VolumeDetachRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, volumeDetachRequest, reqBody)
		resp, _ := json.Marshal(volumeResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Volumes.Detach(ctx, testResourceID, volumeDetachRequest)
	require.NoError(t, err)

	assert.Equal(t, volumeResponse, resp)
}

func TestVolumes_Revert(t *testing.T) {
	setup()
	defer teardown()

	taskResponse := &TaskResponse{Tasks: []string{taskID}}
	URL := fmt.Sprintf("/v1/volumes/%d/%d/%s/%s", projectID, regionID, testResourceID, volumesRevertPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Volumes.Revert(ctx, testResourceID)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestVolumes_MetadataList(t *testing.T) {
	setup()
	defer teardown()

	metadataList := []MetadataDetailed{{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}}
	URL := fmt.Sprintf("/v1/volumes/%d/%d/%s/%s", projectID, regionID, testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(metadataList)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.Volumes.MetadataList(ctx, testResourceID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, metadataList) {
		t.Errorf("Volumes.MetadataList\n returned %+v,\n expected %+v", resp, metadataList)
	}
}

func TestVolumes_MetadataCreate(t *testing.T) {
	setup()
	defer teardown()

	metadataCreateRequest := &MetadataCreateRequest{
		map[string]interface{}{"key": "value"},
	}
	URL := fmt.Sprintf("/v1/volumes/%d/%d/%s/%s", projectID, regionID, testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
	})

	_, err := client.Volumes.MetadataCreate(ctx, testResourceID, metadataCreateRequest)
	require.NoError(t, err)
}

func TestVolumes_MetadataUpdate(t *testing.T) {
	setup()
	defer teardown()

	metadataCreateRequest := &MetadataCreateRequest{
		map[string]interface{}{"key": "value"},
	}
	URL := fmt.Sprintf("/v1/volumes/%d/%d/%s/%s", projectID, regionID, testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
	})

	_, err := client.Volumes.MetadataUpdate(ctx, testResourceID, metadataCreateRequest)
	require.NoError(t, err)
}

func TestVolumes_MetadataDeleteItem(t *testing.T) {
	setup()
	defer teardown()

	URL := fmt.Sprintf("/v1/volumes/%d/%d/%s/%s", projectID, regionID, testResourceID, metadataItemPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.Volumes.MetadataDeleteItem(ctx, testResourceID, nil)
	require.NoError(t, err)
}

func TestVolumes_MetadataGetItem(t *testing.T) {
	setup()
	defer teardown()

	metadata := &MetadataDetailed{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}
	URL := fmt.Sprintf("/v1/volumes/%d/%d/%s/%s", projectID, regionID, testResourceID, metadataItemPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(metadata)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Volumes.MetadataGetItem(ctx, testResourceID, nil)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, metadata) {
		t.Errorf("Volumes.MetadataGetItem\n returned %+v,\n expected %+v", resp, metadata)
	}
}
