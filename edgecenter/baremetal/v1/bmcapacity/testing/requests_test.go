package testing

import (
	"fmt"
	"net/http"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/Edge-Center/edgecentercloud-go/edgecenter/baremetal/v1/bmcapacity"
	th "github.com/Edge-Center/edgecentercloud-go/testhelper"
	fake "github.com/Edge-Center/edgecentercloud-go/testhelper/client"
)

func prepareGetAvailableNodesTestURL() string {
	return fmt.Sprintf("/v1/bmcapacity/%d/%d", fake.ProjectID, fake.RegionID)
}

func TestGetCapacity(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareGetAvailableNodesTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, availableNodesResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("bmcapacity", "v1")
	nodes, err := bmcapacity.GetAvailableNodes(client).Extract()

	require.NoError(t, err)
	require.Equal(t, availableNodes, *nodes)
}
