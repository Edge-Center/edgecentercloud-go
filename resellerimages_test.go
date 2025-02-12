package edgecloud

import (
	"encoding/json"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResellerImages_List(t *testing.T) {
	setup()
	defer teardown()

	ids := ImageIDs{
		"b5b4d65d-945f-4b98-ab6f-332319c724ef",
		"0052a312-e6d8-4177-8e29-b017a3a6b588",
	}

	ri := ResellerImage{
		ImageIDs: &ids,
		RegionID: 8,
	}

	expectedResp := ResellerImageList{
		Count:   1,
		Results: []ResellerImage{ri},
	}

	URL := path.Join(resellerImageBasePathV1, "/", "123456")

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

	respActual, resp, err := client.ResellerImage.List(ctx, 123456)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, *respActual, expectedResp)
}

func TestResellerImages_Update(t *testing.T) {
	setup()
	defer teardown()

	ids := ImageIDs{
		"b5b4d65d-945f-4b98-ab6f-332319c724ef",
		"0052a312-e6d8-4177-8e29-b017a3a6b588",
	}

	updateReq := ResellerImageUpdateRequest{
		ImageIDs:   &ids,
		RegionID:   8,
		ResellerID: 936337,
	}

	expectedResp := ResellerImage{
		ImageIDs: &ids,
		RegionID: 8,
	}

	URL := path.Join(resellerImageBasePathV1)

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

	respActual, resp, err := client.ResellerImage.Update(ctx, &updateReq)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, *respActual, expectedResp)
}

func TestResellerImages_Delete(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(resellerImageBasePathV1, "/", "123456")

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	resp, err := client.ResellerImage.Delete(ctx, 123456)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}
