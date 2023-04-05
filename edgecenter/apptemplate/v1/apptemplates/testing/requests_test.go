package testing

import (
	"fmt"
	"net/http"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/Edge-Center/edgecentercloud-go/edgecenter/apptemplate/v1/apptemplates"
	"github.com/Edge-Center/edgecentercloud-go/pagination"
	th "github.com/Edge-Center/edgecentercloud-go/testhelper"
	fake "github.com/Edge-Center/edgecentercloud-go/testhelper/client"
)

func prepareListTestURL() string {
	return fmt.Sprintf("/v1/apptemplates/%d/%d", fake.ProjectID, fake.RegionID)
}

func prepareGetTestURL(id string) string {
	return fmt.Sprintf("/v1/apptemplates/%d/%d/%s", fake.ProjectID, fake.RegionID, id)
}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ListResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("apptemplates", "v1")
	count := 0

	err := apptemplates.List(client).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := apptemplates.ExtractAppTemplates(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, AppTemplate1, ct)
		require.Equal(t, ExpectedAppTemplateSlice, actual)
		return true, nil
	})

	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestListAll(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, ListResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("apptemplates", "v1")

	groups, err := apptemplates.ListAll(client)
	require.NoError(t, err)
	ct := groups[0]
	require.Equal(t, AppTemplate1, ct)
	require.Equal(t, ExpectedAppTemplateSlice, groups)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(AppTemplate1.ID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, GetResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("apptemplates", "v1")

	ct, err := apptemplates.Get(client, AppTemplate1.ID).Extract()

	require.NoError(t, err)
	require.Equal(t, AppTemplate1, *ct)
}
