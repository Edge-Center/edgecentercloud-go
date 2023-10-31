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

func TestServerGroups_List(t *testing.T) {
	setup()
	defer teardown()

	serverGroups := []ServerGroup{{ID: testResourceID}}
	URL := fmt.Sprintf("/v1/servergroups/%d/%d", projectID, regionID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(serverGroups)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.ServerGroups.List(ctx)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, serverGroups) {
		t.Errorf("ServerGroups.List\n returned %+v,\n expected %+v", resp, serverGroups)
	}
}

func TestServerGroups_Get(t *testing.T) {
	setup()
	defer teardown()

	servergroup := &ServerGroup{ID: testResourceID}
	URL := fmt.Sprintf("/v1/servergroups/%d/%d/%s", projectID, regionID, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(servergroup)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.ServerGroups.Get(ctx, testResourceID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, servergroup) {
		t.Errorf("ServerGroups.Get\n returned %+v,\n expected %+v", resp, servergroup)
	}
}

func TestServerGroups_Create(t *testing.T) {
	setup()
	defer teardown()

	servergroupCreateRequest := &ServerGroupCreateRequest{
		Name: "test-subnet",
	}
	serverGroupResponse := &ServerGroup{ID: testResourceID}
	URL := fmt.Sprintf("/v1/servergroups/%d/%d", projectID, regionID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(ServerGroupCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, servergroupCreateRequest, reqBody)
		resp, _ := json.Marshal(serverGroupResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.ServerGroups.Create(ctx, servergroupCreateRequest)
	require.NoError(t, err)

	assert.Equal(t, serverGroupResponse, resp)
}

func TestServerGroups_Delete(t *testing.T) {
	setup()
	defer teardown()

	URL := fmt.Sprintf("/v1/servergroups/%d/%d/%s", projectID, regionID, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.ServerGroups.Delete(ctx, testResourceID)
	require.NoError(t, err)
}
