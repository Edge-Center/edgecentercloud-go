package util

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

func TestWaitForTaskCompleteWithTestServer(t *testing.T) {
	const (
		taskID  = "f0d19cec-5c3f-4853-886e-304915960ff6"
		timeout = 6 * time.Second
	)

	tests := []struct {
		name          string
		expectedError error
		serverHandler func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name: "Task finishes successfully",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				resp, err := json.Marshal(&edgecloud.Task{ID: taskID, State: edgecloud.TaskStateFinished})
				if err != nil {
					t.Fatalf("failed to marshal JSON: %v", err)
				}
				_, _ = fmt.Fprintf(w, `{"task":%s}`, string(resp))
			},
			expectedError: nil,
		},
		{
			name: "Task state is unknown",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				resp, err := json.Marshal(&edgecloud.Task{ID: taskID, State: "UnknownState"})
				if err != nil {
					t.Fatalf("failed to marshal JSON: %v", err)
				}
				_, _ = fmt.Fprintf(w, `{"task":%s}`, string(resp))
			},
			expectedError: fmt.Errorf("%w: [%s]", errTaskStateUnknown, "UnknownState"),
		},
		{
			name: "Task with error state",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				resp, err := json.Marshal(&edgecloud.Task{ID: taskID, State: edgecloud.TaskStateError})
				if err != nil {
					t.Fatalf("failed to marshal JSON: %v", err)
				}
				_, _ = fmt.Fprintf(w, `{"task":%s}`, string(resp))
			},
			expectedError: errTaskWithErrorState,
		},
		{
			name: "Task timeout exceeded",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				resp, err := json.Marshal(&edgecloud.Task{ID: taskID, State: edgecloud.TaskStateRunning})
				if err != nil {
					t.Fatalf("failed to marshal JSON: %v", err)
				}
				_, _ = fmt.Fprintf(w, `{"task":%s}`, string(resp))
			},
			expectedError: errTaskWaitTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			server := httptest.NewServer(mux)
			defer server.Close()

			mux.HandleFunc(fmt.Sprintf("/v1/tasks/%s", taskID), tt.serverHandler)

			client := edgecloud.NewClient(nil)
			baseURL, _ := url.Parse(server.URL)
			client.BaseURL = baseURL

			err := WaitForTaskComplete(context.Background(), client, taskID, timeout)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
