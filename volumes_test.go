package edgecloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVolumes_List(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []Volume{{ID: testResourceID}}
	URL := path.Join(volumesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Volumes.List(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestVolumes_List_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(volumesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Volumes.List(ctx, nil)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestVolumes_Get(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &Volume{ID: testResourceID}
	URL := path.Join(volumesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Volumes.Get(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestVolumes_Get_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(volumesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Volumes.Get(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestVolumes_Create(t *testing.T) {
	setup()
	defer teardown()

	request := &VolumeCreateRequest{
		Name:     "test-volume",
		Size:     20,
		TypeName: VolumeTypeStandard,
		Source:   VolumeSourceNewVolume,
	}
	expectedResp := &TaskResponse{Tasks: []string{testResourceID}}
	URL := path.Join(volumesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(VolumeCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, reqBody)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Volumes.Create(ctx, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestVolumes_Create_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &VolumeCreateRequest{}
	URL := path.Join(volumesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Volumes.Create(ctx, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestVolumes_Create_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Volumes.Create(ctx, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestVolumes_Delete(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(volumesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Volumes.Delete(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestVolumes_Delete_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(volumesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Volumes.Delete(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestVolumes_ChangeType(t *testing.T) {
	setup()
	defer teardown()

	request := &VolumeChangeTypeRequest{VolumeType: VolumeTypeSsdHiIops}
	expectedResp := &Volume{ID: testResourceID, VolumeType: VolumeTypeSsdHiIops}
	URL := path.Join(volumesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, volumesRetype)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(VolumeChangeTypeRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, reqBody)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Volumes.ChangeType(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestVolumes_ChangeType_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &VolumeChangeTypeRequest{}
	URL := path.Join(volumesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, volumesRetype)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Volumes.ChangeType(ctx, testResourceID, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestVolumes_ChangeType_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Volumes.ChangeType(ctx, testResourceID, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestVolumes_ExtendSize(t *testing.T) {
	setup()
	defer teardown()

	request := &VolumeExtendSizeRequest{Size: 20}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(volumesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, volumesExtend)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(VolumeExtendSizeRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, reqBody)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Volumes.Extend(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestVolumes_ExtendSize_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &VolumeExtendSizeRequest{}
	URL := path.Join(volumesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, volumesExtend)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Volumes.Extend(ctx, testResourceID, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestVolumes_ExtendSize_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Volumes.Extend(ctx, testResourceID, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestVolumes_Rename(t *testing.T) {
	setup()
	defer teardown()

	const (
		name = "new-name"
	)

	request := &Name{Name: name}
	expectedResp := &Volume{ID: testResourceID, Name: name}
	URL := path.Join(volumesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		reqBody := new(Name)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, reqBody)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Volumes.Rename(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestVolumes_Rename_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &Name{}
	URL := path.Join(volumesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Volumes.Rename(ctx, testResourceID, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestVolumes_Rename_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Volumes.Rename(ctx, testResourceID, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestVolumes_Attach(t *testing.T) {
	setup()
	defer teardown()

	request := &VolumeAttachRequest{InstanceID: testResourceID}
	expectedResp := &Volume{ID: testResourceID}
	URL := path.Join(volumesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, volumesAttach)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(VolumeAttachRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, reqBody)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Volumes.Attach(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestVolumes_Attach_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &VolumeAttachRequest{}
	URL := path.Join(volumesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, volumesAttach)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Volumes.Attach(ctx, testResourceID, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestVolumes_Attach_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Volumes.Attach(ctx, testResourceID, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestVolumes_Detach(t *testing.T) {
	setup()
	defer teardown()

	request := &VolumeDetachRequest{InstanceID: testResourceID}
	expectedResp := &Volume{ID: testResourceID}
	URL := path.Join(volumesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, volumesDetach)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(VolumeDetachRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, reqBody)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Volumes.Detach(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestVolumes_Detach_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &VolumeDetachRequest{}
	URL := path.Join(volumesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, volumesDetach)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Volumes.Detach(ctx, testResourceID, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestVolumes_Detach_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Volumes.Detach(ctx, testResourceID, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestVolumes_Revert(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(volumesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, volumesRevert)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Volumes.Revert(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestVolumes_Revert_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(volumesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, volumesRevert)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Volumes.Revert(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestVolumes_MetadataList(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []MetadataDetailed{{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}}
	URL := path.Join(volumesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Volumes.MetadataList(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestVolumes_MetadataCreate(t *testing.T) {
	setup()
	defer teardown()

	request := &Metadata{"key": "value"}
	URL := path.Join(volumesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
	})

	resp, err := client.Volumes.MetadataCreate(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestVolumes_MetadataUpdate(t *testing.T) {
	setup()
	defer teardown()

	request := &Metadata{"key": "value"}
	URL := path.Join(volumesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
	})

	resp, err := client.Volumes.MetadataUpdate(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestVolumes_MetadataDeleteItem(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(volumesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataItemPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	resp, err := client.Volumes.MetadataDeleteItem(ctx, testResourceID, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestVolumes_MetadataGetItem(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &MetadataDetailed{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}
	URL := path.Join(volumesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataItemPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Volumes.MetadataGetItem(ctx, testResourceID, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestVolumes_isValidUUID_Error_Return_Volume(t *testing.T) {
	setup()
	defer teardown()

	tests := []struct {
		name     string
		testFunc func() (*Volume, *Response, error)
	}{
		{
			name: "Get",
			testFunc: func() (*Volume, *Response, error) {
				return client.Volumes.Get(ctx, testResourceIDNotValidUUID)
			},
		},
		{
			name: "ChangeType",
			testFunc: func() (*Volume, *Response, error) {
				return client.Volumes.ChangeType(ctx, testResourceIDNotValidUUID, nil)
			},
		},
		{
			name: "Rename",
			testFunc: func() (*Volume, *Response, error) {
				return client.Volumes.Rename(ctx, testResourceIDNotValidUUID, nil)
			},
		},
		{
			name: "Attach",
			testFunc: func() (*Volume, *Response, error) {
				return client.Volumes.Attach(ctx, testResourceIDNotValidUUID, nil)
			},
		},
		{
			name: "Detach",
			testFunc: func() (*Volume, *Response, error) {
				return client.Volumes.Detach(ctx, testResourceIDNotValidUUID, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respActual, resp, err := tt.testFunc()
			require.Nil(t, respActual)
			require.Equal(t, 400, resp.StatusCode)
			require.EqualError(t, err, NewArgError("volumeID", NotCorrectUUID).Error())
		})
	}
}

func TestVolumes_isValidUUID_Error_Return_TaskResponse(t *testing.T) {
	setup()
	defer teardown()

	tests := []struct {
		name     string
		testFunc func() (*TaskResponse, *Response, error)
	}{
		{
			name: "Delete",
			testFunc: func() (*TaskResponse, *Response, error) {
				return client.Volumes.Delete(ctx, testResourceIDNotValidUUID)
			},
		},
		{
			name: "Revert",
			testFunc: func() (*TaskResponse, *Response, error) {
				return client.Volumes.Revert(ctx, testResourceIDNotValidUUID)
			},
		},
		{
			name: "Extend",
			testFunc: func() (*TaskResponse, *Response, error) {
				return client.Volumes.Extend(ctx, testResourceIDNotValidUUID, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respActual, resp, err := tt.testFunc()
			require.Nil(t, respActual)
			require.Equal(t, 400, resp.StatusCode)
			require.EqualError(t, err, NewArgError("volumeID", NotCorrectUUID).Error())
		})
	}
}

func TestVolumes_Metadata_isValidUUID_Error_Return_Response(t *testing.T) {
	setup()
	defer teardown()

	tests := []struct {
		name     string
		testFunc func() (*Response, error)
	}{
		{
			name: "MetadataCreate",
			testFunc: func() (*Response, error) {
				return client.Volumes.MetadataCreate(ctx, testResourceIDNotValidUUID, nil)
			},
		},
		{
			name: "MetadataUpdate",
			testFunc: func() (*Response, error) {
				return client.Volumes.MetadataUpdate(ctx, testResourceIDNotValidUUID, nil)
			},
		},
		{
			name: "MetadataDeleteItem",
			testFunc: func() (*Response, error) {
				return client.Volumes.MetadataDeleteItem(ctx, testResourceIDNotValidUUID, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.testFunc()
			require.Equal(t, 400, resp.StatusCode)
			require.EqualError(t, err, NewArgError("volumeID", NotCorrectUUID).Error())
		})
	}
}

func TestVolumes_MetadataList_isValidUUID_Error(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Volumes.MetadataList(ctx, testResourceIDNotValidUUID)
	require.Nil(t, respActual)
	require.Equal(t, 400, resp.StatusCode)
	require.EqualError(t, err, NewArgError("volumeID", NotCorrectUUID).Error())
}

func TestVolumes_MetadataGetItem_isValidUUID_Error(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Volumes.MetadataGetItem(ctx, testResourceIDNotValidUUID, nil)
	require.Nil(t, respActual)
	require.Equal(t, 400, resp.StatusCode)
	require.EqualError(t, err, NewArgError("volumeID", NotCorrectUUID).Error())
}
