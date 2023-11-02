package edgecloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjects_Get(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &Project{
		ID:   1,
		Name: "test-project",
	}
	URL := path.Join(projectsBasePath, "1")

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Projects.Get(ctx, "1")
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestProjects_Delete(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(projectsBasePath, "1")

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Projects.Delete(ctx, "1")
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestProjects_Update(t *testing.T) {
	setup()
	defer teardown()

	request := &ProjectUpdateRequest{Name: "new-project"}
	expectedResp := &Project{Name: "new-project"}
	URL := path.Join(projectsBasePath, "1")

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		reqBody := new(ProjectUpdateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, reqBody)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Projects.Update(ctx, "1", request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestProjects_List(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []Project{{Name: "test-projects"}}

	mux.HandleFunc(projectsBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Projects.List(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestProjects_Create(t *testing.T) {
	setup()
	defer teardown()

	request := &ProjectCreateRequest{Name: "new-project"}
	expectedResp := &Project{Name: "new-project"}

	mux.HandleFunc(projectsBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(ProjectCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, reqBody)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Projects.Create(ctx, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}
