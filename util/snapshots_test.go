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

func TestSnapshotsListByStatusAndVolumeID(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	expectedResp := []edgecloud.Snapshot{{Status: "new"}}
	URL := path.Join("/v1/snapshots", strconv.Itoa(projectID), strconv.Itoa(regionID))
	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(expectedResp)
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

	snapshots, err := SnapshotsListByStatusAndVolumeID(context.Background(), client, "new", testResourceID)
	assert.NoError(t, err)
	assert.Len(t, snapshots, 1)
}

func TestSnapshotsListByStatusAndVolumeID_SnapshotsNotFound_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	expectedResp := []edgecloud.Snapshot{{Status: "existing"}}
	URL := path.Join("/v1/snapshots", strconv.Itoa(projectID), strconv.Itoa(regionID))
	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(expectedResp)
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

	snapshots, err := SnapshotsListByStatusAndVolumeID(context.Background(), client, "new", testResourceID)
	assert.ErrorIs(t, err, ErrSnapshotsNotFound)
	assert.Nil(t, snapshots)
}

func TestSnapshotsListByStatusAndVolumeID_SnapshotsList_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	URL := path.Join("/v1/snapshots", strconv.Itoa(projectID), strconv.Itoa(regionID))
	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	snapshots, err := SnapshotsListByStatusAndVolumeID(context.Background(), client, "new", testResourceID)
	assert.Error(t, err)
	assert.Nil(t, snapshots)
}

func TestSnapshotsListByNameAndVolumeID(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	expectedResp := []edgecloud.Snapshot{{Name: testName}}
	URL := path.Join("/v1/snapshots", strconv.Itoa(projectID), strconv.Itoa(regionID))
	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(expectedResp)
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

	snapshots, err := SnapshotsListByNameAndVolumeID(context.Background(), client, testName, testResourceID)
	assert.NoError(t, err)
	assert.Len(t, snapshots, 1)
}

func TestSnapshotsListByNameAndVolumeID_SnapshotsNotFound_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	expectedResp := []edgecloud.Snapshot{{Name: "other_name"}}
	URL := path.Join("/v1/snapshots", strconv.Itoa(projectID), strconv.Itoa(regionID))
	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(expectedResp)
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

	snapshots, err := SnapshotsListByNameAndVolumeID(context.Background(), client, testName, testResourceID)
	assert.ErrorIs(t, err, ErrSnapshotsNotFound)
	assert.Nil(t, snapshots)
}

func TestSnapshotsListByNameAndVolumeID_SnapshotsList_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	URL := path.Join("/v1/snapshots", strconv.Itoa(projectID), strconv.Itoa(regionID))
	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	snapshots, err := SnapshotsListByNameAndVolumeID(context.Background(), client, testName, testResourceID)
	assert.Error(t, err)
	assert.Nil(t, snapshots)
}

func TestWaitSnapshotStatusReady(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	expectedResp := edgecloud.Snapshot{ID: testResourceID, Status: SnapshotReadyStatus}
	URL := path.Join("/v1/snapshots", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)
	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(expectedResp)
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

	err := WaitSnapshotStatusReady(context.Background(), client, testResourceID, nil)
	assert.NoError(t, err)
}

func TestWaitSnapshotStatusReady_SnapshotsGet_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	URL := path.Join("/v1/snapshots", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)
	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	err := WaitSnapshotStatusReady(context.Background(), client, testResourceID, &attempts)
	assert.Error(t, err)
}

func TestWaitSnapshotStatusReady_SnapshotNotReadyError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	expectedResp := edgecloud.Snapshot{ID: testResourceID, Status: "not ready"}
	URL := path.Join("/v1/snapshots", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)
	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(expectedResp)
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

	err := WaitSnapshotStatusReady(context.Background(), client, testResourceID, &attempts)
	assert.ErrorIs(t, err, ErrSnapshotNotReady)
}
