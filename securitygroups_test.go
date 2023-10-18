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

func TestSecurityGroups_Get(t *testing.T) {
	setup()
	defer teardown()

	const (
		sgID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	securityGroup := &SecurityGroup{ID: sgID}
	getSecurityGroupsURL := fmt.Sprintf("/v1/securitygroups/%d/%d/%s", projectID, regionID, sgID)

	mux.HandleFunc(getSecurityGroupsURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(securityGroup)
		_, _ = fmt.Fprintf(w, `{"security_group":%s}`, string(resp))
	})

	resp, _, err := client.SecurityGroups.Get(ctx, sgID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, securityGroup) {
		t.Errorf("SecurityGroups.Get\n returned %+v,\n expected %+v", resp, securityGroup)
	}
}

func TestSecurityGroups_Create(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	securityGroupCreateRequest := &SecurityGroupCreateRequest{}

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	createSecurityGroupsURL := fmt.Sprintf("/v1/securitygroups/%d/%d", projectID, regionID)

	mux.HandleFunc(createSecurityGroupsURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(SecurityGroupCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, securityGroupCreateRequest, reqBody)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprintf(w, `{"tasks":%s}`, string(resp))
	})

	resp, _, err := client.SecurityGroups.Create(ctx, securityGroupCreateRequest)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestSecurityGroups_Delete(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID = "f0d19cec-5c3f-4853-886e-304915960ff6"
		sgID   = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	deleteSecurityGroupsURL := fmt.Sprintf("/v1/securitygroups/%d/%d/%s", projectID, regionID, sgID)

	mux.HandleFunc(deleteSecurityGroupsURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprintf(w, `{"tasks":%s}`, string(resp))
	})

	resp, _, err := client.SecurityGroups.Delete(ctx, sgID)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}
