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

func TestServerGroups_Get(t *testing.T) {
	setup()
	defer teardown()

	const (
		sgID      = "f0d19cec-5c3f-4853-886e-304915960ff6"
		projectID = "27520"
		regionID  = "8"
	)

	servergroup := &ServerGroup{ID: sgID}
	getServerGroupURL := fmt.Sprintf("/v1/servergroups/%s/%s/%s", projectID, regionID, sgID)

	mux.HandleFunc(getServerGroupURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(servergroup)
		_, _ = fmt.Fprintf(w, `{"server_group":%s}`, string(resp))
	})

	opts := ServicePath{Project: projectID, Region: regionID}
	resp, _, err := client.ServerGroups.Get(ctx, sgID, &opts)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, servergroup) {
		t.Errorf("ServerGroups.Get\n returned %+v,\n expected %+v", resp, servergroup)
	}
}

func TestServerGroups_Create(t *testing.T) {
	setup()
	defer teardown()

	const (
		serverGroupID = "f0d19cec-5c3f-4853-886e-304915960ff6"
		projectID     = "27520"
		regionID      = "8"
	)

	servergroupCreateRequest := &ServerGroupCreateRequest{
		Name: "test-subnet",
	}

	serverGroupResponse := &ServerGroup{ID: serverGroupID}

	createServerGroupURL := fmt.Sprintf("/v1/servergroups/%s/%s", projectID, regionID)

	mux.HandleFunc(createServerGroupURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(ServerGroupCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, servergroupCreateRequest, reqBody)
		resp, _ := json.Marshal(serverGroupResponse)
		_, _ = fmt.Fprintf(w, `{"server_group":%s}`, string(resp))
	})

	opts := ServicePath{Project: projectID, Region: regionID}
	resp, _, err := client.ServerGroups.Create(ctx, servergroupCreateRequest, &opts)
	require.NoError(t, err)

	assert.Equal(t, serverGroupResponse, resp)
}

func TestServerGroups_Delete(t *testing.T) {
	setup()
	defer teardown()

	const (
		sgID      = "f0d19cec-5c3f-4853-886e-304915960ff6"
		projectID = "27520"
		regionID  = "8"
	)

	deleteServerGroupURL := fmt.Sprintf("/v1/servergroups/%s/%s/%s", projectID, regionID, sgID)

	mux.HandleFunc(deleteServerGroupURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	opts := ServicePath{Project: projectID, Region: regionID}
	_, err := client.ServerGroups.Delete(ctx, sgID, &opts)
	require.NoError(t, err)
}