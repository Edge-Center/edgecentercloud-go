package util

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

func TestLoadbalancerGetByName(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	lbs := []edgecloud.Loadbalancer{
		{
			Name: testResourceID,
		},
	}
	URL := path.Join("/v1/loadbalancers", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(lbs)
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	lb, err := LoadbalancerGetByName(context.Background(), client, testResourceID)
	assert.NoError(t, err)
	assert.Equal(t, testResourceID, lb.Name)
}

func TestLoadbalancerGetByName_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	URL := path.Join("/v1/loadbalancers", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	lb, err := LoadbalancerGetByName(context.Background(), client, testResourceID)
	assert.Error(t, err)
	assert.Nil(t, lb)
}

func TestLoadbalancerGetByName_CustomErrors(t *testing.T) {
	tests := []struct {
		name      string
		lbs       []edgecloud.Loadbalancer
		expectErr error
	}{
		{
			name:      "ErrLoadbalancersNotFound",
			lbs:       nil,
			expectErr: ErrLoadbalancersNotFound,
		},
		{
			name: "ErrMultipleResults",
			lbs: []edgecloud.Loadbalancer{
				{
					Name: testName,
				},
				{
					Name: testName,
				},
			},
			expectErr: ErrMultipleResults,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mux := http.NewServeMux()
			server := httptest.NewServer(mux)
			defer server.Close()

			URL := path.Join("/v1/loadbalancers", strconv.Itoa(projectID), strconv.Itoa(regionID))

			mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
				resp, err := json.Marshal(tc.lbs)
				if err != nil {
					t.Fatalf("failed to marshal JSON: %v", err)
				}
				_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
			})

			client := edgecloud.NewClient(nil)
			baseURL, _ := url.Parse(server.URL)
			client.BaseURL = baseURL
			client.Project = projectID
			client.Region = regionID

			lbs, err := LoadbalancerGetByName(context.Background(), client, testName)
			assert.ErrorIs(t, err, tc.expectErr)
			assert.Nil(t, lbs)
		})
	}
}

func TestLBPoolGetByName(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	pools := []edgecloud.Pool{
		{
			ID:   testResourceID,
			Name: testName,
		},
	}
	URL := path.Join("/v1/lbpools", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(pools)
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	pool, err := LBPoolGetByName(context.Background(), client, testName, testResourceID)
	assert.NoError(t, err)
	assert.Equal(t, testName, pool.Name)
}

func TestLBPoolGetByName_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	URL := path.Join("/v1/lbpools", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	lb, err := LBPoolGetByName(context.Background(), client, testName, testResourceID)
	assert.Error(t, err)
	assert.Nil(t, lb)
}

func TestLBPoolGetByName_CustomErrors(t *testing.T) {
	tests := []struct {
		name      string
		pools     []edgecloud.Pool
		expectErr error
	}{
		{
			name:      "ErrLoadbalancerPoolsNotFound",
			pools:     nil,
			expectErr: ErrLoadbalancerPoolsNotFound,
		},
		{
			name: "ErrMultipleResults",
			pools: []edgecloud.Pool{
				{
					ID:   testResourceID,
					Name: testName,
				},
				{
					ID:   testResourceID,
					Name: testName,
				},
			},
			expectErr: ErrMultipleResults,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mux := http.NewServeMux()
			server := httptest.NewServer(mux)
			defer server.Close()

			URL := path.Join("/v1/lbpools", strconv.Itoa(projectID), strconv.Itoa(regionID))

			mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
				resp, err := json.Marshal(tc.pools)
				if err != nil {
					t.Fatalf("failed to marshal JSON: %v", err)
				}
				_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
			})

			client := edgecloud.NewClient(nil)
			baseURL, _ := url.Parse(server.URL)
			client.BaseURL = baseURL
			client.Project = projectID
			client.Region = regionID

			pools, err := LBPoolGetByName(context.Background(), client, testName, testResourceID)
			assert.ErrorIs(t, err, tc.expectErr)
			assert.Nil(t, pools)
		})
	}
}

func TestWaitLoadBalancerProvisioningStatusActive(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	expectedResp := edgecloud.Loadbalancer{
		Name:               testResourceID,
		ProvisioningStatus: edgecloud.ProvisioningStatusActive,
	}
	URL := path.Join("/v1/loadbalancers", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	err := WaitLoadBalancerProvisioningStatusActive(context.Background(), client, testResourceID)
	assert.NoError(t, err)
}
