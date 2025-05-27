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

func TestKeyPairs_List(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []KeyPair{{PublicKey: "ssh-key"}}
	URL := path.Join(keypairsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.KeyPairs.List(ctx)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestKeyPairsV2_List(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []KeyPairV2{{PublicKey: "ssh-key"}}
	URL := path.Join(keypairsBasePathV2)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.KeyPairs.ListV2(ctx)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestKeyPairs_Get(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &KeyPair{SSHKeyID: testResourceID}
	URL := path.Join(keypairsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.KeyPairs.Get(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestKeyPairsV2_Get(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &KeyPairV2{SSHKeyID: testResourceID}
	URL := path.Join(keypairsBasePathV2, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.KeyPairs.GetV2(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestKeyPairs_Create(t *testing.T) {
	setup()
	defer teardown()

	request := &KeyPairCreateRequest{SSHKeyName: "ssh-key"}
	expectedResp := &KeyPair{SSHKeyName: "ssh-key"}
	URL := path.Join(keypairsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(KeyPairCreateRequest)
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

	respActual, resp, err := client.KeyPairs.Create(ctx, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestKeyPairsV2_Create(t *testing.T) {
	setup()
	defer teardown()

	request := &KeyPairCreateRequestV2{SSHKeyName: "ssh-key", ProjectID: 1234}
	expectedResp := &KeyPairV2{SSHKeyName: "ssh-key"}
	URL := path.Join(keypairsBasePathV2)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(KeyPairCreateRequestV2)
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

	respActual, resp, err := client.KeyPairs.CreateV2(ctx, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestKeyPairs_Delete(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(keypairsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.KeyPairs.Delete(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestKeyPairsV2_Delete(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(keypairsBasePathV2, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	resp, err := client.KeyPairs.DeleteV2(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestKeyPairs_Share(t *testing.T) {
	setup()
	defer teardown()

	request := &KeyPairShareRequest{SharedInProject: true}
	expectedResp := &KeyPair{SSHKeyID: testResourceID}
	URL := path.Join(keypairsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, keypairsShare)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		reqBody := new(KeyPairShareRequest)
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

	respActual, resp, err := client.KeyPairs.Share(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}
