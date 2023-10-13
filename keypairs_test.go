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

func TestKeyPairs_Get(t *testing.T) {
	setup()
	defer teardown()

	const (
		keypairID = "f0d19cec-5c3f-4853-886e-304915960ff6"
		projectID = "27520"
		regionID  = "8"
	)

	keypair := &KeyPair{SSHKeyID: keypairID}
	getKeyPairsURL := fmt.Sprintf("/v1/keypairs/%s/%s/%s", projectID, regionID, keypairID)

	mux.HandleFunc(getKeyPairsURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(keypair)
		_, _ = fmt.Fprintf(w, `{"keypair":%s}`, string(resp))
	})

	opts := ServicePath{Project: projectID, Region: regionID}
	resp, _, err := client.KeyPairs.Get(ctx, keypairID, &opts)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, keypair) {
		t.Errorf("KeyPairs.Get\n returned %+v,\n expected %+v", resp, keypair)
	}
}

func TestKeyPairs_Create(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID    = "f0d19cec-5c3f-4853-886e-304915960ff6"
		projectID = "27520"
		regionID  = "8"
	)

	keyPairCreateRequest := &KeyPairCreateRequest{
		SSHKeyName: "ssh-key",
	}

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	createKeyPairsURL := fmt.Sprintf("/v1/keypairs/%s/%s", projectID, regionID)

	mux.HandleFunc(createKeyPairsURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(KeyPairCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, keyPairCreateRequest, reqBody)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprintf(w, `{"tasks":%s}`, string(resp))
	})

	opts := ServicePath{Project: projectID, Region: regionID}
	resp, _, err := client.KeyPairs.Create(ctx, keyPairCreateRequest, &opts)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestKeyPairs_Delete(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID    = "f0d19cec-5c3f-4853-886e-304915960ff6"
		keypairID = "f0d19cec-5c3f-4853-886e-304915960ff6"
		projectID = "27520"
		regionID  = "8"
	)

	taskResponse := &TaskResponse{Tasks: []string{taskID}}

	deleteKeyPairsURL := fmt.Sprintf("/v1/keypairs/%s/%s/%s", projectID, regionID, keypairID)

	mux.HandleFunc(deleteKeyPairsURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprintf(w, `{"tasks":%s}`, string(resp))
	})

	opts := ServicePath{Project: projectID, Region: regionID}
	resp, _, err := client.KeyPairs.Delete(ctx, keypairID, &opts)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}
