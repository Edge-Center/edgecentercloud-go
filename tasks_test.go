package edgecloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTasks_ListActive(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []Task{{ID: taskID}}
	URL := path.Join(tasksBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), tasksActive)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Tasks.ListActive(ctx)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestTasks_Acknowledge(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &Task{ID: taskID}
	URL := path.Join(tasksBasePathV1, taskID, tasksAcknowledge)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Tasks.Acknowledge(ctx, taskID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestTasks_AcknowledgeAll(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(tasksBasePathV1, tasksAcknowledge)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
	})

	resp, err := client.Tasks.AcknowledgeAll(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestTasks_Get(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &Task{ID: taskID}
	URL := path.Join(tasksBasePathV1, taskID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Tasks.Get(ctx, taskID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestTasks_List(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []Task{{ID: taskID}}

	mux.HandleFunc(tasksBasePathV1, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Tasks.List(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}
