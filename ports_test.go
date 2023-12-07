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

func TestPorts_Assign(t *testing.T) {
	setup()
	defer teardown()

	request := &AllowedAddressPairsRequest{IPAddress: "192.168.0.1"}
	expectedResp := &Port{PortID: testResourceID}
	URL := path.Join(portsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, portsAllowAddressPairs)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		reqBody := new(AllowedAddressPairsRequest)
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

	respActual, resp, err := client.Ports.Assign(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestPorts_EnablePortSecurity(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &InstancePortInterface{PortID: testResourceID}
	URL := path.Join(portsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, portsEnableSecurity)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Ports.EnablePortSecurity(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestPorts_DisablePortSecurity(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &InstancePortInterface{PortID: testResourceID}
	URL := path.Join(portsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, portsDisableSecurity)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Ports.DisablePortSecurity(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}
