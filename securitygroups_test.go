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

func TestSecurityGroups_Update(t *testing.T) {
	setup()
	defer teardown()

	securityGroupUpdateRequest := &SecurityGroupUpdateRequest{}
	securityGroup := &SecurityGroup{ID: testResourceID}
	URL := fmt.Sprintf("/v1/securitygroups/%d/%d/%s", projectID, regionID, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		reqBody := new(SecurityGroupUpdateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, securityGroupUpdateRequest, reqBody)
		resp, _ := json.Marshal(securityGroup)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.SecurityGroups.Update(ctx, testResourceID, securityGroupUpdateRequest)
	require.NoError(t, err)

	assert.Equal(t, securityGroup, resp)
}

func TestSecurityGroups_DeepCopy(t *testing.T) {
	setup()
	defer teardown()

	securityGroupDeepCopyRequest := &SecurityGroupDeepCopyRequest{}
	URL := fmt.Sprintf("/v1/securitygroups/%d/%d/%s/%s", projectID, regionID, testResourceID, securitygroupsCopyPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(SecurityGroupDeepCopyRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, securityGroupDeepCopyRequest, reqBody)
	})

	resp, err := client.SecurityGroups.DeepCopy(ctx, testResourceID, securityGroupDeepCopyRequest)
	require.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestSecurityGroups_RuleCreate(t *testing.T) {
	setup()
	defer teardown()

	ruleCreateRequest := &RuleCreateRequest{}
	securityGroupRule := &SecurityGroupRule{ID: testResourceID}
	URL := fmt.Sprintf("/v1/securitygroups/%d/%d/%s/%s", projectID, regionID, testResourceID, securitygroupsRulesPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(RuleCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, ruleCreateRequest, reqBody)
		resp, _ := json.Marshal(securityGroupRule)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.SecurityGroups.RuleCreate(ctx, testResourceID, ruleCreateRequest)
	require.NoError(t, err)

	assert.Equal(t, securityGroupRule, resp)
}

func TestSecurityGroups_RuleDelete(t *testing.T) {
	setup()
	defer teardown()

	taskResponse := &TaskResponse{Tasks: []string{taskID}}
	URL := fmt.Sprintf("/v1/securitygrouprules/%d/%d/%s", projectID, regionID, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.SecurityGroups.RuleDelete(ctx, testResourceID)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestSecurityGroups_RuleUpdate(t *testing.T) {
	setup()
	defer teardown()

	ruleUpdateRequest := &RuleUpdateRequest{}
	securityGroupRule := &SecurityGroupRule{ID: testResourceID}
	URL := fmt.Sprintf("/v1/securitygrouprules/%d/%d/%s", projectID, regionID, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		reqBody := new(RuleUpdateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, ruleUpdateRequest, reqBody)
		resp, _ := json.Marshal(securityGroupRule)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.SecurityGroups.RuleUpdate(ctx, testResourceID, ruleUpdateRequest)
	require.NoError(t, err)

	assert.Equal(t, securityGroupRule, resp)
}

func TestSecurityGroups_MetadataList(t *testing.T) {
	setup()
	defer teardown()

	metadataList := []MetadataDetailed{{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}}
	URL := fmt.Sprintf("/v1/securitygroups/%d/%d/%s/%s", projectID, regionID, testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(metadataList)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.SecurityGroups.MetadataList(ctx, testResourceID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, metadataList) {
		t.Errorf("SecurityGroups.MetadataList\n returned %+v,\n expected %+v", resp, metadataList)
	}
}

func TestSecurityGroups_MetadataCreate(t *testing.T) {
	setup()
	defer teardown()

	metadataCreateRequest := &MetadataCreateRequest{
		map[string]interface{}{"key": "value"},
	}
	URL := fmt.Sprintf("/v1/securitygroups/%d/%d/%s/%s", projectID, regionID, testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
	})

	_, err := client.SecurityGroups.MetadataCreate(ctx, testResourceID, metadataCreateRequest)
	require.NoError(t, err)
}

func TestSecurityGroups_MetadataUpdate(t *testing.T) {
	setup()
	defer teardown()

	metadataCreateRequest := &MetadataCreateRequest{
		map[string]interface{}{"key": "value"},
	}
	URL := fmt.Sprintf("/v1/securitygroups/%d/%d/%s/%s", projectID, regionID, testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
	})

	_, err := client.SecurityGroups.MetadataUpdate(ctx, testResourceID, metadataCreateRequest)
	require.NoError(t, err)
}

func TestSecurityGroups_MetadataDeleteItem(t *testing.T) {
	setup()
	defer teardown()

	URL := fmt.Sprintf("/v1/securitygroups/%d/%d/%s/%s", projectID, regionID, testResourceID, metadataItemPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.SecurityGroups.MetadataDeleteItem(ctx, testResourceID, nil)
	require.NoError(t, err)
}

func TestSecurityGroups_MetadataGetItem(t *testing.T) {
	setup()
	defer teardown()

	metadata := &MetadataDetailed{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}
	URL := fmt.Sprintf("/v1/securitygroups/%d/%d/%s/%s", projectID, regionID, testResourceID, metadataItemPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(metadata)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.SecurityGroups.MetadataGetItem(ctx, testResourceID, nil)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, metadata) {
		t.Errorf("SecurityGroups.MetadataGetItem\n returned %+v,\n expected %+v", resp, metadata)
	}
}
