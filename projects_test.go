package edgecloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjects_Get(t *testing.T) {
	setup()
	defer teardown()

	project := &Project{
		ID:   1,
		Name: "test-project",
	}
	URL := fmt.Sprintf("%s/1", projectsBasePath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(project)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Projects.Get(ctx, "1")
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, project) {
		t.Errorf("Projects.Get\n returned %+v,\n expected %+v", resp, project)
	}
}

func TestProjects_Delete(t *testing.T) {
	setup()
	defer teardown()

	taskResponse := &TaskResponse{Tasks: []string{taskID}}
	URL := fmt.Sprintf("%s/1", projectsBasePath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Projects.Delete(ctx, "1")
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestProjects_Update(t *testing.T) {
	setup()
	defer teardown()

	updateRequest := &ProjectUpdateRequest{
		Name: "new-project",
	}
	projectResponse := &Project{Name: "new-project"}
	URL := fmt.Sprintf("%s/1", projectsBasePath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		reqBody := new(ProjectUpdateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, updateRequest, reqBody)
		resp, _ := json.Marshal(projectResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Projects.Update(ctx, "1", updateRequest)
	require.NoError(t, err)

	assert.Equal(t, projectResponse, resp)
}

func TestProjects_List(t *testing.T) {
	setup()
	defer teardown()

	projects := []Project{{Name: "test-projects"}}

	mux.HandleFunc(projectsBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(projects)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.Projects.List(ctx, nil)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, projects) {
		t.Errorf("Projects.List\n returned %+v,\n expected %+v", resp, projects)
	}
}

func TestProjects_Create(t *testing.T) {
	setup()
	defer teardown()

	createRequest := &ProjectCreateRequest{
		Name: "new-project",
	}
	projectResponse := &Project{Name: "new-project"}

	mux.HandleFunc(projectsBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(ProjectCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, createRequest, reqBody)
		resp, _ := json.Marshal(projectResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Projects.Create(ctx, createRequest)
	require.NoError(t, err)

	assert.Equal(t, projectResponse, resp)
}
