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

func TestFloatingips_List(t *testing.T) {
	setup()
	defer teardown()

	floatingIPs := []FloatingIP{{ID: testResourceID}}
	URL := fmt.Sprintf("/v1/floatingips/%d/%d", projectID, regionID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(&floatingIPs)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.Floatingips.List(ctx)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, floatingIPs) {
		t.Errorf("Floatingips.List\n returned %+v,\n expected %+v", resp, floatingIPs)
	}
}

func TestFloatingips_Get(t *testing.T) {
	setup()
	defer teardown()

	floatingIP := &FloatingIP{ID: testResourceID}
	URL := fmt.Sprintf("/v1/floatingips/%d/%d/%s", projectID, regionID, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(floatingIP)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Floatingips.Get(ctx, testResourceID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, floatingIP) {
		t.Errorf("Floatingips.Get\n returned %+v,\n expected %+v", resp, floatingIP)
	}
}

func TestFloatingips_Create(t *testing.T) {
	setup()
	defer teardown()

	floatingIPCreateRequest := &FloatingIPCreateRequest{}
	taskResponse := &TaskResponse{Tasks: []string{taskID}}
	URL := fmt.Sprintf("/v1/floatingips/%d/%d", projectID, regionID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(FloatingIPCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, floatingIPCreateRequest, reqBody)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Floatingips.Create(ctx, floatingIPCreateRequest)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestFloatingips_Delete(t *testing.T) {
	setup()
	defer teardown()

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	URL := fmt.Sprintf("/v1/floatingips/%d/%d/%s", projectID, regionID, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Floatingips.Delete(ctx, testResourceID)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestFloatingips_Assign(t *testing.T) {
	setup()
	defer teardown()

	assignRequest := &AssignRequest{PortID: testResourceID}

	floatingIP := &FloatingIP{
		ID:     testResourceID,
		PortID: testResourceID,
	}
	URL := fmt.Sprintf("/v1/floatingips/%d/%d/%s/%s", projectID, regionID, testResourceID, floatingipsAssign)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, _ := json.Marshal(&floatingIP)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Floatingips.Assign(ctx, testResourceID, assignRequest)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, floatingIP) {
		t.Errorf("Floatingips.Assign\n returned %+v,\n expected %+v", resp, floatingIP)
	}
}

func TestFloatingips_UnAssign(t *testing.T) {
	setup()
	defer teardown()

	floatingIP := &FloatingIP{ID: testResourceID}
	URL := fmt.Sprintf("/v1/floatingips/%d/%d/%s/%s", projectID, regionID, testResourceID, floatingipsUnAssign)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, _ := json.Marshal(&floatingIP)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Floatingips.UnAssign(ctx, testResourceID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, floatingIP) {
		t.Errorf("Floatingips.UnAssign\n returned %+v,\n expected %+v", resp, floatingIP)
	}
}

func TestFloatingips_ListAvailable(t *testing.T) {
	setup()
	defer teardown()

	floatingIPs := []FloatingIP{{ID: testResourceID}}
	URL := fmt.Sprintf("/v1/availablefloatingips/%d/%d", projectID, regionID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(&floatingIPs)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.Floatingips.ListAvailable(ctx)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, floatingIPs) {
		t.Errorf("Floatingips.ListAvailable\n returned %+v,\n expected %+v", resp, floatingIPs)
	}
}

func TestFloatingips_MetadataList(t *testing.T) {
	setup()
	defer teardown()

	metadataList := []MetadataDetailed{{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}}
	URL := fmt.Sprintf("/v1/floatingips/%d/%d/%s/%s", projectID, regionID, testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(metadataList)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.Floatingips.MetadataList(ctx, testResourceID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, metadataList) {
		t.Errorf("Floatingips.MetadataList\n returned %+v,\n expected %+v", resp, metadataList)
	}
}

func TestFloatingips_MetadataCreate(t *testing.T) {
	setup()
	defer teardown()

	metadataCreateRequest := &MetadataCreateRequest{
		map[string]interface{}{"key": "value"},
	}
	URL := fmt.Sprintf("/v1/floatingips/%d/%d/%s/%s", projectID, regionID, testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
	})

	_, err := client.Floatingips.MetadataCreate(ctx, testResourceID, metadataCreateRequest)
	require.NoError(t, err)
}

func TestFloatingips_MetadataUpdate(t *testing.T) {
	setup()
	defer teardown()

	metadataCreateRequest := &MetadataCreateRequest{
		map[string]interface{}{"key": "value"},
	}
	URL := fmt.Sprintf("/v1/floatingips/%d/%d/%s/%s", projectID, regionID, testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
	})

	_, err := client.Floatingips.MetadataUpdate(ctx, testResourceID, metadataCreateRequest)
	require.NoError(t, err)
}

func TestFloatingips_MetadataDeleteItem(t *testing.T) {
	setup()
	defer teardown()

	URL := fmt.Sprintf("/v1/floatingips/%d/%d/%s/%s", projectID, regionID, testResourceID, metadataItemPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.Floatingips.MetadataDeleteItem(ctx, testResourceID, nil)
	require.NoError(t, err)
}

func TestFloatingips_MetadataGetItem(t *testing.T) {
	setup()
	defer teardown()

	metadata := &MetadataDetailed{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}
	URL := fmt.Sprintf("/v1/floatingips/%d/%d/%s/%s", projectID, regionID, testResourceID, metadataItemPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(metadata)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Floatingips.MetadataGetItem(ctx, testResourceID, nil)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, metadata) {
		t.Errorf("Floatingips.MetadataGetItem\n returned %+v,\n expected %+v", resp, metadata)
	}
}
