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

var attempts uint = 2

func TestVolumesListByName(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	volumes := []edgecloud.Volume{
		{
			Name: testResourceID,
		},
		{
			Name: testResourceID,
		},
	}
	URL := path.Join("/v1/volumes", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(volumes)
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

	volumes, err := VolumesListByName(context.Background(), client, testResourceID)
	assert.NoError(t, err)
	assert.Len(t, volumes, 2)
}

func TestVolumesListByName_VolumeList_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	URL := path.Join("/v1/volumes", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	volumes, err := VolumesListByName(context.Background(), client, testResourceID)
	assert.Error(t, err)
	assert.Nil(t, volumes)
}

func TestVolumesListByName_VolumesNotFound_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	URL := path.Join("/v1/volumes", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(nil)
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

	volumes, err := VolumesListByName(context.Background(), client, testResourceID)
	assert.ErrorIs(t, err, ErrVolumesNotFound)
	assert.Nil(t, volumes)
}

func TestWaitVolumeAttachedToInstance(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	expectedResp := edgecloud.Volume{
		ID:          testResourceID,
		Attachments: []edgecloud.Attachment{{ServerID: testResourceID}},
	}
	URL := path.Join("/v1/volumes", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

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

	err := WaitVolumeAttachedToInstance(context.Background(), client, testResourceID, testResourceID, nil)
	assert.NoError(t, err)
}

func TestWaitVolumeAttachedToInstance_VolumeGet_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	URL := path.Join("/v1/volumes", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)
	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	err := WaitVolumeAttachedToInstance(context.Background(), client, testResourceID, testResourceID, &attempts)
	assert.Error(t, err)
}

func TestWaitVolumeAttachedToInstance_VolumesNotAttached_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	expectedResp := edgecloud.Volume{ID: testResourceID}
	URL := path.Join("/v1/volumes", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

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

	err := WaitVolumeAttachedToInstance(context.Background(), client, testResourceID, testResourceID, &attempts)
	assert.ErrorIs(t, err, ErrVolumesNotAttached)
}

func TestWaitVolumeDetachedFromInstance(t *testing.T) {
	tests := []struct {
		name         string
		expectedResp edgecloud.Volume
	}{
		{
			name:         "empty Attachments",
			expectedResp: edgecloud.Volume{ID: testResourceID},
		},
		{
			name: "no Attachments",
			expectedResp: edgecloud.Volume{
				ID:          testResourceID,
				Attachments: []edgecloud.Attachment{{ServerID: "123"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			server := httptest.NewServer(mux)
			defer server.Close()

			URL := path.Join("/v1/volumes", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)
			mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
				resp, err := json.Marshal(tt.expectedResp)
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

			err := WaitVolumeDetachedFromInstance(context.Background(), client, testResourceID, testResourceID, nil)
			assert.NoError(t, err)
		})
	}
}

func TestWaitVolumeDetachedFromInstance_VolumeGet_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	URL := path.Join("/v1/volumes", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	err := WaitVolumeDetachedFromInstance(context.Background(), client, testResourceID, testResourceID, &attempts)
	assert.Error(t, err)
}

func TestWaitVolumeDetachedFromInstance_VolumesNotDetachedError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	expectedResp := edgecloud.Volume{
		ID:          testResourceID,
		Attachments: []edgecloud.Attachment{{ServerID: testResourceID}},
	}
	URL := path.Join("/v1/volumes", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)
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

	err := WaitVolumeDetachedFromInstance(context.Background(), client, testResourceID, testResourceID, &attempts)
	assert.ErrorIs(t, err, ErrVolumesNotDetached)
}
