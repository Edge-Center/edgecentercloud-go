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

func TestAvailabilityZonesServiceOp_List(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(AvailabilityZoneBasePath, strconv.Itoa(regionID))
	expResp := AvailabilityZonesList{
		RegionID:          regionID,
		AvailabilityZones: []string{"availability-zone-1", "availability-zone-2"},
	}

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, err = w.Write(resp)
		if err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	})

	respActual, resp, err := client.AvailabilityZones.List(ctx)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, expResp, *respActual)
}

func TestAvailabilityZonesServiceOp_List_NotFoundError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(AvailabilityZoneBasePath, strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, "Region not found")
	})

	respActual, resp, err := client.AvailabilityZones.List(ctx)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}
