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

func TestServerGroupGetByInstance(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	sgs := []edgecloud.ServerGroup{
		{
			ID:        testResourceID,
			Instances: []edgecloud.ServerGroupInstance{{InstanceID: testResourceID}},
		},
	}
	URL := path.Join("/v1/servergroups", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(sgs)
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

	sg, err := ServerGroupGetByInstance(context.Background(), client, testResourceID)
	assert.NoError(t, err)
	assert.Equal(t, testResourceID, sg.ID)
}

func TestServerGroupGetByInstance_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	URL := path.Join("/v1/servergroups", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	sg, err := ServerGroupGetByInstance(context.Background(), client, testResourceID)
	assert.Error(t, err)
	assert.Nil(t, sg)
}

func TestServerGroupGetByInstance_CustomError(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	URL := path.Join("/v1/servergroups", strconv.Itoa(projectID), strconv.Itoa(regionID))

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

	sgs, err := ServerGroupGetByInstance(context.Background(), client, testResourceID)
	assert.ErrorIs(t, err, ErrServerGroupNotFound)
	assert.Nil(t, sgs)
}
