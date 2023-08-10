package testing

import (
	"fmt"
	"net/http"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/Edge-Center/edgecentercloud-go/edgecenter/quota/v2/quotas"
	th "github.com/Edge-Center/edgecentercloud-go/testhelper"
	fake "github.com/Edge-Center/edgecentercloud-go/testhelper/client"
)

func prepareListCombinedTestURL() string {
	return "/v2/quotas_client"
}

func prepareListGlobalTestURL(clientID int) string {
	return fmt.Sprintf("/v2/quotas_global/%d", clientID)
}

func prepareListRegionalTestURL(clientID, regionID int) string {
	return fmt.Sprintf("/v2/quotas_regional/%d/%d", clientID, regionID)
}

func TestListCombined(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareListCombinedTestURL()

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, CombinedResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("quotas", "v2")

	ct, err := quotas.ListCombined(client, quotas.ListCombinedOpts{}).Extract()
	require.NoError(t, err)
	require.Equal(t, CombinedQuota1, *ct)
}

func TestListGlobal(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareListGlobalTestURL(clientID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GlobalResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("quotas", "v2")

	ct, err := quotas.ListGlobal(client, clientID).Extract()
	require.NoError(t, err)
	require.Equal(t, CombinedQuota1.GlobalQuotas, *ct)
}

func TestListRegional(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareListRegionalTestURL(clientID, regionID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, RegionalResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("quotas", "v2")

	ct, err := quotas.ListRegional(client, clientID, regionID).Extract()
	require.NoError(t, err)
	require.Equal(t, CombinedQuota1.RegionalQuotas[0], *ct)
}
