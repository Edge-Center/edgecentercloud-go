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

func TestL7Rules_List(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []L7Rule{{ID: testResourceID}}
	URL := path.Join(l7policiesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, l7rulesPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.L7Rules.List(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestL7Rules_Create(t *testing.T) {
	setup()
	defer teardown()

	request := &L7RuleCreateRequest{Key: "test-l7rule"}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(l7policiesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, l7rulesPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(L7RuleCreateRequest)
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

	respActual, resp, err := client.L7Rules.Create(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestL7Rules_Delete(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(l7policiesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, l7rulesPath, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.L7Rules.Delete(ctx, testResourceID, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestL7Rules_Get(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &L7Rule{ID: testResourceID}
	URL := path.Join(l7policiesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, l7rulesPath, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.L7Rules.Get(ctx, testResourceID, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestL7Rules_Update(t *testing.T) {
	setup()
	defer teardown()

	request := &L7RuleUpdateRequest{}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(l7policiesBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, l7rulesPath, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		reqBody := new(L7RuleUpdateRequest)
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

	respActual, resp, err := client.L7Rules.Update(ctx, testResourceID, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}
