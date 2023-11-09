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
