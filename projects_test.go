package edgecloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProjects_Get(t *testing.T) {
	setup()
	defer teardown()

	project := &Project{
		ID:   1,
		Name: "test-project",
	}
	URL := "/v1/projects/1"

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
