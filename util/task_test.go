package util

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

const (
	timeout = 6 * time.Second
)

func TestWaitTask(t *testing.T) {
	tests := []struct {
		name          string
		expectedError error
		serverHandler func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name: "Task finishes successfully",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				resp, err := json.Marshal(&edgecloud.Task{ID: testResourceID, State: edgecloud.TaskStateFinished})
				if err != nil {
					t.Fatalf("failed to marshal JSON: %v", err)
				}
				_, _ = fmt.Fprint(w, string(resp))
			},
			expectedError: nil,
		},
		{
			name: "Task state is unknown",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				resp, err := json.Marshal(&edgecloud.Task{ID: testResourceID, State: "UnknownState"})
				if err != nil {
					t.Fatalf("failed to marshal JSON: %v", err)
				}
				_, _ = fmt.Fprint(w, string(resp))
			},
			expectedError: fmt.Errorf("%w: [%s]", errTaskStateUnknown, "UnknownState"),
		},
		{
			name: "Task with error state",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				resp, err := json.Marshal(&edgecloud.Task{ID: testResourceID, State: edgecloud.TaskStateError})
				if err != nil {
					t.Fatalf("failed to marshal JSON: %v", err)
				}
				_, _ = fmt.Fprint(w, string(resp))
			},
			expectedError: edgecloud.NewArgError("taskID", errTaskWithErrorState.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			server := httptest.NewServer(mux)
			defer server.Close()

			URL := path.Join("/v1/tasks", testResourceID)
			mux.HandleFunc(URL, tt.serverHandler)

			client := edgecloud.NewClient(nil)
			baseURL, _ := url.Parse(server.URL)
			client.BaseURL = baseURL

			_, err := waitTask(context.Background(), client, testResourceID)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestWaitForTaskComplete_TaskTimeoutExceeded(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	URL := path.Join("/v1/tasks", testResourceID)
	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(&edgecloud.Task{ID: testResourceID, State: edgecloud.TaskStateRunning})
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL

	err := WaitForTaskComplete(context.Background(), client, testResourceID, timeout)
	assert.Equal(t, edgecloud.NewArgError("taskID", errTaskWaitTimeout.Error()), err)
}

func TestWaitAndGetTaskInfo_Success(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	URL := path.Join("/v1/tasks", testResourceID)
	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(&edgecloud.Task{ID: testResourceID, State: edgecloud.TaskStateFinished})
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL

	taskInfo, err := WaitAndGetTaskInfo(context.Background(), client, testResourceID, timeout)
	assert.NoError(t, err)
	assert.Equal(t, testResourceID, taskInfo.ID)
	assert.Equal(t, edgecloud.TaskStateFinished, taskInfo.State)
}

func TestExecuteAndExtractTaskResult_Success(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	URL := path.Join("/v1/tasks", testResourceID)
	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(&edgecloud.Task{ID: testResourceID, State: edgecloud.TaskStateFinished})
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL

	result, err := ExecuteAndExtractTaskResult(
		context.Background(),
		func(ctx context.Context, opt interface{}) (*edgecloud.TaskResponse, *edgecloud.Response, error) {
			return &edgecloud.TaskResponse{Tasks: []string{testResourceID}}, &edgecloud.Response{}, nil
		},
		"testOpt",
		client,
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}
