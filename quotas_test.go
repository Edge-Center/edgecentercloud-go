package edgecloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuotas_ListCombined(t *testing.T) {
	setup()
	defer teardown()

	combinedQuotaResponse := &CombinedQuota{}

	mux.HandleFunc(quotasClientBasePathV2, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(combinedQuotaResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Quotas.ListCombined(ctx, nil)
	require.NoError(t, err)

	assert.Equal(t, combinedQuotaResponse, resp)
}

func TestQuotas_ListCombinedWithOptions(t *testing.T) {
	setup()
	defer teardown()

	combinedQuotaResponse := &CombinedQuota{}
	listCombinedOptions := &ListCombinedOptions{
		ClientID: 123,
	}

	mux.HandleFunc(quotasClientBasePathV2, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(combinedQuotaResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Quotas.ListCombined(ctx, listCombinedOptions)
	require.NoError(t, err)

	assert.Equal(t, combinedQuotaResponse, resp)
}

func TestQuotas_ListGlobal(t *testing.T) {
	setup()
	defer teardown()

	quotaResponse := &Quota{}
	listCombinedOptions := &ListGlobalOptions{
		ClientID: 123,
	}

	mux.HandleFunc(quotasGlobalBasePathV2, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(quotaResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Quotas.ListGlobal(ctx, listCombinedOptions)
	require.NoError(t, err)

	assert.Equal(t, quotaResponse, resp)
}

func TestQuotas_ListRegional(t *testing.T) {
	setup()
	defer teardown()

	quotaResponse := &Quota{}
	listRegionalOptions := &ListRegionalOptions{
		ClientID: 123,
		RegionID: regionID,
	}

	mux.HandleFunc(quotasRegionalBasePathV2, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(quotaResponse)
		_, _ = fmt.Fprint(w, string(resp))
	})

	resp, _, err := client.Quotas.ListRegional(ctx, listRegionalOptions)
	require.NoError(t, err)

	assert.Equal(t, quotaResponse, resp)
}
