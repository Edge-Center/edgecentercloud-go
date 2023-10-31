package edgecloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPorts_Assign(t *testing.T) {
	setup()
	defer teardown()

	allowedAddressPairsRequest := &AllowedAddressPairsRequest{
		IPAddress: "192.168.0.1",
	}
	portResponse := &Port{PortID: testResourceID}
	URL := fmt.Sprintf("/v1/ports/%d/%d/%s/%s", projectID, regionID, testResourceID, portsAllowAddressPairs)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		reqBody := new(AllowedAddressPairsRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, allowedAddressPairsRequest, reqBody)
		resp, _ := json.Marshal(portResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Ports.Assign(ctx, testResourceID, allowedAddressPairsRequest)
	require.NoError(t, err)

	assert.Equal(t, portResponse, resp)
}

func TestPorts_EnablePortSecurity(t *testing.T) {
	setup()
	defer teardown()

	instancePortInterfaceResponse := &InstancePortInterface{PortID: testResourceID}
	URL := fmt.Sprintf("/v1/ports/%d/%d/%s/%s", projectID, regionID, testResourceID, portsEnableSecurity)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, _ := json.Marshal(instancePortInterfaceResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Ports.EnablePortSecurity(ctx, testResourceID)
	require.NoError(t, err)

	assert.Equal(t, instancePortInterfaceResponse, resp)
}

func TestPorts_DisablePortSecurity(t *testing.T) {
	setup()
	defer teardown()

	instancePortInterfaceResponse := &InstancePortInterface{PortID: testResourceID}
	URL := fmt.Sprintf("/v1/ports/%d/%d/%s/%s", projectID, regionID, testResourceID, portsDisableSecurity)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, _ := json.Marshal(instancePortInterfaceResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Ports.DisablePortSecurity(ctx, testResourceID)
	require.NoError(t, err)

	assert.Equal(t, instancePortInterfaceResponse, resp)
}
