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

func TestSubnetworks_List(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []Subnetwork{{ID: testResourceID}}
	URL := path.Join(subnetsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Subnetworks.List(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestSubnetworks_List_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(subnetsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))
	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Subnetworks.List(ctx, nil)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestSubnetworks_Get(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &Subnetwork{ID: testResourceID}
	URL := path.Join(subnetsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Subnetworks.Get(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestSubnetworks_Get_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(subnetsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)
	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Subnetworks.Get(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestSubnetworks_Create(t *testing.T) {
	setup()
	defer teardown()

	request := &SubnetworkCreateRequest{
		Name:      "test-subnet",
		NetworkID: testResourceID,
	}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(subnetsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(SubnetworkCreateRequest)
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

	respActual, resp, err := client.Subnetworks.Create(ctx, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestSubnetworks_Create_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &SubnetworkCreateRequest{}
	URL := path.Join(subnetsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Subnetworks.Create(ctx, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestSubnetworks_Create_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Subnetworks.Create(ctx, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestSubnetworks_Delete(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(subnetsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Subnetworks.Delete(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestSubnetworks_Delete_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(subnetsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)
	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Subnetworks.Delete(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestSubnetworks_Update(t *testing.T) {
	setup()
	defer teardown()

	request := &SubnetworkUpdateRequest{
		Name: "test-subnet",
	}
	expectedResp := &Subnetwork{ID: testResourceID}
	URL := path.Join(subnetsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		reqBody := new(SubnetworkUpdateRequest)
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

	respActual, resp, err := client.Subnetworks.Update(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestSubnetworks_Update_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &SubnetworkUpdateRequest{}
	URL := path.Join(subnetsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Subnetworks.Update(ctx, testResourceID, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestSubnetworks_Update_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Subnetworks.Update(ctx, testResourceID, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestSubnetworks_MetadataList(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []MetadataDetailed{{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}}
	URL := path.Join(subnetsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Subnetworks.MetadataList(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestSubnetworks_MetadataCreate(t *testing.T) {
	setup()
	defer teardown()

	request := &Metadata{"key": "value"}
	URL := path.Join(subnetsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
	})

	resp, err := client.Subnetworks.MetadataCreate(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestSubnetworks_MetadataUpdate(t *testing.T) {
	setup()
	defer teardown()

	request := &Metadata{"key": "value"}
	URL := path.Join(subnetsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
	})

	resp, err := client.Subnetworks.MetadataUpdate(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestSubnetworks_MetadataDeleteItem(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(subnetsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataItemPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	resp, err := client.Subnetworks.MetadataDeleteItem(ctx, testResourceID, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestSubnetworks_MetadataGetItem(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &MetadataDetailed{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}
	URL := path.Join(subnetsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataItemPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Subnetworks.MetadataGetItem(ctx, testResourceID, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestSubnetworks_isValidUUID_Error_Return_Subnetwork(t *testing.T) {
	setup()
	defer teardown()

	tests := []struct {
		name     string
		testFunc func() (*Subnetwork, *Response, error)
	}{
		{
			name: "Get",
			testFunc: func() (*Subnetwork, *Response, error) {
				return client.Subnetworks.Get(ctx, testResourceIDNotValidUUID)
			},
		},
		{
			name: "Update",
			testFunc: func() (*Subnetwork, *Response, error) {
				return client.Subnetworks.Update(ctx, testResourceIDNotValidUUID, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respActual, resp, err := tt.testFunc()
			require.Nil(t, respActual)
			require.Equal(t, 400, resp.StatusCode)
			require.EqualError(t, err, NewArgError("subnetworkID", NotCorrectUUID).Error())
		})
	}
}

func TestSubnetworks_Delete_isValidUUID_Error(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Subnetworks.Delete(ctx, testResourceIDNotValidUUID)
	require.Nil(t, respActual)
	require.Equal(t, 400, resp.StatusCode)
	require.EqualError(t, err, NewArgError("subnetworkID", NotCorrectUUID).Error())
}

func TestSubnetworks_Metadata_isValidUUID_Error_Return_Response(t *testing.T) {
	setup()
	defer teardown()

	tests := []struct {
		name     string
		testFunc func() (*Response, error)
	}{
		{
			name: "MetadataCreate",
			testFunc: func() (*Response, error) {
				return client.Subnetworks.MetadataCreate(ctx, testResourceIDNotValidUUID, nil)
			},
		},
		{
			name: "MetadataUpdate",
			testFunc: func() (*Response, error) {
				return client.Subnetworks.MetadataUpdate(ctx, testResourceIDNotValidUUID, nil)
			},
		},
		{
			name: "MetadataDeleteItem",
			testFunc: func() (*Response, error) {
				return client.Subnetworks.MetadataDeleteItem(ctx, testResourceIDNotValidUUID, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.testFunc()
			require.Equal(t, 400, resp.StatusCode)
			require.EqualError(t, err, NewArgError("subnetworkID", NotCorrectUUID).Error())
		})
	}
}

func TestSubnetworks_MetadataList_isValidUUID_Error(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Subnetworks.MetadataList(ctx, testResourceIDNotValidUUID)
	require.Nil(t, respActual)
	require.Equal(t, 400, resp.StatusCode)
	require.EqualError(t, err, NewArgError("subnetworkID", NotCorrectUUID).Error())
}

func TestSubnetworks_MetadataGetItem_isValidUUID_Error(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Subnetworks.MetadataGetItem(ctx, testResourceIDNotValidUUID, nil)
	require.Nil(t, respActual)
	require.Equal(t, 400, resp.StatusCode)
	require.EqualError(t, err, NewArgError("subnetworkID", NotCorrectUUID).Error())
}
