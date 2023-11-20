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

const (
	testResourceID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	testName       = "test-name"
	projectID      = 2750
	regionID       = 8
)

func TestResourceIsDeleted(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		expected   error
	}{
		{
			name:       "OK resource not found",
			statusCode: http.StatusNotFound,
			expected:   nil,
		},
		{
			name:       "resource not deleted",
			statusCode: http.StatusOK,
			expected:   errResourceNotDeleted,
		},
		{
			name:       "error retrieving resource",
			statusCode: http.StatusInternalServerError,
			expected:   errGetResourceInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			server := httptest.NewServer(mux)
			defer server.Close()

			URL := path.Join("/v1/floatingips", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)
			mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				resp, _ := json.MarshalIndent(edgecloud.Response{}, "", "    ")
				_, _ = fmt.Fprint(w, string(resp))
			})

			client := edgecloud.NewClient(nil)
			baseURL, _ := url.Parse(server.URL)
			client.BaseURL = baseURL
			client.Project = projectID
			client.Region = regionID

			GetResourceFunc := func(ctx context.Context, id string) (*edgecloud.FloatingIP, *edgecloud.Response, error) {
				return client.Floatingips.Get(ctx, id)
			}

			err := ResourceIsDeleted(context.Background(), GetResourceFunc, testResourceID)
			assert.Equal(t, tt.expected, err)
		})
	}
}

func TestResourceIsExist(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		expectedError error
		exist         bool
	}{
		{
			name:          "exist",
			statusCode:    http.StatusOK,
			expectedError: nil,
			exist:         true,
		},
		{
			name:          "not exist",
			statusCode:    http.StatusNotFound,
			expectedError: nil,
			exist:         false,
		},
		{
			name:          "default case error",
			statusCode:    http.StatusInternalServerError,
			expectedError: fmt.Errorf("%w, status code: %d, details: %s", errGetResourceInfo, 500, "500"),
			exist:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			server := httptest.NewServer(mux)
			defer server.Close()

			URL := path.Join("/v1/networks", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)
			mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				resp, _ := json.MarshalIndent(edgecloud.Response{}, "", "    ")
				_, _ = fmt.Fprint(w, string(resp))
			})

			client := edgecloud.NewClient(nil)
			baseURL, _ := url.Parse(server.URL)
			client.BaseURL = baseURL
			client.Project = projectID
			client.Region = regionID

			getFunc := func(ctx context.Context, id string) (*edgecloud.Network, *edgecloud.Response, error) {
				return client.Networks.Get(ctx, id)
			}

			exist, err := ResourceIsExist(context.Background(), getFunc, testResourceID)
			if tt.name == "default case error" {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expectedError, err)
			}
			assert.Equal(t, tt.exist, exist)
		})
	}
}

func TestDeleteResourceIfExist(t *testing.T) {
	client := edgecloud.NewClient(nil)

	tests := []struct {
		name     string
		urlPath  string
		resource interface{}
	}{
		{
			name:     "LoadbalancersService",
			urlPath:  "/v1/loadbalancers",
			resource: client.Loadbalancers,
		},
		{
			name:     "FloatingIPsService",
			urlPath:  "/v1/floatingips",
			resource: client.Floatingips,
		},
		{
			name:     "VolumesService",
			urlPath:  "/v1/volumes",
			resource: client.Volumes,
		},
		{
			name:     "L7PoliciesService",
			urlPath:  "/v1/l7policies",
			resource: client.L7Policies,
		},
		{
			name:     "SnapshotsService",
			urlPath:  "/v1/snapshots",
			resource: client.Snapshots,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			server := httptest.NewServer(mux)
			defer server.Close()

			baseURL, _ := url.Parse(server.URL)
			client.BaseURL = baseURL
			client.Project = projectID
			client.Region = regionID

			URL := path.Join(tt.urlPath, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)
			mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodDelete:
					resp, _ := json.Marshal(&edgecloud.TaskResponse{Tasks: []string{testResourceID}})
					_, _ = fmt.Fprint(w, string(resp))
				case http.MethodGet:
					w.WriteHeader(http.StatusNotFound)
					resp, _ := json.MarshalIndent(edgecloud.Response{}, "", "    ")
					_, _ = fmt.Fprint(w, string(resp))
				}
			})

			mux.HandleFunc(path.Join("/v1/tasks", testResourceID), func(w http.ResponseWriter, r *http.Request) {
				resp, err := json.Marshal(&edgecloud.Task{ID: testResourceID, State: edgecloud.TaskStateFinished})
				if err != nil {
					t.Fatalf("failed to marshal JSON: %v", err)
				}
				_, _ = fmt.Fprint(w, string(resp))
			})

			err := DeleteResourceIfExist(context.Background(), client, tt.resource, testResourceID)
			assert.NoError(t, err)
		})
	}
}

func TestDeleteResourceIfExist_deleteAndWait_Error(t *testing.T) {
	client := edgecloud.NewClient(nil)

	tests := []struct {
		name     string
		urlPath  string
		resource interface{}
	}{
		{
			name:     "LoadbalancersService",
			urlPath:  "/v1/loadbalancers",
			resource: client.Loadbalancers,
		},
		{
			name:     "FloatingIPsService",
			urlPath:  "/v1/floatingips",
			resource: client.Floatingips,
		},
		{
			name:     "VolumesService",
			urlPath:  "/v1/volumes",
			resource: client.Volumes,
		},
		{
			name:     "L7PoliciesService",
			urlPath:  "/v1/l7policies",
			resource: client.L7Policies,
		},
		{
			name:     "SnapshotsService",
			urlPath:  "/v1/snapshots",
			resource: client.Snapshots,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			server := httptest.NewServer(mux)
			defer server.Close()

			baseURL, _ := url.Parse(server.URL)
			client.BaseURL = baseURL
			client.Project = projectID
			client.Region = regionID

			URL := path.Join(tt.urlPath, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)
			mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			mux.HandleFunc(path.Join("/v1/tasks", testResourceID), func(w http.ResponseWriter, r *http.Request) {
				resp, err := json.Marshal(&edgecloud.Task{ID: testResourceID, State: edgecloud.TaskStateFinished})
				if err != nil {
					t.Fatalf("failed to marshal JSON: %v", err)
				}
				_, _ = fmt.Fprint(w, string(resp))
			})

			err := DeleteResourceIfExist(context.Background(), client, tt.resource, testResourceID)
			assert.Error(t, err)
		})
	}
}

func TestDeleteResourceIfExist_ResourceNotSupported(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := edgecloud.NewClient(nil)

	err := DeleteResourceIfExist(context.Background(), nil, client.Flavors, testResourceID)
	assert.Equal(t, err, errDeleteResourceIfExistIsNotSupported)
}
