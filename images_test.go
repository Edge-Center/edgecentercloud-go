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

func TestImages_List(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []Image{{ID: testResourceID}}
	URL := path.Join(imagesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Images.List(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestImages_List_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(imagesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Images.List(ctx, nil)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestImages_Create(t *testing.T) {
	setup()
	defer teardown()

	request := &ImageCreateRequest{Name: "test-image"}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(imagesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(ImageCreateRequest)
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

	respActual, resp, err := client.Images.Create(ctx, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestImages_Create_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &ImageCreateRequest{}
	URL := path.Join(imagesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Images.Create(ctx, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestImages_Create_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Images.Create(ctx, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestImages_Get(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &Image{ID: testResourceID}
	URL := path.Join(imagesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(expectedResp)
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Images.Get(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestImages_Get_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(imagesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Images.Get(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestImages_Delete(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(imagesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Images.Delete(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestImages_Delete_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(imagesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Images.Delete(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestImages_Update(t *testing.T) {
	setup()
	defer teardown()

	request := &ImageUpdateRequest{Name: "test-image"}
	expectedResp := &Image{ID: testResourceID}
	URL := path.Join(imagesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		reqBody := new(ImageUpdateRequest)
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

	respActual, resp, err := client.Images.Update(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestImages_Update_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(imagesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Images.Update(ctx, testResourceID, nil)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestImages_Upload(t *testing.T) {
	setup()
	defer teardown()

	request := &ImageUploadRequest{Name: "test-image"}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(downloadimageBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(ImageUploadRequest)
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

	respActual, resp, err := client.Images.Upload(ctx, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestImages_Upload_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(downloadimageBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Images.Upload(ctx, nil)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestImages_ImagesBaremetalList(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []Image{{ID: testResourceID}}
	URL := path.Join(bmimagesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Images.ImagesBaremetalList(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestImages_ImagesBaremetalList_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(bmimagesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Images.ImagesBaremetalList(ctx, nil)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestImages_ImagesBaremetalCreate(t *testing.T) {
	setup()
	defer teardown()

	request := &ImageCreateRequest{Name: "test-image"}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(bmimagesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(ImageCreateRequest)
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

	respActual, resp, err := client.Images.ImagesBaremetalCreate(ctx, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestImages_ImagesBaremetalCreate_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &ImageCreateRequest{}
	URL := path.Join(bmimagesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Images.ImagesBaremetalCreate(ctx, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestImages_ImagesBaremetalCreate_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Images.ImagesBaremetalCreate(ctx, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestImages_ImagesProjectList(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []Image{{ID: testResourceID}}
	URL := path.Join(projectimagesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Images.ImagesProjectList(ctx)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestImages_ImagesProjectList_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(projectimagesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Images.ImagesProjectList(ctx)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestImages_MetadataList(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []MetadataDetailed{{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}}
	URL := path.Join(imagesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Images.MetadataList(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestImages_MetadataCreate(t *testing.T) {
	setup()
	defer teardown()

	request := &Metadata{"key": "value"}
	URL := path.Join(imagesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
	})

	resp, err := client.Images.MetadataCreate(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestImages_MetadataUpdate(t *testing.T) {
	setup()
	defer teardown()

	request := &Metadata{"key": "value"}
	URL := path.Join(imagesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
	})

	resp, err := client.Images.MetadataUpdate(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestImages_MetadataDeleteItem(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(imagesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataItemPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	resp, err := client.Images.MetadataDeleteItem(ctx, testResourceID, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestImages_MetadataGetItem(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &MetadataDetailed{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}
	URL := path.Join(imagesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataItemPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Images.MetadataGetItem(ctx, testResourceID, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestImages_isValidUUID_Error_Return_Image(t *testing.T) {
	setup()
	defer teardown()

	tests := []struct {
		name     string
		testFunc func() (*Image, *Response, error)
	}{
		{
			name: "Get",
			testFunc: func() (*Image, *Response, error) {
				return client.Images.Get(ctx, testResourceIDNotValidUUID)
			},
		},
		{
			name: "Update",
			testFunc: func() (*Image, *Response, error) {
				return client.Images.Update(ctx, testResourceIDNotValidUUID, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respActual, resp, err := tt.testFunc()
			require.Nil(t, respActual)
			require.Equal(t, 400, resp.StatusCode)
			require.EqualError(t, err, NewArgError("imageID", NotCorrectUUID).Error())
		})
	}
}

func TestImages_Delete_isValidUUID_Error(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Images.Delete(ctx, testResourceIDNotValidUUID)
	require.Nil(t, respActual)
	require.Equal(t, 400, resp.StatusCode)
	require.EqualError(t, err, NewArgError("imageID", NotCorrectUUID).Error())
}

func TestImages_Metadata_isValidUUID_Error_Return_Response(t *testing.T) {
	setup()
	defer teardown()

	tests := []struct {
		name     string
		testFunc func() (*Response, error)
	}{
		{
			name: "MetadataCreate",
			testFunc: func() (*Response, error) {
				return client.Images.MetadataCreate(ctx, testResourceIDNotValidUUID, nil)
			},
		},
		{
			name: "MetadataUpdate",
			testFunc: func() (*Response, error) {
				return client.Images.MetadataUpdate(ctx, testResourceIDNotValidUUID, nil)
			},
		},
		{
			name: "MetadataDeleteItem",
			testFunc: func() (*Response, error) {
				return client.Images.MetadataDeleteItem(ctx, testResourceIDNotValidUUID, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.testFunc()
			require.Equal(t, 400, resp.StatusCode)
			require.EqualError(t, err, NewArgError("imageID", NotCorrectUUID).Error())
		})
	}
}

func TestImages_MetadataList_isValidUUID_Error(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Images.MetadataList(ctx, testResourceIDNotValidUUID)
	require.Nil(t, respActual)
	require.Equal(t, 400, resp.StatusCode)
	require.EqualError(t, err, NewArgError("imageID", NotCorrectUUID).Error())
}

func TestImages_MetadataGetItem_isValidUUID_Error(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Images.MetadataGetItem(ctx, testResourceIDNotValidUUID, nil)
	require.Nil(t, respActual)
	require.Equal(t, 400, resp.StatusCode)
	require.EqualError(t, err, NewArgError("imageID", NotCorrectUUID).Error())
}
