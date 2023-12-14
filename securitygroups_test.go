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

func TestSecurityGroups_List(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []SecurityGroup{{ID: testResourceID}}
	URL := path.Join(securitygroupsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.SecurityGroups.List(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestSecurityGroups_List_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(securitygroupsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.SecurityGroups.List(ctx, nil)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestSecurityGroups_Get(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &SecurityGroup{ID: testResourceID}
	URL := path.Join(securitygroupsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.SecurityGroups.Get(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestSecurityGroups_Get_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(securitygroupsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.SecurityGroups.Get(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestSecurityGroups_Create(t *testing.T) {
	setup()
	defer teardown()

	request := &SecurityGroupCreateRequest{}
	expectedResp := &SecurityGroup{ID: testResourceID}
	URL := path.Join(securitygroupsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(SecurityGroupCreateRequest)
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

	respActual, resp, err := client.SecurityGroups.Create(ctx, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestSecurityGroups_Create_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &SecurityGroupCreateRequest{}
	URL := path.Join(securitygroupsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.SecurityGroups.Create(ctx, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestSecurityGroups_Create_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.SecurityGroups.Create(ctx, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestSecurityGroups_Delete(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(securitygroupsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, err := client.SecurityGroups.Delete(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestSecurityGroups_Delete_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(securitygroupsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)
	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	resp, err := client.SecurityGroups.Delete(ctx, testResourceID)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestSecurityGroups_Update(t *testing.T) {
	setup()
	defer teardown()

	request := &SecurityGroupUpdateRequest{}
	expectedResp := &SecurityGroup{ID: testResourceID}
	URL := path.Join(securitygroupsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		reqBody := new(SecurityGroupUpdateRequest)
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

	respActual, resp, err := client.SecurityGroups.Update(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestSecurityGroups_Update_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &SecurityGroupUpdateRequest{}
	URL := path.Join(securitygroupsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.SecurityGroups.Update(ctx, testResourceID, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestSecurityGroups_Update_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.SecurityGroups.Update(ctx, testResourceID, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestSecurityGroups_DeepCopy(t *testing.T) {
	setup()
	defer teardown()

	request := &Name{}
	URL := path.Join(securitygroupsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, securitygroupsCopy)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(Name)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, reqBody)
	})

	resp, err := client.SecurityGroups.DeepCopy(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestSecurityGroups_DeepCopy_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &Name{}
	URL := path.Join(securitygroupsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, securitygroupsCopy)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	resp, err := client.SecurityGroups.DeepCopy(ctx, testResourceID, request)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestSecurityGroups_DeepCopy_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Floatingips.Assign(ctx, testResourceID, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestSecurityGroups_RuleCreate(t *testing.T) {
	setup()
	defer teardown()

	request := &RuleCreateRequest{}
	expectedResp := &SecurityGroupRule{ID: testResourceID}
	URL := path.Join(securitygroupsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, securitygroupsRules)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(RuleCreateRequest)
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

	respActual, resp, err := client.SecurityGroups.RuleCreate(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestSecurityGroups_RuleCreate_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &RuleCreateRequest{}
	URL := path.Join(securitygroupsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, securitygroupsRules)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.SecurityGroups.RuleCreate(ctx, testResourceID, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestSecurityGroups_RuleCreate_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.SecurityGroups.RuleCreate(ctx, testResourceID, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestSecurityGroups_RuleDelete(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(securitygroupsRulesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.SecurityGroups.RuleDelete(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestSecurityGroups_RuleDelete_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(securitygroupsRulesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.SecurityGroups.RuleDelete(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestSecurityGroups_RuleUpdate(t *testing.T) {
	setup()
	defer teardown()

	request := &RuleUpdateRequest{}
	expectedResp := &SecurityGroupRule{ID: testResourceID}
	URL := path.Join(securitygroupsRulesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		reqBody := new(RuleUpdateRequest)
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

	respActual, resp, err := client.SecurityGroups.RuleUpdate(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestSecurityGroups_RuleUpdate_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &RuleUpdateRequest{}
	URL := path.Join(securitygroupsRulesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.SecurityGroups.RuleUpdate(ctx, testResourceID, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestSecurityGroups_RuleUpdate_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.SecurityGroups.RuleUpdate(ctx, testResourceID, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestSecurityGroups_MetadataList(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []MetadataDetailed{{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}}
	URL := path.Join(securitygroupsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.SecurityGroups.MetadataList(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestSecurityGroups_MetadataCreate(t *testing.T) {
	setup()
	defer teardown()

	request := &Metadata{"key": "value"}
	URL := path.Join(securitygroupsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
	})

	resp, err := client.SecurityGroups.MetadataCreate(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestSecurityGroups_MetadataUpdate(t *testing.T) {
	setup()
	defer teardown()

	request := &Metadata{"key": "value"}
	URL := path.Join(securitygroupsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
	})

	resp, err := client.SecurityGroups.MetadataUpdate(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestSecurityGroups_MetadataDeleteItem(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(securitygroupsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataItemPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	resp, err := client.SecurityGroups.MetadataDeleteItem(ctx, testResourceID, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestSecurityGroups_MetadataGetItem(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &MetadataDetailed{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}
	URL := path.Join(securitygroupsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataItemPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.SecurityGroups.MetadataGetItem(ctx, testResourceID, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestSecurityGroups_isValidUUID_Error_Return_SecurityGroup(t *testing.T) {
	setup()
	defer teardown()

	tests := []struct {
		name     string
		testFunc func() (*SecurityGroup, *Response, error)
	}{
		{
			name: "Get",
			testFunc: func() (*SecurityGroup, *Response, error) {
				return client.SecurityGroups.Get(ctx, testResourceIDNotValidUUID)
			},
		},
		{
			name: "Update",
			testFunc: func() (*SecurityGroup, *Response, error) {
				return client.SecurityGroups.Update(ctx, testResourceIDNotValidUUID, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respActual, resp, err := tt.testFunc()
			require.Nil(t, respActual)
			require.Equal(t, 400, resp.StatusCode)
			require.EqualError(t, err, NewArgError("securityGroupID", NotCorrectUUID).Error())
		})
	}
}

func TestSecurityGroupsRule_isValidUUID_Error_Return_SecurityGroupRule(t *testing.T) {
	setup()
	defer teardown()

	tests := []struct {
		name     string
		testFunc func() (*SecurityGroupRule, *Response, error)
	}{
		{
			name: "RuleCreate",
			testFunc: func() (*SecurityGroupRule, *Response, error) {
				return client.SecurityGroups.RuleCreate(ctx, testResourceIDNotValidUUID, nil)
			},
		},
		{
			name: "RuleUpdate",
			testFunc: func() (*SecurityGroupRule, *Response, error) {
				return client.SecurityGroups.RuleUpdate(ctx, testResourceIDNotValidUUID, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respActual, resp, err := tt.testFunc()
			require.Nil(t, respActual)
			require.Equal(t, 400, resp.StatusCode)
			require.EqualError(t, err, NewArgError("securityGroupID", NotCorrectUUID).Error())
		})
	}
}

func TestSecurityGroupsRule_Delete_isValidUUID_Error(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.SecurityGroups.RuleDelete(ctx, testResourceIDNotValidUUID)
	assert.Nil(t, respActual)
	require.Equal(t, 400, resp.StatusCode)
	require.EqualError(t, err, NewArgError("securityGroupID", NotCorrectUUID).Error())
}

func TestSecurityGroups_Delete_isValidUUID_Error(t *testing.T) {
	setup()
	defer teardown()

	resp, err := client.SecurityGroups.Delete(ctx, testResourceIDNotValidUUID)
	require.Equal(t, 400, resp.StatusCode)
	require.EqualError(t, err, NewArgError("securityGroupID", NotCorrectUUID).Error())
}

func TestSecurityGroups_DeepCopy_isValidUUID_Error(t *testing.T) {
	setup()
	defer teardown()

	resp, err := client.SecurityGroups.DeepCopy(ctx, testResourceIDNotValidUUID, nil)
	require.Equal(t, 400, resp.StatusCode)
	require.EqualError(t, err, NewArgError("securityGroupID", NotCorrectUUID).Error())
}

func TestSecurityGroups_Metadata_isValidUUID_Error_Return_Response(t *testing.T) {
	setup()
	defer teardown()

	tests := []struct {
		name     string
		testFunc func() (*Response, error)
	}{
		{
			name: "MetadataCreate",
			testFunc: func() (*Response, error) {
				return client.SecurityGroups.MetadataCreate(ctx, testResourceIDNotValidUUID, nil)
			},
		},
		{
			name: "MetadataUpdate",
			testFunc: func() (*Response, error) {
				return client.SecurityGroups.MetadataUpdate(ctx, testResourceIDNotValidUUID, nil)
			},
		},
		{
			name: "MetadataDeleteItem",
			testFunc: func() (*Response, error) {
				return client.SecurityGroups.MetadataDeleteItem(ctx, testResourceIDNotValidUUID, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.testFunc()
			require.Equal(t, 400, resp.StatusCode)
			require.EqualError(t, err, NewArgError("securityGroupID", NotCorrectUUID).Error())
		})
	}
}

func TestSecurityGroups_MetadataList_isValidUUID_Error(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.SecurityGroups.MetadataList(ctx, testResourceIDNotValidUUID)
	require.Nil(t, respActual)
	require.Equal(t, 400, resp.StatusCode)
	require.EqualError(t, err, NewArgError("securityGroupID", NotCorrectUUID).Error())
}

func TestSecurityGroups_MetadataGetItem_isValidUUID_Error(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.SecurityGroups.MetadataGetItem(ctx, testResourceIDNotValidUUID, nil)
	require.Nil(t, respActual)
	require.Equal(t, 400, resp.StatusCode)
	require.EqualError(t, err, NewArgError("securityGroupID", NotCorrectUUID).Error())
}
