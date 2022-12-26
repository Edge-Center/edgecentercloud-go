package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Edge-Center/edgecentercloud-go/edgecenter/keystone/v1/types"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"

	"github.com/Edge-Center/edgecentercloud-go/edgecenter/keystone/v1/keystones"
	fake "github.com/Edge-Center/edgecentercloud-go/testhelper/client"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	"github.com/Edge-Center/edgecentercloud-go/pagination"
	th "github.com/Edge-Center/edgecentercloud-go/testhelper"
)

func prepareListTestURL() string {
	return "/v1/keystones"
}

func prepareGetTestURL(id int) string {
	return fmt.Sprintf("/v1/keystones/%d", id)
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

	client := fake.ServiceTokenClient("keystones", "v1")
	count := 0

	err := keystones.List(client).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := keystones.ExtractKeystones(page)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, Keystone1, ct)
		require.Equal(t, ExpectedKeystoneSlice, actual)
		return true, nil
	})

	require.NoError(t, err)

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

	client := fake.ServiceTokenClient("keystones", "v1")

	results, err := keystones.ListAll(client)
	require.NoError(t, err)
	ct := results[0]
	require.Equal(t, Keystone1, ct)
	require.Equal(t, ExpectedKeystoneSlice, results)

}

func TestGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(Keystone1.ID)

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

	client := fake.ServiceTokenClient("keystones", "v1")

	ct, err := keystones.Get(client, Keystone1.ID).Extract()

	require.NoError(t, err)
	require.Equal(t, Keystone1, *ct)
	require.Equal(t, createdTime, ct.CreatedOn)

}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc(prepareListTestURL(), func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		_, err := fmt.Fprint(w, CreateResponse)
		if err != nil {
			log.Error(err)
		}
	})

	options := keystones.CreateOpts{
		URL:                       Keystone1.URL,
		State:                     types.KeystoneStateNew,
		KeystoneFederatedDomainID: Keystone1.KeystoneFederatedDomainID,
		AdminPassword:             "",
	}

	err := edgecloud.TranslateValidationError(options.Validate())
	require.NoError(t, err)

	client := fake.ServiceTokenClient("keystones", "v1")
	keystone, err := keystones.Create(client, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Keystone1, *keystone)
	require.Equal(t, createdTime, keystone.CreatedOn)
}

func TestUpdate(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	testURL := prepareGetTestURL(Keystone1.ID)

	th.Mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(w, UpdateResponse)
		if err != nil {
			log.Error(err)
		}
	})

	client := fake.ServiceTokenClient("keystones", "v1")

	options := keystones.UpdateOpts{
		URL:   &Keystone1.URL,
		State: types.KeystoneStateDeleted,
	}

	err := edgecloud.TranslateValidationError(options.Validate())
	require.NoError(t, err)

	keystone, err := keystones.Update(client, Keystone1.ID, options).Extract()
	require.NoError(t, err)
	require.Equal(t, Keystone1, *keystone)
	require.Equal(t, createdTime, keystone.CreatedOn)

}
