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

	flavors := []Flavor{{FlavorID: testResourceID}}
	URL := fmt.Sprintf("/v1/flavors/%d/%d", projectID, regionID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
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

func TestFlavors_ListBaremetal(t *testing.T) {
	setup()
	defer teardown()

	flavors := []Flavor{{FlavorID: testResourceID}}
	URL := fmt.Sprintf("/v1/bmflavors/%d/%d", projectID, regionID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(flavors)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.Flavors.ListBaremetal(ctx, nil)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, flavors) {
		t.Errorf("Flavors.ListBaremetal\n returned %+v,\n expected %+v", resp, flavors)
	}
}

func TestFlavors_ListBaremetalForClient(t *testing.T) {
	setup()
	defer teardown()

	flavors := []Flavor{{FlavorID: testResourceID}}
	URL := fmt.Sprintf("/v1/bmflavors/%d", regionID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(flavors)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	resp, _, err := client.Flavors.ListBaremetalForClient(ctx, nil)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, flavors) {
		t.Errorf("Flavors.ListBaremetalForClient\n returned %+v,\n expected %+v", resp, flavors)
	}
}
