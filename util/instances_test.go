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
	"github.com/stretchr/testify/require"

	edgecloud "github.com/Edge-Center/edgecentercloud-go/v2"
)

func TestWaitForInstanceShutoff(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	expectedRespStop := edgecloud.Instance{ID: testResourceID}
	URLStop := path.Join("/v1/instances", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, "stop")
	mux.HandleFunc(URLStop, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(expectedRespStop)
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	expectedRespGet := edgecloud.Instance{ID: testResourceID, Status: InstanceShutoffStatus}
	URLGet := path.Join("/v1/instances", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)
	mux.HandleFunc(URLGet, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(expectedRespGet)
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	err := WaitForInstanceShutoff(context.Background(), client, testResourceID, nil)
	assert.NoError(t, err)
}

func TestWaitForInstanceShutoff_InstanceGet_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	expectedRespStop := edgecloud.Instance{ID: testResourceID}
	URLStop := path.Join("/v1/instances", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, "stop")
	mux.HandleFunc(URLStop, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(expectedRespStop)
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	URL := path.Join("/v1/instances", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)
	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	err := WaitForInstanceShutoff(context.Background(), client, testResourceID, &attempts)
	assert.Error(t, err)
}

func TestWaitForInstanceShutoff_InstanceStop_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	URL := path.Join("/v1/instances", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, "stop")
	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	err := WaitForInstanceShutoff(context.Background(), client, testResourceID, &attempts)
	assert.Error(t, err)
}

func TestWaitForInstanceShutoff_ErrInstanceNotShutOff(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	expectedRespStop := edgecloud.Instance{ID: testResourceID}
	URLStop := path.Join("/v1/instances", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, "stop")
	mux.HandleFunc(URLStop, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(expectedRespStop)
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	expectedRespGet := edgecloud.Instance{ID: testResourceID}
	URLGet := path.Join("/v1/instances", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)
	mux.HandleFunc(URLGet, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(expectedRespGet)
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	err := WaitForInstanceShutoff(context.Background(), client, testResourceID, &attempts)
	assert.ErrorIs(t, err, ErrInstanceNotShutOff)
}

func TestInstanceNetworkInterfaceByID(t *testing.T) {
	ctx := context.Background()
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	instanceID := testResourceID

	instanceIfaces := []edgecloud.InstanceInterface{{PortID: testResourceID}, {PortID: testResourceID2}}
	URL := path.Join("/v1/instances/", strconv.Itoa(projectID), strconv.Itoa(regionID), instanceID, "interfaces")

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(instanceIfaces)
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

	iface, err := InstanceNetworkInterfaceByID(ctx, client, instanceID, testResourceID2)
	require.NoError(t, err)
	assert.Equal(t, testResourceID2, iface.PortID)
}
