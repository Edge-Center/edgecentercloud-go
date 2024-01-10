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

	edgecloud "github.com/Edge-Center/edgecentercloud-go/v2"
)

func TestFloatingIPsListByPortID(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	floatingIPs := []edgecloud.FloatingIP{
		{
			PortID: testResourceID,
		},
		{
			PortID: testResourceID,
		},
	}
	URL := path.Join("/v1/floatingips", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(floatingIPs)
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

	floatingIPs, err := FloatingIPsListByPortID(context.Background(), client, testResourceID)
	assert.NoError(t, err)
	assert.Len(t, floatingIPs, 2)
}

func TestFloatingIPsListByPortID_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	URL := path.Join("/v1/floatingips", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	floatingIPs, err := FloatingIPsListByPortID(context.Background(), client, testResourceID)
	assert.Error(t, err)
	assert.Nil(t, floatingIPs)
}

func TestFloatingIPsListByPortID_ErrFloatingIPsNotFound(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	var floatingIPs []edgecloud.FloatingIP
	URL := path.Join("/v1/floatingips", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(floatingIPs)
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

	floatingIPs, err := FloatingIPsListByPortID(context.Background(), client, testResourceID)
	assert.ErrorIs(t, err, ErrFloatingIPsNotFound)
	assert.Nil(t, floatingIPs)
}

func TestFloatingIPDetailedByIPAddress(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	floatingIPs := []edgecloud.FloatingIP{{FloatingIPAddress: testResourceID}}
	URL := path.Join("/v1/floatingips", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(floatingIPs)
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

	floatingIP, err := FloatingIPDetailedByIPAddress(context.Background(), client, testResourceID)
	assert.NoError(t, err)
	assert.NotNil(t, floatingIP)
}

func TestFloatingIPDetailedByIPAddress_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	URL := path.Join("/v1/floatingips", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	floatingIP, err := FloatingIPDetailedByIPAddress(context.Background(), client, testResourceID)
	assert.Error(t, err)
	assert.Nil(t, floatingIP)
}

func TestFloatingIPDetailedByIPAddress_ErrFloatingIPNotFound(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	var floatingIPs []edgecloud.FloatingIP
	URL := path.Join("/v1/floatingips", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(floatingIPs)
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

	floatingIP, err := FloatingIPDetailedByIPAddress(context.Background(), client, testResourceID)
	assert.ErrorIs(t, err, ErrFloatingIPNotFound)
	assert.Nil(t, floatingIP)
}

func TestFloatingIPDetailedByID(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	floatingIPs := []edgecloud.FloatingIP{{ID: testResourceID}}
	URL := path.Join("/v1/floatingips", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(floatingIPs)
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

	floatingIP, err := FloatingIPDetailedByID(context.Background(), client, testResourceID)
	assert.NoError(t, err)
	assert.NotNil(t, floatingIP)
}

func TestFloatingIPDetailedByID_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	URL := path.Join("/v1/floatingips", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	floatingIP, err := FloatingIPDetailedByID(context.Background(), client, testResourceID)
	assert.Error(t, err)
	assert.Nil(t, floatingIP)
}

func TestFloatingIPDetailedByID_ErrFloatingIPNotFound(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	var floatingIPs []edgecloud.FloatingIP
	URL := path.Join("/v1/floatingips", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(floatingIPs)
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

	floatingIP, err := FloatingIPDetailedByID(context.Background(), client, testResourceID)
	assert.ErrorIs(t, err, ErrFloatingIPNotFound)
	assert.Nil(t, floatingIP)
}
