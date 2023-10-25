package edgecloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegions_Get(t *testing.T) {
	setup()
	defer teardown()

	region := &Region{
		ID: 1,
	}

	mux.HandleFunc("/v1/regions/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(region)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Regions.Get(ctx, "1", nil)
	require.NoError(t, err)

	if !reflect.DeepEqual(resp, region) {
		t.Errorf("Regions.Get\n returned %+v,\n expected %+v", resp, region)
	}
}
