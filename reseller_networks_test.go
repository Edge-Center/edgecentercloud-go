package edgecloud

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResellerNetworks_List(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := ResellerNetworks{
		Count: 1,
		Results: []ResellerNetwork{
			{ID: testResourceID},
		},
	}

	mux.HandleFunc(resellerNetworksBasePathV1, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, err = w.Write(resp)
		if err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	})

	respActual, resp, err := client.ResellerNetworks.List(context.Background(), nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, *respActual, expectedResp)
}
