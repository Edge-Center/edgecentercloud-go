package edgecloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestProjects_Get(t *testing.T) {
	setup()
	defer teardown()

	project := &Project{
		ID:   1,
		Name: "test-project",
	}

	mux.HandleFunc("/v1/projects/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(project)
		_, _ = fmt.Fprintf(w, `{"project":%s}`, string(resp))
	})

	resp, _, err := client.Projects.Get(ctx, "1")
	if err != nil {
		t.Errorf("Projects.Get returned error: %v", err)
	}

	if !reflect.DeepEqual(resp, project) {
		t.Errorf("Projects.Get\n returned %+v,\n expected %+v", resp, project)
	}
}
