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
		_, _ = fmt.Fprintf(w, `{"floating_ip":%s}`, string(resp))
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
		_, _ = fmt.Fprintf(w, `{"tasks":%s}`, string(resp))
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
		_, _ = fmt.Fprintf(w, `{"tasks":%s}`, string(resp))
	})

	resp, _, err := client.Floatingips.Delete(ctx, fipID)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}
