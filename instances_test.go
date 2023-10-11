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

const (
	testID        = "f0d19cec-5c3f-4853-886e-304915960ff6"
	testProjectID = "27520"
	testRegionID  = "8"
)

func TestInstances_Get(t *testing.T) {
	setup()
	defer teardown()

	instance := &Instance{ID: testID}
	getInstanceURL := fmt.Sprintf("/v1/instances/%s/%s/%s", testProjectID, testRegionID, testID)

	mux.HandleFunc(getInstanceURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(instance)
		fmt.Fprintf(w, `{"instance":%s}`, string(resp))
	})

	opts := ServicePath{Project: testProjectID, Region: testRegionID}
	resp, _, err := client.Instances.Get(ctx, testID, &opts)
	if err != nil {
		t.Errorf("Instances.Get returned error: %v", err)
	}

	if !reflect.DeepEqual(resp, instance) {
		t.Errorf("Instances.Get\n returned %+v,\n expected %+v", resp, instance)
	}
}

func TestInstances_Create(t *testing.T) {
	setup()
	defer teardown()

	instanceCreateRequest := &InstanceCreateRequest{
		Names:          []string{"test-instance"},
		Flavor:         "g1-standard-1-2",
		Interfaces:     []InstanceInterface{{Type: ExternalInterfaceType}},
		SecurityGroups: []InstanceSecurityGroupsCreate{{ID: testID}},
		Volumes: []InstanceVolumeCreate{
			{
				TypeName:  SsdHiIops,
				Size:      5,
				BootIndex: 0,
				Source:    Image,
				ImageID:   testID,
			},
		},
	}

	taskResponse := &TaskResponse{Tasks: []string{testID}}

	createInstanceURL := fmt.Sprintf("/v2/instances/%s/%s", testProjectID, testRegionID)

	mux.HandleFunc(createInstanceURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(InstanceCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}
		assert.Equal(t, instanceCreateRequest, reqBody)
		resp, _ := json.Marshal(taskResponse)
		fmt.Fprintf(w, `{"tasks":%s}`, string(resp))
	})

	opts := ServicePath{Project: testProjectID, Region: testRegionID}
	resp, _, err := client.Instances.Create(ctx, instanceCreateRequest, &opts)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}

func TestInstances_Delete(t *testing.T) {
	setup()
	defer teardown()

	taskResponse := &TaskResponse{Tasks: []string{testID}}

	deleteInstanceURL := fmt.Sprintf("/v1/instances/%s/%s/%s", testProjectID, testRegionID, testID)

	mux.HandleFunc(deleteInstanceURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(taskResponse)
		fmt.Fprintf(w, `{"tasks":%s}`, string(resp))
	})

	opts := ServicePath{Project: testProjectID, Region: testRegionID}
	resp, _, err := client.Instances.Delete(ctx, testID, &opts)
	require.NoError(t, err)

	assert.Equal(t, taskResponse, resp)
}
