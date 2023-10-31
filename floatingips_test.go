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

	const (
		fipID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	floatingIPs := []FloatingIP{{ID: fipID}}
	getFloatingipsURL := fmt.Sprintf("/v1/floatingips/%d/%d", projectID, regionID)

	mux.HandleFunc(getFloatingipsURL, func(w http.ResponseWriter, r *http.Request) {
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

	const (
		fipID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	floatingIP := &FloatingIP{ID: fipID}
	getFloatingipsURL := fmt.Sprintf("/v1/floatingips/%d/%d/%s", projectID, regionID, fipID)

	mux.HandleFunc(getFloatingipsURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(floatingIP)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Floatingips.Get(ctx, fipID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, floatingIP) {
		t.Errorf("Floatingips.Get\n returned %+v,\n expected %+v", resp, floatingIP)
	}
}

func TestFloatingips_Create(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	floatingIPCreateRequest := &FloatingIPCreateRequest{}

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	createFloatingipsURL := fmt.Sprintf("/v1/floatingips/%d/%d", projectID, regionID)

	mux.HandleFunc(createFloatingipsURL, func(w http.ResponseWriter, r *http.Request) {
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

	const (
		taskID = "f0d19cec-5c3f-4853-886e-304915960ff6"
		fipID  = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	deleteFloatingipsURL := fmt.Sprintf("/v1/floatingips/%d/%d/%s", projectID, regionID, fipID)

	mux.HandleFunc(deleteFloatingipsURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Floatingips.Delete(ctx, fipID)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestFloatingips_Assign(t *testing.T) {
	setup()
	defer teardown()

	const (
		fipID  = "f0d19cec-5c3f-4853-886e-304915960ff6"
		portID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	assignRequest := &AssignRequest{PortID: portID}

	floatingIP := &FloatingIP{
		ID:     fipID,
		PortID: portID,
	}
	assignFloatingipsURL := fmt.Sprintf("/v1/floatingips/%d/%d/%s/%s", projectID, regionID, fipID, floatingipsAssign)

	mux.HandleFunc(assignFloatingipsURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, _ := json.Marshal(&floatingIP)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Floatingips.Assign(ctx, fipID, assignRequest)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, floatingIP) {
		t.Errorf("Floatingips.Assign\n returned %+v,\n expected %+v", resp, floatingIP)
	}
}

func TestFloatingips_UnAssign(t *testing.T) {
	setup()
	defer teardown()

	const (
		fipID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	floatingIP := &FloatingIP{ID: fipID}
	assignFloatingipsURL := fmt.Sprintf("/v1/floatingips/%d/%d/%s/%s", projectID, regionID, fipID, floatingipsUnAssign)

	mux.HandleFunc(assignFloatingipsURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, _ := json.Marshal(&floatingIP)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Floatingips.UnAssign(ctx, fipID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, floatingIP) {
		t.Errorf("Floatingips.UnAssign\n returned %+v,\n expected %+v", resp, floatingIP)
	}
}

func TestFloatingips_ListAvailable(t *testing.T) {
	setup()
	defer teardown()

	const (
		fipID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	floatingIPs := []FloatingIP{{ID: fipID}}
	getFloatingipsURL := fmt.Sprintf("/v1/availablefloatingips/%d/%d", projectID, regionID)

	mux.HandleFunc(getFloatingipsURL, func(w http.ResponseWriter, r *http.Request) {
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

	const (
		fipID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	metadataList := []MetadataDetailed{{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}}
	getFloatingipsMetadataListURL := fmt.Sprintf("/v1/floatingips/%d/%d/%s/%s", projectID, regionID, fipID, metadataPath)

	mux.HandleFunc(getFloatingipsMetadataListURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(metadataList)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.Floatingips.MetadataList(ctx, fipID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, metadataList) {
		t.Errorf("Floatingips.MetadataList\n returned %+v,\n expected %+v", resp, metadataList)
	}
}

func TestFloatingips_MetadataCreate(t *testing.T) {
	setup()
	defer teardown()

	const (
		fipID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	metadataCreateRequest := &MetadataCreateRequest{
		map[string]interface{}{"key": "value"},
	}

	createFloatingipsMetadataURL := fmt.Sprintf("/v1/floatingips/%d/%d/%s/%s", projectID, regionID, fipID, metadataPath)

	mux.HandleFunc(createFloatingipsMetadataURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
	})

	_, err := client.Floatingips.MetadataCreate(ctx, fipID, metadataCreateRequest)
	require.NoError(t, err)
}

func TestFloatingips_MetadataUpdate(t *testing.T) {
	setup()
	defer teardown()

	const (
		fipID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	metadataCreateRequest := &MetadataCreateRequest{
		map[string]interface{}{"key": "value"},
	}

	createFloatingipsMetadataURL := fmt.Sprintf("/v1/floatingips/%d/%d/%s/%s", projectID, regionID, fipID, metadataPath)

	mux.HandleFunc(createFloatingipsMetadataURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
	})

	_, err := client.Floatingips.MetadataUpdate(ctx, fipID, metadataCreateRequest)
	require.NoError(t, err)
}

func TestFloatingips_MetadataDeleteItem(t *testing.T) {
	setup()
	defer teardown()

	const (
		fipID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	deleteFloatingipsMetadataItemURL := fmt.Sprintf("/v1/floatingips/%d/%d/%s/%s", projectID, regionID, fipID, metadataItemPath)

	mux.HandleFunc(deleteFloatingipsMetadataItemURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.Floatingips.MetadataDeleteItem(ctx, fipID, nil)
	require.NoError(t, err)
}

func TestFloatingips_MetadataGetItem(t *testing.T) {
	setup()
	defer teardown()

	const (
		fipID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	metadata := &MetadataDetailed{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}
	getFloatingipsMetadataItemURL := fmt.Sprintf("/v1/floatingips/%d/%d/%s/%s", projectID, regionID, fipID, metadataItemPath)

	mux.HandleFunc(getFloatingipsMetadataItemURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(metadata)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Floatingips.MetadataGetItem(ctx, fipID, nil)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, metadata) {
		t.Errorf("Floatingips.MetadataGetItem\n returned %+v,\n expected %+v", resp, metadata)
	}
}
