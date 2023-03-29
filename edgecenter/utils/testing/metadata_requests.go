package testing

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/utils/metadata"
	th "github.com/Edge-Center/edgecentercloud-go/testhelper"
	fake "github.com/Edge-Center/edgecentercloud-go/testhelper/client"
)

type TestFunc func(t *testing.T)

type URLFunc func(c *edgecloud.ServiceClient, id string, args ...string)

func prepareTestParams(resourceName string, urlFunc func(c *edgecloud.ServiceClient) string, extraParams ...string) (client *edgecloud.ServiceClient, relativeURL string) {
	version := "v1"
	if extraParams != nil {
		version = extraParams[0]
	}

	client = fake.ServiceTokenClient(resourceName, version)

	resourceURL := ""
	if urlFunc == nil {
		resourceURL = client.ResourceBaseURL()
	} else {
		resourceURL = urlFunc(client)
	}

	parsedURL, err := url.Parse(resourceURL)
	if err != nil {
		panic(err)
	}

	relativeURL = parsedURL.Path

	return
}

func BuildTestMetadataListAll(resourceName string, resourceID string, extraParams ...string) TestFunc {
	return func(t *testing.T) {
		t.Helper()
		th.SetupHTTP()
		defer th.TeardownHTTP()

		client, relativeURL := prepareTestParams(resourceName, nil, extraParams...)

		th.Mux.HandleFunc(relativeURL, func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := fmt.Fprint(w, MetadataListResponse)
			if err != nil {
				logrus.Error(err)
			}
		})

		actual, err := metadata.ResourceMetadataListAll(client, resourceID)
		require.NoError(t, err)
		ct := actual[0]
		require.Equal(t, Metadata1, ct)
		require.Equal(t, ExpectedMetadataList, actual)
	}
}

func BuildTestMetadataGet(resourceName string, resourceID string, extraParams ...string) TestFunc {
	return func(t *testing.T) {
		t.Helper()
		th.SetupHTTP()
		defer th.TeardownHTTP()

		client, relativeURL := prepareTestParams(resourceName, func(c *edgecloud.ServiceClient) string {
			return metadata.ResourceMetadataItemURL(c, resourceID, ResourceMetadataReadOnly.Key)
		}, extraParams...)

		th.Mux.HandleFunc(relativeURL, func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := fmt.Fprint(w, MetadataResponse)
			if err != nil {
				logrus.Error(err)
			}
		})

		actual, err := metadata.ResourceMetadataGet(client, resourceID, ResourceMetadataReadOnly.Key).Extract()
		require.NoError(t, err)
		require.Equal(t, &ResourceMetadataReadOnly, actual)
	}
}

func BuildTestMetadataCreate(resourceName string, resourceID string, extraParams ...string) TestFunc {
	return func(t *testing.T) {
		t.Helper()
		th.SetupHTTP()
		defer th.TeardownHTTP()

		client, relativeURL := prepareTestParams(resourceName, func(c *edgecloud.ServiceClient) string {
			return metadata.ResourceMetadataURL(c, resourceID)
		}, extraParams...)

		th.Mux.HandleFunc(relativeURL, func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "POST")
			th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

			th.TestHeader(t, r, "Content-Type", "application/json")
			th.TestHeader(t, r, "Accept", "application/json")
			th.TestJSONRequest(t, r, MetadataCreateRequest)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
		})

		err := metadata.ResourceMetadataCreateOrUpdate(client, resourceID, map[string]string{
			"test1": "test1",
			"test2": "test2",
		}).ExtractErr()
		require.NoError(t, err)
	}
}

func BuildTestMetadataUpdate(resourceName string, resourceID string, extraParams ...string) TestFunc {
	return func(t *testing.T) {
		t.Helper()
		th.SetupHTTP()
		defer th.TeardownHTTP()

		client, relativeURL := prepareTestParams(resourceName, func(c *edgecloud.ServiceClient) string {
			return metadata.ResourceMetadataURL(c, resourceID)
		}, extraParams...)

		th.Mux.HandleFunc(relativeURL, func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "PUT")
			th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))

			th.TestHeader(t, r, "Content-Type", "application/json")
			th.TestHeader(t, r, "Accept", "application/json")
			th.TestJSONRequest(t, r, MetadataCreateRequest)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
		})

		err := metadata.ResourceMetadataReplace(client, resourceID, map[string]string{
			"test1": "test1",
			"test2": "test2",
		}).ExtractErr()
		require.NoError(t, err)
	}
}

func BuildTestMetadataDelete(resourceName string, resourceID string, extraParams ...string) TestFunc {
	return func(t *testing.T) {
		t.Helper()
		th.SetupHTTP()
		defer th.TeardownHTTP()

		client, relativeURL := prepareTestParams(resourceName, func(c *edgecloud.ServiceClient) string {
			return metadata.ResourceMetadataItemURL(c, resourceID, ResourceMetadataReadOnly.Key)
		}, extraParams...)

		th.Mux.HandleFunc(relativeURL, func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "DELETE")
			th.TestHeader(t, r, "Authorization", fmt.Sprintf("Bearer %s", fake.AccessToken))
			th.TestHeader(t, r, "Accept", "application/json")
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
		})

		err := metadata.ResourceMetadataDelete(client, resourceID, Metadata1.Key).ExtractErr()
		require.NoError(t, err)
	}
}
