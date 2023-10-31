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

func TestSecurityGroups_List(t *testing.T) {
	setup()
	defer teardown()

	securityGroups := []SecurityGroup{{ID: testResourceID}}
	URL := fmt.Sprintf("/v1/securitygroups/%d/%d", projectID, regionID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(securityGroups)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.SecurityGroups.List(ctx, nil)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, securityGroups) {
		t.Errorf("SecurityGroups.List\n returned %+v,\n expected %+v", resp, securityGroups)
	}
}

func TestSecurityGroups_Get(t *testing.T) {
	setup()
	defer teardown()

	securityGroup := &SecurityGroup{ID: testResourceID}
	URL := fmt.Sprintf("/v1/securitygroups/%d/%d/%s", projectID, regionID, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(securityGroup)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.SecurityGroups.Get(ctx, testResourceID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, securityGroup) {
		t.Errorf("SecurityGroups.Get\n returned %+v,\n expected %+v", resp, securityGroup)
	}
}

func TestSecurityGroups_Create(t *testing.T) {
	setup()
	defer teardown()

	securityGroupCreateRequest := &SecurityGroupCreateRequest{}
	taskResponse := &TaskResponse{Tasks: []string{taskID}}
	URL := fmt.Sprintf("/v1/securitygroups/%d/%d", projectID, regionID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(SecurityGroupCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, securityGroupCreateRequest, reqBody)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.SecurityGroups.Create(ctx, securityGroupCreateRequest)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestSecurityGroups_Delete(t *testing.T) {
	setup()
	defer teardown()

	taskResponse := &TaskResponse{Tasks: []string{taskID}}
	URL := fmt.Sprintf("/v1/securitygroups/%d/%d/%s", projectID, regionID, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.SecurityGroups.Delete(ctx, testResourceID)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}
