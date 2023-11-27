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
