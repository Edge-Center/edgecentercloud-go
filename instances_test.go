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

func TestInstances_List(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []Instance{{ID: testResourceID}}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Instances.List(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestInstances_List_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Instances.List(ctx, nil)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestInstances_Get(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &Instance{
		ID: testResourceID,
	}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(expectedResp)
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Instances.Get(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestInstances_Get_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Instances.Get(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestInstances_Create(t *testing.T) {
	setup()
	defer teardown()

	bootIndex := 0

	request := &InstanceCreateRequest{
		Names:          []string{"test-instance"},
		Flavor:         "g1-standard-1-2",
		Interfaces:     []InstanceInterface{{Type: InterfaceTypeExternal}},
		SecurityGroups: []ID{{ID: "f0d19cec-5c3f-4853-886e-304915960ff6"}},
		Volumes: []InstanceVolumeCreate{
			{
				TypeName:  VolumeTypeSsdHiIops,
				Size:      5,
				BootIndex: &bootIndex,
				Source:    VolumeSourceImage,
				ImageID:   "f0d19cec-5c3f-4853-886e-304915960ff6",
			},
		},
	}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(instancesBasePathV2, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(InstanceCreateRequest)
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

	respActual, resp, err := client.Instances.Create(ctx, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestInstances_Create_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Instances.Create(ctx, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestInstances_Create_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(instancesBasePathV2, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	request := &InstanceCreateRequest{}
	respActual, resp, err := client.Instances.Create(ctx, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestInstances_Delete(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Instances.Delete(ctx, testResourceID, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestInstances_Delete_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Instances.Delete(ctx, testResourceID, nil)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestInstances_DeleteWithOptions(t *testing.T) {
	setup()
	defer teardown()

	instanceDeleteOptions := InstanceDeleteOptions{
		DeleteFloatings: true,
		Volumes:         []string{"f0d19cec-5c3f-4853-886e-304915960ff6"},
	}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Instances.Delete(ctx, testResourceID, &instanceDeleteOptions)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestInstances_MetadataGet(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &MetadataDetailed{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Instances.MetadataGet(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestInstances_MetadataGet_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Instances.MetadataGet(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestInstances_MetadataList(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []MetadataDetailed{{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}}

	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Instances.MetadataList(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestInstances_MetadataUpdate(t *testing.T) {
	setup()
	defer teardown()

	request := &Metadata{"key": "value"}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
	})

	resp, err := client.Instances.MetadataUpdate(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestInstances_MetadataDeleteItem(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(instancesBasePathV2, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataItemPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	resp, err := client.Instances.MetadataDeleteItem(ctx, testResourceID, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestInstances_MetadataGetItem(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &MetadataDetailed{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}
	URL := path.Join(instancesBasePathV2, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataItemPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Instances.MetadataGetItem(ctx, testResourceID, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestInstances_MetadataCreate(t *testing.T) {
	setup()
	defer teardown()

	request := &Metadata{"key": "value"}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
	})

	resp, err := client.Instances.MetadataCreate(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestInstances_CheckLimits(t *testing.T) {
	setup()
	defer teardown()

	request := &InstanceCheckLimitsRequest{
		Names:      []string{""},
		Interfaces: []InstanceInterface{},
		Volumes:    []InstanceCheckLimitsVolume{},
	}
	URL := path.Join(instancesBasePathV2, strconv.Itoa(projectID), strconv.Itoa(regionID), instancesCheckLimits)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, err := json.Marshal(map[string]int{})
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	_, resp, err := client.Instances.CheckLimits(ctx, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestInstances_CheckLimits_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(instancesBasePathV2, strconv.Itoa(projectID), strconv.Itoa(regionID), instancesCheckLimits)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Instances.CheckLimits(ctx, nil)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestInstances_UpdateFlavor(t *testing.T) {
	setup()
	defer teardown()

	request := &InstanceFlavorUpdateRequest{FlavorID: "g1-standard-1-2"}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesChangeFlavor)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Instances.UpdateFlavor(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, expectedResp, respActual)
}

func TestInstances_UpdateFlavor_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &InstanceFlavorUpdateRequest{}

	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesChangeFlavor)
	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Instances.UpdateFlavor(ctx, testResourceID, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestInstances_UpdateFlavor_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Instances.UpdateFlavor(ctx, testResourceID, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestInstances_AvailableFlavors(t *testing.T) {
	setup()
	defer teardown()

	request := &InstanceCheckFlavorVolumeRequest{
		Volumes: []InstanceVolumeCreate{{Source: VolumeSourceExistingVolume}},
	}
	expectedResp := []Flavor{{FlavorID: "g1-standard-2-8"}}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), instancesAvailableFlavors)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Instances.AvailableFlavors(ctx, request, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, expectedResp, respActual)
}

func TestInstances_AvailableFlavors_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &InstanceCheckFlavorVolumeRequest{}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), instancesAvailableFlavors)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Instances.AvailableFlavors(ctx, request, nil)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestInstances_AvailableFlavors_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Instances.AvailableFlavors(ctx, nil, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestInstances_AvailableFlavorsToResize(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []Flavor{{FlavorID: "g1-standard-2-8"}}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesAvailableFlavors)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Instances.AvailableFlavorsToResize(ctx, testResourceID, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, expectedResp, respActual)
}

func TestInstances_AvailableFlavorsToResize_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesAvailableFlavors)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Instances.AvailableFlavorsToResize(ctx, testResourceID, nil)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestInstances_AvailableNames(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &InstanceAvailableNames{}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), instancesAvailableNames)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Instances.AvailableNames(ctx)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, expectedResp, respActual)
}

func TestInstances_AvailableNames_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), instancesAvailableNames)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Instances.AvailableNames(ctx)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestInstances_Rename(t *testing.T) {
	setup()
	defer teardown()

	const (
		name = "new-name"
	)

	request := &Name{Name: name}
	expectedResp := &Instance{ID: testResourceID, Name: name}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

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

	respActual, resp, err := client.Instances.Rename(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestInstances_Rename_ResponseError(t *testing.T) {
	setup()
	defer teardown()
	request := &Name{}

	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Instances.Rename(ctx, testResourceID, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestInstances_Rename_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Instances.Rename(ctx, testResourceID, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestInstances_PortsList(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []InstancePort{{ID: testResourceID}}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesPorts)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Instances.PortsList(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestInstances_PortsList_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesPorts)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Instances.PortsList(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestInstances_InstanceStart(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &Instance{ID: testResourceID}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesStart)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Instances.InstanceStart(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestInstances_InstanceStart_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesStart)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Instances.InstanceStart(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestInstances_InstanceStop(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &Instance{ID: testResourceID}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesStop)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Instances.InstanceStop(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestInstances_InstanceStop_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesStop)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Instances.InstanceStop(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestInstances_InstancePowercycle(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &Instance{ID: testResourceID}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesPowercycle)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Instances.InstancePowercycle(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestInstances_InstancePowercycle_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesPowercycle)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Instances.InstancePowercycle(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestInstances_InstanceReboot(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &Instance{ID: testResourceID}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesReboot)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Instances.InstanceReboot(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestInstances_InstanceReboot_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesReboot)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Instances.InstanceReboot(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestInstances_InstanceSuspend(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &Instance{ID: testResourceID}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesSuspend)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Instances.InstanceSuspend(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestInstances_InstanceSuspend_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesSuspend)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Instances.InstanceSuspend(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestInstances_InstanceResume(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &Instance{ID: testResourceID}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesResume)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Instances.InstanceResume(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestInstances_InstanceResume_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesResume)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Instances.InstanceResume(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestInstances_MetricsList(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []InstanceMetrics{{}}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesMetrics)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Instances.MetricsList(ctx, testResourceID, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestInstances_MetricsList_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesMetrics)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Instances.MetricsList(ctx, testResourceID, nil)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestInstances_FilterBySecurityGroup(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []Instance{{}}
	URL := path.Join(securitygroupsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesInstances)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Instances.FilterBySecurityGroup(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestInstances_FilterBySecurityGroup_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesInstances)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Instances.FilterBySecurityGroup(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestInstances_SecurityGroupList(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []IDName{{}}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesSecurityGroups)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Instances.SecurityGroupList(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestInstances_SecurityGroupList_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesSecurityGroups)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Instances.SecurityGroupList(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestInstances_Assign(t *testing.T) {
	setup()
	defer teardown()

	request := &AssignSecurityGroupRequest{Name: "test-name"}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesAddSecurityGroup)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(AssignSecurityGroupRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, reqBody)
	})

	resp, err := client.Instances.SecurityGroupAssign(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestInstances_Assign_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	resp, err := client.Instances.SecurityGroupAssign(ctx, testResourceID, nil)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestInstances_UnAssign(t *testing.T) {
	setup()
	defer teardown()

	request := &AssignSecurityGroupRequest{Name: "test-name"}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesDelSecurityGroup)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(AssignSecurityGroupRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, reqBody)
	})

	resp, err := client.Instances.SecurityGroupUnAssign(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestInstances_UnAssign_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	resp, err := client.Instances.SecurityGroupUnAssign(ctx, testResourceID, nil)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestInstances_GetConsole(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &RemoteConsole{}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesGetConsole)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Instances.GetConsole(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestInstances_GetConsole_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesGetConsole)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Instances.GetConsole(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestInstances_AttachInterface(t *testing.T) {
	setup()
	defer teardown()

	request := &InstanceAttachInterfaceRequest{}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesAttachInterface)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(InstanceAttachInterfaceRequest)
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

	respActual, resp, err := client.Instances.AttachInterface(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, expectedResp, respActual)
}

func TestInstances_AttachInterface_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &InstanceAttachInterfaceRequest{}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesAttachInterface)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Instances.AttachInterface(ctx, testResourceID, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestInstances_AttachInterface_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Instances.AttachInterface(ctx, testResourceID, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestInstances_DetachInterface(t *testing.T) {
	setup()
	defer teardown()

	request := &InstanceDetachInterfaceRequest{}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesDetachInterface)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(InstanceDetachInterfaceRequest)
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

	respActual, resp, err := client.Instances.DetachInterface(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, expectedResp, respActual)
}

func TestInstances_DetachInterface_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &InstanceDetachInterfaceRequest{}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesDetachInterface)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Instances.DetachInterface(ctx, testResourceID, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestInstances_DetachInterface_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Instances.DetachInterface(ctx, testResourceID, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestInstances_InterfaceList(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []InstancePortInterface{{}}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesInterfaces)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Instances.InterfaceList(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestInstances_InterfaceList_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesInterfaces)
	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Instances.InterfaceList(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestInstances_PutIntoServerGroup(t *testing.T) {
	setup()
	defer teardown()

	request := &InstancePutIntoServerGroupRequest{}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesPutIntoServerGroup)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(InstancePutIntoServerGroupRequest)
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

	respActual, resp, err := client.Instances.PutIntoServerGroup(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, expectedResp, respActual)
}

func TestInstances_PutIntoServerGroup_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &InstancePutIntoServerGroupRequest{}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesPutIntoServerGroup)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Instances.PutIntoServerGroup(ctx, testResourceID, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestInstances_PutIntoServerGroup_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Instances.PutIntoServerGroup(ctx, testResourceID, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestInstances_RemoveFromServerGroup(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesRemoveFromServerGroup)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Instances.RemoveFromServerGroup(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, expectedResp, respActual)
}

func TestInstances_RemoveFromServerGroup_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(instancesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, instancesRemoveFromServerGroup)
	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})
	respActual, resp, err := client.Instances.RemoveFromServerGroup(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}
