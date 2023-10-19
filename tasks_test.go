package edgecloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTasks_Get(t *testing.T) {
	setup()
	defer teardown()

	const (
		taskID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	task := &Task{ID: taskID}
	getVolumesURL := fmt.Sprintf("/v1/tasks/%s", taskID)

	mux.HandleFunc(getVolumesURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(task)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Tasks.Get(ctx, taskID)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, task) {
		t.Errorf("Tasks.Get\n returned %+v,\n expected %+v", resp, task)
	}
}
