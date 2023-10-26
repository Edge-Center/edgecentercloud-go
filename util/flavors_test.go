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

func TestLoadbalancerFlavorIsExist(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	const (
		projectID = 27520
		regionID  = 8
	)

	flavor := "g1-gpu-1-2-1"
	flavors := []edgecloud.Flavor{{
		FlavorName: flavor,
		FlavorID:   flavor,
	}}

	flavorsListURL := fmt.Sprintf("/v1/lbflavors/%d/%d", projectID, regionID)
	mux.HandleFunc(flavorsListURL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(flavors)
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

	exist, err := LoadbalancerFlavorIsExist(context.Background(), client, flavor)
	assert.NoError(t, err)
	assert.Equal(t, exist, true)
}

func TestFlavorIsExist(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	const (
		projectID = 27520
		regionID  = 8
	)

	flavor := "g1-gpu-1-2-1"
	flavors := []edgecloud.Flavor{{
		FlavorName: flavor,
		FlavorID:   flavor,
	}}

	flavorsListURL := fmt.Sprintf("/v1/flavors/%d/%d", projectID, regionID)
	mux.HandleFunc(flavorsListURL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(flavors)
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

	exist, err := FlavorIsExist(context.Background(), client, flavor)
	assert.NoError(t, err)
	assert.Equal(t, exist, true)
}
