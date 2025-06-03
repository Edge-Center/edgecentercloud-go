package edgecloud

import (
	"encoding/json"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResellerImagesV2_ListByRole(t *testing.T) {
	setup()
	defer teardown()

	ids := ImageIDs{
		"b5b4d65d-945f-4b98-ab6f-332319c724ef",
		"0052a312-e6d8-4177-8e29-b017a3a6b588",
	}

	ri := ResellerImageV2{
		ImageIDs:   &ids,
		RegionID:   8,
		EntityType: ResellerType,
		EntityID:   123456,
	}

	expectedResp := ResellerImageV2List{
		Count:   1,
		Results: []ResellerImageV2{ri},
	}

	mux.HandleFunc(resellerImageBasePathV2, func(w http.ResponseWriter, r *http.Request) {
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

	respActual, resp, err := client.ResellerImageV2.ListByRole(ctx)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, *respActual, expectedResp)
}

func TestResellerImagesV2_List(t *testing.T) {
	setup()
	defer teardown()

	ids := ImageIDs{
		"b5b4d65d-945f-4b98-ab6f-332319c724ef",
		"0052a312-e6d8-4177-8e29-b017a3a6b588",
	}

	ri := ResellerImageV2{
		ImageIDs:   &ids,
		RegionID:   8,
		EntityType: ResellerType,
		EntityID:   123456,
	}

	expectedResp := ResellerImageV2List{
		Count:   1,
		Results: []ResellerImageV2{ri},
	}

	URL := path.Join(resellerImageBasePathV2, "/", "reseller", "/", "123456")

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
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

	respActual, resp, err := client.ResellerImageV2.List(ctx, "reseller", 123456)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, *respActual, expectedResp)
}

func TestResellerImagesV2_Update(t *testing.T) {
	setup()
	defer teardown()

	ids := ImageIDs{
		"b5b4d65d-945f-4b98-ab6f-332319c724ef",
		"0052a312-e6d8-4177-8e29-b017a3a6b588",
	}

	updateReq := ResellerImageV2UpdateRequest{
		ImageIDs:   &ids,
		RegionID:   8,
		EntityType: ResellerType,
		EntityID:   123456,
	}

	expectedResp := ResellerImageV2{
		ImageIDs:   &ids,
		RegionID:   8,
		EntityType: ResellerType,
		EntityID:   123456,
	}

	URL := path.Join(resellerImageBasePathV2)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, err = w.Write(resp)
		if err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	})

	respActual, resp, err := client.ResellerImageV2.Update(ctx, &updateReq)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, *respActual, expectedResp)
}

func TestResellerImagesV2_Delete(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(resellerImageBasePathV2, "/", "reseller", "/", "123456")

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	resp, err := client.ResellerImageV2.Delete(ctx, ResellerType, 123456, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}
