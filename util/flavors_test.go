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
	testFlavor = "g1-gpu-1-2-1"
)

func TestLoadbalancerFlavorIsExist(t *testing.T) {
	tests := []struct {
		name       string
		flavors    []edgecloud.Flavor
		flavorName string
		want       bool
	}{
		{
			name:       "flavor exists",
			flavors:    []edgecloud.Flavor{{FlavorName: testFlavor, FlavorID: testFlavor}},
			flavorName: testFlavor,
			want:       true,
		},
		{
			name:       "flavor does not exist",
			flavors:    []edgecloud.Flavor{{FlavorName: testFlavor, FlavorID: testFlavor}},
			flavorName: "non-existent-flavor",
			want:       false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mux := http.NewServeMux()
			server := httptest.NewServer(mux)
			defer server.Close()

			URL := path.Join("/v1/lbflavors", strconv.Itoa(projectID), strconv.Itoa(regionID))

			mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
				resp, err := json.Marshal(tc.flavors)
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

			exist, err := LoadbalancerFlavorIsExist(context.Background(), client, tc.flavorName)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, exist)
		})
	}
}

func TestLoadbalancerFlavorIsExist_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	URL := path.Join("/v1/lbflavors", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	exist, err := LoadbalancerFlavorIsExist(context.Background(), client, testFlavor)
	assert.Error(t, err)
	assert.Equal(t, false, exist)
}

func TestFlavorIsExist(t *testing.T) {
	tests := []struct {
		name       string
		flavors    []edgecloud.Flavor
		flavorName string
		want       bool
	}{
		{
			name:       "flavor exists",
			flavors:    []edgecloud.Flavor{{FlavorName: testFlavor, FlavorID: testFlavor}},
			flavorName: testFlavor,
			want:       true,
		},
		{
			name:       "flavor does not exist",
			flavors:    []edgecloud.Flavor{{FlavorName: testFlavor, FlavorID: testFlavor}},
			flavorName: "non-existent-flavor",
			want:       false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mux := http.NewServeMux()
			server := httptest.NewServer(mux)
			defer server.Close()

			URL := path.Join("/v1/flavors", strconv.Itoa(projectID), strconv.Itoa(regionID))

			mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
				resp, err := json.Marshal(tc.flavors)
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

			exist, err := FlavorIsExist(context.Background(), client, tc.flavorName)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, exist)
		})
	}
}

func TestFlavorIsExist_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	URL := path.Join("/v1/flavors", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	exist, err := FlavorIsExist(context.Background(), client, testFlavor)
	assert.Error(t, err)
	assert.Equal(t, exist, false)
}

func TestFlavorIsAvailable(t *testing.T) {
	tests := []struct {
		name       string
		flavors    []edgecloud.Flavor
		flavorName string
		want       bool
	}{
		{
			name:       "flavor is available",
			flavors:    []edgecloud.Flavor{{FlavorName: testFlavor, FlavorID: testFlavor}},
			flavorName: testFlavor,
			want:       true,
		},
		{
			name:       "flavor does not available",
			flavors:    []edgecloud.Flavor{{FlavorName: testFlavor, FlavorID: testFlavor}},
			flavorName: "not-available-flavor",
			want:       false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mux := http.NewServeMux()
			server := httptest.NewServer(mux)
			defer server.Close()

			instanceCheckFlavorVolumeRequest := &edgecloud.InstanceCheckFlavorVolumeRequest{
				Volumes: []edgecloud.InstanceVolumeCreate{{Source: edgecloud.VolumeSourceExistingVolume}},
			}
			URL := path.Join("/v1/instances", strconv.Itoa(projectID), strconv.Itoa(regionID), "available_flavors")

			mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
				resp, err := json.Marshal(tc.flavors)
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

			available, err := FlavorIsAvailable(context.Background(), client, tc.flavorName, instanceCheckFlavorVolumeRequest)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, available)
		})
	}
}

func TestFlavorIsAvailable_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	instanceCheckFlavorVolumeRequest := &edgecloud.InstanceCheckFlavorVolumeRequest{
		Volumes: []edgecloud.InstanceVolumeCreate{{Source: edgecloud.VolumeSourceExistingVolume}},
	}
	URL := path.Join("/v1/instances", strconv.Itoa(projectID), strconv.Itoa(regionID), "available_flavors")

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	available, err := FlavorIsAvailable(context.Background(), client, testFlavor, instanceCheckFlavorVolumeRequest)
	assert.Error(t, err)
	assert.Equal(t, false, available)
}
