package util

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

func TestResourceIsDeleted(t *testing.T) {
	resourceID := "f0d19cec-5c3f-4853-886e-304915960ff6"
	projectID := 2750
	regionID := 8

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

			getResourceURL := fmt.Sprintf("/v1/floatingips/%d/%d/%s", projectID, regionID, resourceID)
			mux.HandleFunc(getResourceURL, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				resp, _ := json.MarshalIndent(edgecloud.Response{}, "", "    ")
				_, _ = fmt.Fprint(w, string(resp))
			})

			client := edgecloud.NewClient(nil)
			baseURL, _ := url.Parse(server.URL)
			client.BaseURL = baseURL
			client.Project = projectID
			client.Region = regionID

			retrieveFunc := func(ctx context.Context, id string) (*edgecloud.FloatingIP, *edgecloud.Response, error) {
				return client.Floatingips.Get(ctx, id)
			}

			err := ResourceIsDeleted(context.Background(), retrieveFunc, resourceID)
			assert.Equal(t, tt.expected, err)
		})
	}
}