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

	const (
		portID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	allowedAddressPairsRequest := &AllowedAddressPairsRequest{
		IPAddress: "192.168.0.1",
	}

	portResponse := &Port{PortID: portID}

	assignPortURL := fmt.Sprintf("/v1/ports/%d/%d/%s/%s", projectID, regionID, portID, portsAllowAddressPairs)

	mux.HandleFunc(assignPortURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		reqBody := new(AllowedAddressPairsRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, allowedAddressPairsRequest, reqBody)
		resp, _ := json.Marshal(portResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Ports.Assign(ctx, portID, allowedAddressPairsRequest)
	require.NoError(t, err)

	assert.Equal(t, portResponse, resp)
}

func TestPorts_EnablePortSecurity(t *testing.T) {
	setup()
	defer teardown()

	const (
		portID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	instancePortInterfaceResponse := &InstancePortInterface{PortID: portID}

	enablePortSecurityURL := fmt.Sprintf("/v1/ports/%d/%d/%s/%s", projectID, regionID, portID, portsEnableSecurity)

	mux.HandleFunc(enablePortSecurityURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, _ := json.Marshal(instancePortInterfaceResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Ports.EnablePortSecurity(ctx, portID)
	require.NoError(t, err)

	assert.Equal(t, instancePortInterfaceResponse, resp)
}

func TestPorts_DisablePortSecurity(t *testing.T) {
	setup()
	defer teardown()

	const (
		portID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	instancePortInterfaceResponse := &InstancePortInterface{PortID: portID}

	enablePortSecurityURL := fmt.Sprintf("/v1/ports/%d/%d/%s/%s", projectID, regionID, portID, portsDisableSecurity)

	mux.HandleFunc(enablePortSecurityURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, _ := json.Marshal(instancePortInterfaceResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Ports.DisablePortSecurity(ctx, portID)
	require.NoError(t, err)

	assert.Equal(t, instancePortInterfaceResponse, resp)
}
