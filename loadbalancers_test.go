package edgecloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadbalancers_List(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []Loadbalancer{{ID: testResourceID}}
	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.List(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_List_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Loadbalancers.List(ctx, nil)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestLoadbalancers_Get(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &Loadbalancer{ID: testResourceID}
	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.Get(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_Get_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Loadbalancers.Get(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestLoadbalancers_Create(t *testing.T) {
	setup()
	defer teardown()

	request := &LoadbalancerCreateRequest{
		Name:   "test-loadbalancer",
		Flavor: "lb1-1-2",
	}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(LoadbalancerCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, reqBody)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.Create(ctx, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_Create_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &LoadbalancerCreateRequest{}
	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Loadbalancers.Create(ctx, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestLoadbalancers_Create_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Loadbalancers.Create(ctx, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestLoadbalancers_Delete(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.Delete(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_Delete_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Loadbalancers.Delete(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestLoadbalancers_ListenerList(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []Listener{{ID: testResourceID}}
	URL := path.Join(lblistenersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.ListenerList(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_ListenerList_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(lblistenersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Loadbalancers.ListenerList(ctx, nil)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestLoadbalancers_ListenerGet(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &Listener{ID: testResourceID}
	URL := path.Join(lblistenersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.ListenerGet(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_ListenerGet_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(lblistenersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Loadbalancers.ListenerGet(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestLoadbalancers_ListenerCreate(t *testing.T) {
	setup()
	defer teardown()

	request := &ListenerCreateRequest{
		Name:           "test-loadbalancer",
		Protocol:       ListenerProtocolTCP,
		ProtocolPort:   80,
		LoadbalancerID: testResourceID,
	}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(lblistenersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(ListenerCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, reqBody)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.ListenerCreate(ctx, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_ListenerCreate_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &ListenerCreateRequest{}
	URL := path.Join(lblistenersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Loadbalancers.ListenerCreate(ctx, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestLoadbalancers_ListenerCreate_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Loadbalancers.ListenerCreate(ctx, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestLoadbalancers_ListenerDelete(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(lblistenersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.ListenerDelete(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_ListenerDelete_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(lblistenersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Loadbalancers.ListenerDelete(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestLoadbalancers_ListenerUpdate(t *testing.T) {
	setup()
	defer teardown()

	request := &ListenerUpdateRequest{Name: "test-listener"}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(lblistenersBasePathV2, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		reqBody := new(ListenerUpdateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, reqBody)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.ListenerUpdate(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_ListenerUpdate_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &ListenerUpdateRequest{}
	URL := path.Join(lblistenersBasePathV2, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Loadbalancers.ListenerUpdate(ctx, testResourceID, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestLoadbalancers_ListenerUpdate_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Loadbalancers.ListenerUpdate(ctx, testResourceID, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestLoadbalancers_ListenerRename(t *testing.T) {
	setup()
	defer teardown()

	request := &Name{Name: "test-listener"}
	expectedResp := &Listener{ID: testResourceID}
	URL := path.Join(lblistenersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		reqBody := new(Name)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, reqBody)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.ListenerRename(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_ListenerRename_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &Name{}
	URL := path.Join(lblistenersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Loadbalancers.ListenerRename(ctx, testResourceID, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestLoadbalancers_ListenerRename_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Loadbalancers.ListenerRename(ctx, testResourceID, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestLoadbalancers_PoolGet(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &Pool{ID: testResourceID}
	URL := path.Join(lbpoolsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.PoolGet(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_PoolGet_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(lbpoolsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Loadbalancers.PoolGet(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestLoadbalancers_PoolCreate(t *testing.T) {
	setup()
	defer teardown()

	request := &PoolCreateRequest{
		LoadbalancerPoolCreateRequest: LoadbalancerPoolCreateRequest{
			Name:       "test-loadbalancer",
			ListenerID: testResourceID,
		},
	}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(lbpoolsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(PoolCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, reqBody)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.PoolCreate(ctx, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_PoolCreate_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &PoolCreateRequest{}
	URL := path.Join(lbpoolsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Loadbalancers.PoolCreate(ctx, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestLoadbalancers_PoolCreate_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Loadbalancers.PoolCreate(ctx, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestLoadbalancers_PoolDelete(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(lbpoolsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.PoolDelete(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_PoolDelete_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(lbpoolsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)
	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Loadbalancers.PoolDelete(ctx, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestLoadbalancers_PoolUpdate(t *testing.T) {
	setup()
	defer teardown()

	request := &PoolUpdateRequest{
		ID:   testResourceID,
		Name: "test-lbpool",
	}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(lbpoolsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		reqBody := new(PoolUpdateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, reqBody)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.PoolUpdate(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_PoolUpdate_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &PoolUpdateRequest{}
	URL := path.Join(lbpoolsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Loadbalancers.PoolUpdate(ctx, testResourceID, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestLoadbalancers_PoolList(t *testing.T) {
	setup()
	defer teardown()

	options := PoolListOptions{LoadbalancerID: testResourceID}
	expectedResp := []Pool{{ID: testResourceID}}
	URL := path.Join(lbpoolsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.PoolList(ctx, &options)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_PoolList_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	options := PoolListOptions{}
	URL := path.Join(lbpoolsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Loadbalancers.PoolList(ctx, &options)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestLoadbalancers_PoolMemberCreate(t *testing.T) {
	setup()
	defer teardown()

	request := &PoolMemberCreateRequest{ID: testResourceID}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(lbpoolsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, loadbalancersMember)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(PoolMemberCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, reqBody)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.PoolMemberCreate(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_PoolMemberCreate_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &PoolMemberCreateRequest{}
	URL := path.Join(lbpoolsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, loadbalancersMember)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Loadbalancers.PoolMemberCreate(ctx, testResourceID, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestLoadbalancers_PoolMemberCreate_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Loadbalancers.PoolMemberCreate(ctx, testResourceID, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestLoadbalancers_PoolMemberDelete(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(lbpoolsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, loadbalancersMember, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.PoolMemberDelete(ctx, testResourceID, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_PoolMemberDelete_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(lbpoolsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, loadbalancersMember, testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Loadbalancers.PoolMemberDelete(ctx, testResourceID, testResourceID)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestLoadbalancers_HealthMonitorCreate(t *testing.T) {
	setup()
	defer teardown()

	request := &HealthMonitorCreateRequest{ID: testResourceID}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(lbpoolsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, loadbalancersHealthMonitor)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := new(HealthMonitorCreateRequest)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, reqBody)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.HealthMonitorCreate(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_HealthMonitorCreate_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &HealthMonitorCreateRequest{}
	URL := path.Join(lbpoolsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, loadbalancersHealthMonitor)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Loadbalancers.HealthMonitorCreate(ctx, testResourceID, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestLoadbalancers_HealthMonitorCreate_reqBodyNil(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Loadbalancers.HealthMonitorCreate(ctx, testResourceID, nil)
	assert.Nil(t, respActual)
	assert.Nil(t, resp)
	assert.EqualError(t, err, NewArgError("reqBody", "cannot be nil").Error())
}

func TestLoadbalancers_HealthMonitorDelete(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(lbpoolsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, loadbalancersHealthMonitor)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	resp, err := client.Loadbalancers.HealthMonitorDelete(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestLoadbalancers_HealthMonitorDelete_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(lbpoolsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, loadbalancersHealthMonitor)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	resp, err := client.Loadbalancers.HealthMonitorDelete(ctx, testResourceID)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestLoadbalancers_CheckLimits(t *testing.T) {
	setup()
	defer teardown()

	request := &LoadbalancerCheckLimitsRequest{}
	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), loadbalancersCheckLimits)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		resp, err := json.Marshal(map[string]int{})
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	_, resp, err := client.Loadbalancers.CheckLimits(ctx, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestLoadbalancers_CheckLimits_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &LoadbalancerCheckLimitsRequest{}
	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), loadbalancersCheckLimits)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Loadbalancers.CheckLimits(ctx, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestLoadbalancers_FlavorList(t *testing.T) {
	setup()
	defer teardown()

	options := FlavorsOptions{IncludePrices: true}
	expectedResp := []Flavor{{FlavorID: "g1-gpu-1-2-1"}}
	URL := path.Join(lbflavorsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.FlavorList(ctx, &options)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_FlavorList_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	options := FlavorsOptions{IncludePrices: true}
	URL := path.Join(lbflavorsBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Loadbalancers.FlavorList(ctx, &options)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestLoadbalancers_Rename(t *testing.T) {
	setup()
	defer teardown()

	request := &Name{Name: "test-loadbalancer"}
	expectedResp := &Loadbalancer{ID: testResourceID}
	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		reqBody := new(Name)
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, reqBody)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.Rename(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_Rename_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	request := &Name{}
	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Loadbalancers.Rename(ctx, testResourceID, request)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestLoadbalancers_MetadataList(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []MetadataDetailed{{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}}
	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.MetadataList(ctx, testResourceID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_MetadataCreate(t *testing.T) {
	setup()
	defer teardown()

	request := &Metadata{"key": "value"}
	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
	})

	resp, err := client.Loadbalancers.MetadataCreate(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestLoadbalancers_MetadataUpdate(t *testing.T) {
	setup()
	defer teardown()

	request := &Metadata{"key": "value"}
	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
	})

	resp, err := client.Loadbalancers.MetadataUpdate(ctx, testResourceID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestLoadbalancers_MetadataDeleteItem(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataItemPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	resp, err := client.Loadbalancers.MetadataDeleteItem(ctx, testResourceID, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func TestLoadbalancers_MetadataGetItem(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &MetadataDetailed{
		Key:      "image_id",
		Value:    "b3c52ece-147e-4af5-8d7c-84691309b879",
		ReadOnly: true,
	}
	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, metadataItemPath)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.MetadataGetItem(ctx, testResourceID, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_MetricsList(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []LoadbalancerMetrics{{}}
	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, loadbalancersMetrics)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.Loadbalancers.MetricsList(ctx, testResourceID, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestLoadbalancers_MetricsList_ResponseError(t *testing.T) {
	setup()
	defer teardown()

	URL := path.Join(loadbalancersBasePathV1, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID, loadbalancersMetrics)
	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Bad request")
	})

	respActual, resp, err := client.Loadbalancers.MetricsList(ctx, testResourceID, nil)
	assert.Nil(t, respActual)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 400)
}

func TestLoadbalancers_isValidUUID_Error_Return_Loadbalancer(t *testing.T) {
	setup()
	defer teardown()

	tests := []struct {
		name     string
		testFunc func() (*Loadbalancer, *Response, error)
	}{
		{
			name: "Get",
			testFunc: func() (*Loadbalancer, *Response, error) {
				return client.Loadbalancers.Get(ctx, testResourceIDNotValidUUID)
			},
		},
		{
			name: "Rename",
			testFunc: func() (*Loadbalancer, *Response, error) {
				return client.Loadbalancers.Rename(ctx, testResourceIDNotValidUUID, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respActual, resp, err := tt.testFunc()
			require.Nil(t, respActual)
			require.Equal(t, 400, resp.StatusCode)
			require.EqualError(t, err, NewArgError("loadbalancerID", NotCorrectUUID).Error())
		})
	}
}

func TestLoadbalancers_isValidUUID_Error_Return_Listener(t *testing.T) {
	setup()
	defer teardown()

	tests := []struct {
		name     string
		testFunc func() (*Listener, *Response, error)
	}{
		{
			name: "ListenerGet",
			testFunc: func() (*Listener, *Response, error) {
				return client.Loadbalancers.ListenerGet(ctx, testResourceIDNotValidUUID)
			},
		},
		{
			name: "Rename",
			testFunc: func() (*Listener, *Response, error) {
				return client.Loadbalancers.ListenerRename(ctx, testResourceIDNotValidUUID, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respActual, resp, err := tt.testFunc()
			require.Nil(t, respActual)
			require.Equal(t, 400, resp.StatusCode)
			require.EqualError(t, err, NewArgError("listenerID", NotCorrectUUID).Error())
		})
	}
}

func TestLoadbalancers_isValidUUID_Error_Return_TaskResponse(t *testing.T) {
	setup()
	defer teardown()

	tests := []struct {
		name     string
		testFunc func() (*TaskResponse, *Response, error)
		argError *ArgError
	}{
		{
			name: "LoadBalancerDelete",
			testFunc: func() (*TaskResponse, *Response, error) {
				return client.Loadbalancers.Delete(ctx, testResourceIDNotValidUUID)
			},
			argError: NewArgError("loadbalancerID", NotCorrectUUID),
		},
		{
			name: "ListenerDelete",
			testFunc: func() (*TaskResponse, *Response, error) {
				return client.Loadbalancers.ListenerDelete(ctx, testResourceIDNotValidUUID)
			},
			argError: NewArgError("listenerID", NotCorrectUUID),
		},
		{
			name: "ListenerUpdate",
			testFunc: func() (*TaskResponse, *Response, error) {
				return client.Loadbalancers.ListenerUpdate(ctx, testResourceIDNotValidUUID, nil)
			},
			argError: NewArgError("listenerID", NotCorrectUUID),
		},
		{
			name: "PoolDelete",
			testFunc: func() (*TaskResponse, *Response, error) {
				return client.Loadbalancers.PoolDelete(ctx, testResourceIDNotValidUUID)
			},
			argError: NewArgError("poolID", NotCorrectUUID),
		},
		{
			name: "PoolUpdate",
			testFunc: func() (*TaskResponse, *Response, error) {
				return client.Loadbalancers.PoolUpdate(ctx, testResourceIDNotValidUUID, nil)
			},
			argError: NewArgError("poolID", NotCorrectUUID),
		},
		{
			name: "PoolMemberCreate",
			testFunc: func() (*TaskResponse, *Response, error) {
				return client.Loadbalancers.PoolMemberCreate(ctx, testResourceIDNotValidUUID, nil)
			},
			argError: NewArgError("poolID", NotCorrectUUID),
		},
		{
			name: "HealthMonitorCreate",
			testFunc: func() (*TaskResponse, *Response, error) {
				return client.Loadbalancers.HealthMonitorCreate(ctx, testResourceIDNotValidUUID, nil)
			},
			argError: NewArgError("poolID", NotCorrectUUID),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respActual, resp, err := tt.testFunc()
			require.Nil(t, respActual)
			require.Equal(t, 400, resp.StatusCode)
			require.EqualError(t, err, tt.argError.Error())
		})
	}
}

func TestLoadbalancers_PoolMemberDelete_isValidUUID_Error(t *testing.T) {
	setup()
	defer teardown()

	tests := []struct {
		name     string
		testFunc func() (*TaskResponse, *Response, error)
		argError *ArgError
	}{
		{
			name: "PoolMemberDeletePoolIdError",
			testFunc: func() (*TaskResponse, *Response, error) {
				return client.Loadbalancers.PoolMemberDelete(ctx, testResourceIDNotValidUUID, testResourceID)
			},
			argError: NewArgError("poolID", NotCorrectUUID),
		},
		{
			name: "PoolMemberDeleteMemberIdError",
			testFunc: func() (*TaskResponse, *Response, error) {
				return client.Loadbalancers.PoolMemberDelete(ctx, testResourceID, testResourceIDNotValidUUID)
			},
			argError: NewArgError("memberID", NotCorrectUUID),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respActual, resp, err := tt.testFunc()
			require.Nil(t, respActual)
			require.Equal(t, 400, resp.StatusCode)
			require.EqualError(t, err, tt.argError.Error())
		})
	}
}

func TestLoadbalancers_PoolGet_isValidUUID_Error(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Loadbalancers.PoolGet(ctx, testResourceIDNotValidUUID)
	require.Nil(t, respActual)
	require.Equal(t, 400, resp.StatusCode)
	require.EqualError(t, err, NewArgError("poolID", NotCorrectUUID).Error())
}

func TestLoadbalancers_HealthMonitorDelete_isValidUUID_Error(t *testing.T) {
	setup()
	defer teardown()

	resp, err := client.Loadbalancers.HealthMonitorDelete(ctx, testResourceIDNotValidUUID)
	require.Equal(t, 400, resp.StatusCode)
	require.EqualError(t, err, NewArgError("poolID", NotCorrectUUID).Error())
}

func TestLoadbalancers_MetricsList_isValidUUID_Error(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Loadbalancers.MetricsList(ctx, testResourceIDNotValidUUID, nil)
	require.Nil(t, respActual)
	require.Equal(t, 400, resp.StatusCode)
	require.EqualError(t, err, NewArgError("loadbalancerID", NotCorrectUUID).Error())
}

func TestLoadbalancers_Metadata_isValidUUID_Error_Return_Response(t *testing.T) {
	setup()
	defer teardown()

	tests := []struct {
		name     string
		testFunc func() (*Response, error)
	}{
		{
			name: "MetadataCreate",
			testFunc: func() (*Response, error) {
				return client.Loadbalancers.MetadataCreate(ctx, testResourceIDNotValidUUID, nil)
			},
		},
		{
			name: "MetadataUpdate",
			testFunc: func() (*Response, error) {
				return client.Loadbalancers.MetadataUpdate(ctx, testResourceIDNotValidUUID, nil)
			},
		},
		{
			name: "MetadataDeleteItem",
			testFunc: func() (*Response, error) {
				return client.Loadbalancers.MetadataDeleteItem(ctx, testResourceIDNotValidUUID, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.testFunc()
			require.Equal(t, 400, resp.StatusCode)
			require.EqualError(t, err, NewArgError("loadbalancerID", NotCorrectUUID).Error())
		})
	}
}

func TestLoadbalancers_MetadataList_isValidUUID_Error(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Loadbalancers.MetadataList(ctx, testResourceIDNotValidUUID)
	require.Nil(t, respActual)
	require.Equal(t, 400, resp.StatusCode)
	require.EqualError(t, err, NewArgError("loadbalancerID", NotCorrectUUID).Error())
}

func TestLoadbalancers_MetadataGetItem_isValidUUID_Error(t *testing.T) {
	setup()
	defer teardown()

	respActual, resp, err := client.Loadbalancers.MetadataGetItem(ctx, testResourceIDNotValidUUID, nil)
	require.Nil(t, respActual)
	require.Equal(t, 400, resp.StatusCode)
	require.EqualError(t, err, NewArgError("loadbalancerID", NotCorrectUUID).Error())
}
