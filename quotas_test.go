package edgecloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQuotas_ListCombined(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &CombinedQuota{}

	mux.HandleFunc(quotasClientBasePathV2, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Quotas.ListCombined(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestQuotas_ListCombinedWithOptions(t *testing.T) {
	setup()
	defer teardown()

	options := &ListCombinedOptions{ClientID: 123}
	expectedResp := &CombinedQuota{}

	mux.HandleFunc(quotasClientBasePathV2, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Quotas.ListCombined(ctx, options)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestQuotas_ListGlobal(t *testing.T) {
	setup()
	defer teardown()

	options := &ListGlobalOptions{ClientID: 123}
	expectedResp := &Quota{}

	mux.HandleFunc(quotasGlobalBasePathV2, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Quotas.ListGlobal(ctx, options)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestQuotas_ListRegional(t *testing.T) {
	setup()
	defer teardown()

	options := &ListRegionalOptions{
		ClientID: 123,
		RegionID: regionID,
	}
	expectedResp := &Quota{}

	mux.HandleFunc(quotasRegionalBasePathV2, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Quotas.ListRegional(ctx, options)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestQuotas_DeleteNotificationThreshold(t *testing.T) {
	setup()
	defer teardown()

	clientID := 8
	URL := path.Join(quotasClientBasePathV2, strconv.Itoa(clientID), quotasNotificationThreshold)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	resp, err := client.Quotas.DeleteNotificationThreshold(ctx, clientID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestNetworks_GetNotificationThreshold(t *testing.T) {
	setup()
	defer teardown()

	clientID := 8
	expectedResp := &QuotaNotificationThreshold{Threshold: 1}
	URL := path.Join(quotasClientBasePathV2, strconv.Itoa(clientID), quotasNotificationThreshold)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Quotas.GetNotificationThreshold(ctx, clientID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestNetworks_UpdateNotificationThreshold(t *testing.T) {
	setup()
	defer teardown()

	clientID := 8
	request := &NotificationThresholdUpdateRequest{Threshold: 1}
	expectedResp := &QuotaNotificationThreshold{Threshold: 1}
	URL := path.Join(quotasClientBasePathV2, strconv.Itoa(clientID), quotasNotificationThreshold)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Quotas.UpdateNotificationThreshold(ctx, clientID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}
