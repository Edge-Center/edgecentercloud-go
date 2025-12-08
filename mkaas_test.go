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

const (
	mkaasTestNetworkID = "f0d19cec-5c3f-4853-886e-304915960ff4"
	mkaasSubnetID      = "rtb19cec-5c3f-4853-886e-3045fk9g5kg9"
	mkaasClusterName   = "test-cluster"
	mkaasPoolName      = "test-pool"
	mkaasKeypairName   = "keypair"
	mkaasFlavor        = "g1-standard-2-4"
	mkaasVersion       = "v1.31.0"
)

func TestMKaaSServiceOp_ClusterCreate(t *testing.T) {
	setup()
	defer teardown()

	request := MKaaSClusterCreateRequest{
		Name:           mkaasClusterName,
		SSHKeyPairName: mkaasKeypairName,
		NetworkID:      mkaasTestNetworkID,
		SubnetID:       mkaasSubnetID,
		ControlPlane: ControlPlaneCreateRequest{
			Flavor:     mkaasFlavor,
			NodeCount:  1,
			VolumeSize: 10,
			Version:    mkaasVersion,
		},
		Pools: []MKaaSPoolCreateRequest{
			{
				Name:         mkaasPoolName,
				Flavor:       mkaasFlavor,
				MaxNodeCount: PtrTo(5),
				MinNodeCount: PtrTo(1),
				NodeCount:    2,
				VolumeSize:   10,
				Labels:       nil,
				Taints:       nil,
			},
		},
	}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(MKaaSClustersBasePathV2, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := &MKaaSClusterCreateRequest{}
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, *reqBody)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.MkaaS.ClusterCreate(ctx, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestMKaaSServiceOp_ClustersList(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []MKaaSCluster{{ID: testResourceIntID}}
	URL := path.Join(MKaaSClustersBasePathV2, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.MkaaS.ClustersList(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestMKaaSServiceOp_ClustersList_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(MKaaSClustersBasePathV2, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.MkaaS.ClustersList(ctx, nil)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestMKaaSServiceOp_ClustersGet(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &MKaaSCluster{
		ID: testResourceIntID,
	}
	URL := path.Join(MKaaSClustersBasePathV2, strconv.Itoa(projectID), strconv.Itoa(regionID), strconv.Itoa(testResourceIntID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(expectedResp)
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.MkaaS.ClusterGet(ctx, testResourceIntID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestMKaaSServiceOp_ClustersDelete(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(MKaaSClustersBasePathV2, strconv.Itoa(projectID), strconv.Itoa(regionID), strconv.Itoa(testResourceIntID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.MkaaS.ClusterDelete(ctx, testResourceIntID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestMKaaSServiceOp_PoolCreate(t *testing.T) {
	setup()
	defer teardown()

	request := MKaaSPoolCreateRequest{
		Name:       mkaasClusterName,
		Flavor:     mkaasFlavor,
		NodeCount:  1,
		VolumeSize: 10,
	}

	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(MKaaSClustersBasePathV2, strconv.Itoa(projectID), strconv.Itoa(regionID), strconv.Itoa(testResourceIntID), "pools")

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := &MKaaSPoolCreateRequest{}
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, *reqBody)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.MkaaS.PoolCreate(ctx, testResourceIntID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestMKaaSServiceOp_PoolsList(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []MKaaSPool{{ID: testResourceIntID}}
	URL := path.Join(MKaaSClustersBasePathV2, strconv.Itoa(projectID), strconv.Itoa(regionID), strconv.Itoa(testResourceIntID), "pools")

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.MkaaS.PoolsList(ctx, testResourceIntID, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestMKaaSServiceOp_PoolGet(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &MKaaSPool{
		ID: testResourceIntID,
	}
	URL := path.Join(MKaaSClustersBasePathV2, strconv.Itoa(projectID), strconv.Itoa(regionID), strconv.Itoa(testResourceIntID), "pools", strconv.Itoa(testResourceIntID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(expectedResp)
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.MkaaS.PoolGet(ctx, testResourceIntID, testResourceIntID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestMKaaSServiceOp_PoolDelete(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(MKaaSClustersBasePathV2, strconv.Itoa(projectID), strconv.Itoa(regionID),
		strconv.Itoa(testResourceIntID), "pools", strconv.Itoa(testResourceIntID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.MkaaS.PoolDelete(ctx, testResourceIntID, testResourceIntID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestMKaaSServiceOp_ClusterUpdate(t *testing.T) {
	setup()
	defer teardown()

	request := MKaaSClusterUpdateRequest{
		Name:            "updated-cluster",
		MasterNodeCount: 3,
	}

	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(MKaaSClustersBasePathV2, strconv.Itoa(projectID), strconv.Itoa(regionID), strconv.Itoa(testResourceIntID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		reqBody := &MKaaSClusterUpdateRequest{}
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, *reqBody)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.MkaaS.ClusterUpdate(ctx, testResourceIntID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestMKaaSServiceOp_PoolUpdate(t *testing.T) {
	setup()
	defer teardown()

	request := MKaaSPoolUpdateRequest{
		Name:      PtrTo("updated-pool"),
		NodeCount: PtrTo(4),
	}

	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(MKaaSClustersBasePathV2, strconv.Itoa(projectID), strconv.Itoa(regionID), strconv.Itoa(testResourceIntID), "pools", strconv.Itoa(testResourceIntID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		reqBody := &MKaaSPoolUpdateRequest{}
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		// Compare selected fields because request contains pointer fields
		if reqBody.Name == nil || *reqBody.Name != *request.Name {
			t.Errorf("expected name %v, got %v", *request.Name, reqBody.Name)
		}
		if reqBody.NodeCount == nil || *reqBody.NodeCount != *request.NodeCount {
			t.Errorf("expected node count %v, got %v", *request.NodeCount, reqBody.NodeCount)
		}
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.MkaaS.PoolUpdate(ctx, testResourceIntID, testResourceIntID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}
