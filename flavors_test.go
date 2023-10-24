package edgecloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFlavors_List(t *testing.T) {
	setup()
	defer teardown()

	const (
		flavorID = "f0d19cec-5c3f-4853-886e-304915960ff6"
	)

	flavors := []Flavor{{FlavorID: flavorID}}
	getFlavorsURL := fmt.Sprintf("/v1/flavors/%d/%d", projectID, regionID)

	mux.HandleFunc(getFlavorsURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(flavors)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.Flavors.List(ctx, nil)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, flavors) {
		t.Errorf("Flavors.List\n returned %+v,\n expected %+v", resp, flavors)
	}
}
