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

func TestL7PoliciesListByListenerID(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	l7Policies := []edgecloud.L7Policy{
		{
			ListenerID: testResourceID,
		},
		{
			ListenerID: testResourceID,
		},
	}
	URL := path.Join("/v1/l7policies", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(l7Policies)
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

	l7Policies, err := L7PoliciesListByListenerID(context.Background(), client, testResourceID)
	assert.NoError(t, err)
	assert.Len(t, l7Policies, 2)
}

func TestL7PoliciesListByListenerID_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	URL := path.Join("/v1/l7policies", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	l7Policies, err := L7PoliciesListByListenerID(context.Background(), client, testResourceID)
	assert.Error(t, err)
	assert.Nil(t, l7Policies)
}

func TestL7PoliciesListByListenerID_ErrL7PoliciesNotFound(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	var l7Policies []edgecloud.L7Policy
	URL := path.Join("/v1/l7policies", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(l7Policies)
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

	l7Policies, err := L7PoliciesListByListenerID(context.Background(), client, testResourceID)
	assert.ErrorIs(t, err, ErrL7PoliciesNotFound)
	assert.Nil(t, l7Policies)
}

func TestGetLbL7PolicyFromName_Ok(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	l7PolicyToFind := edgecloud.L7Policy{ID: testResourceID, Name: testName}
	l7PolicyOther := edgecloud.L7Policy{ID: testResourceID2, Name: "other_name"}

	l7Policies := []edgecloud.L7Policy{l7PolicyToFind, l7PolicyOther}
	URL := path.Join("/v1/l7policies", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(l7Policies)
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

	l7Policy, err := GetLbL7PolicyFromName(context.Background(), client, testName)
	assert.NoError(t, err)
	assert.Equal(t, l7PolicyToFind, *l7Policy)
}
