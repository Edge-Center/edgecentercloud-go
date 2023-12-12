package util

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

func TestLoadbalancerGetByName(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	lbs := []edgecloud.Loadbalancer{
		{
			Name: testResourceID,
		},
	}
	URL := path.Join("/v1/loadbalancers", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(lbs)
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	lb, err := LoadbalancerGetByName(context.Background(), client, testResourceID)
	assert.NoError(t, err)
	assert.Equal(t, testResourceID, lb.Name)
}

func TestLoadbalancerGetByName_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	URL := path.Join("/v1/loadbalancers", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	lb, err := LoadbalancerGetByName(context.Background(), client, testResourceID)
	assert.Error(t, err)
	assert.Nil(t, lb)
}

func TestLoadbalancerGetByName_CustomErrors(t *testing.T) {
	tests := []struct {
		name      string
		lbs       []edgecloud.Loadbalancer
		expectErr error
	}{
		{
			name:      "ErrLoadbalancersNotFound",
			lbs:       nil,
			expectErr: ErrLoadbalancersNotFound,
		},
		{
			name: "ErrMultipleResults",
			lbs: []edgecloud.Loadbalancer{
				{
					Name: testName,
				},
				{
					Name: testName,
				},
			},
			expectErr: ErrMultipleResults,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mux := http.NewServeMux()
			server := httptest.NewServer(mux)
			defer server.Close()

			URL := path.Join("/v1/loadbalancers", strconv.Itoa(projectID), strconv.Itoa(regionID))

			mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
				resp, err := json.Marshal(tc.lbs)
				if err != nil {
					t.Fatalf("failed to marshal JSON: %v", err)
				}
				_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
			})

			client := edgecloud.NewClient(nil)
			baseURL, _ := url.Parse(server.URL)
			client.BaseURL = baseURL
			client.Project = projectID
			client.Region = regionID

			lbs, err := LoadbalancerGetByName(context.Background(), client, testName)
			assert.ErrorIs(t, err, tc.expectErr)
			assert.Nil(t, lbs)
		})
	}
}

func TestLBListenerGetByName(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	listeners := []edgecloud.Listener{
		{
			ID:   testResourceID,
			Name: testName,
		},
	}
	URL := path.Join("/v1/lblisteners", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(listeners)
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	lis, err := LBListenerGetByName(context.Background(), client, testName, testResourceID)
	assert.NoError(t, err)
	assert.Equal(t, testName, lis.Name)
}

func TestLBListenerGetByName_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	URL := path.Join("/v1/lblisteners", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	lis, err := LBListenerGetByName(context.Background(), client, testName, testResourceID)
	assert.Error(t, err)
	assert.Nil(t, lis)
}

func TestLBListenerGetByName_CustomErrors(t *testing.T) {
	tests := []struct {
		name      string
		listeners []edgecloud.Listener
		expectErr error
	}{
		{
			name:      "ErrLoadbalancerListenerNotFound",
			listeners: nil,
			expectErr: ErrLoadbalancerListenerNotFound,
		},
		{
			name: "ErrMultipleResults",
			listeners: []edgecloud.Listener{
				{
					ID:   testResourceID,
					Name: testName,
				},
				{
					ID:   testResourceID,
					Name: testName,
				},
			},
			expectErr: ErrMultipleResults,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mux := http.NewServeMux()
			server := httptest.NewServer(mux)
			defer server.Close()

			URL := path.Join("/v1/lblisteners", strconv.Itoa(projectID), strconv.Itoa(regionID))

			mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
				resp, err := json.Marshal(tc.listeners)
				if err != nil {
					t.Fatalf("failed to marshal JSON: %v", err)
				}
				_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
			})

			client := edgecloud.NewClient(nil)
			baseURL, _ := url.Parse(server.URL)
			client.BaseURL = baseURL
			client.Project = projectID
			client.Region = regionID

			listeners, err := LBListenerGetByName(context.Background(), client, testName, testResourceID)
			assert.ErrorIs(t, err, tc.expectErr)
			assert.Nil(t, listeners)
		})
	}
}

func TestLBPoolGetByName(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	pools := []edgecloud.Pool{
		{
			ID:   testResourceID,
			Name: testName,
		},
	}
	URL := path.Join("/v1/lbpools", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(pools)
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	pool, err := LBPoolGetByName(context.Background(), client, testName, testResourceID)
	assert.NoError(t, err)
	assert.Equal(t, testName, pool.Name)
}

func TestLBPoolGetByName_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	URL := path.Join("/v1/lbpools", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	pool, err := LBPoolGetByName(context.Background(), client, testName, testResourceID)
	assert.Error(t, err)
	assert.Nil(t, pool)
}

func TestLBPoolGetByName_CustomErrors(t *testing.T) {
	tests := []struct {
		name      string
		pools     []edgecloud.Pool
		expectErr error
	}{
		{
			name:      "ErrLoadbalancerPoolsNotFound",
			pools:     nil,
			expectErr: ErrLoadbalancerPoolsNotFound,
		},
		{
			name: "ErrMultipleResults",
			pools: []edgecloud.Pool{
				{
					ID:   testResourceID,
					Name: testName,
				},
				{
					ID:   testResourceID,
					Name: testName,
				},
			},
			expectErr: ErrMultipleResults,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mux := http.NewServeMux()
			server := httptest.NewServer(mux)
			defer server.Close()

			URL := path.Join("/v1/lbpools", strconv.Itoa(projectID), strconv.Itoa(regionID))

			mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
				resp, err := json.Marshal(tc.pools)
				if err != nil {
					t.Fatalf("failed to marshal JSON: %v", err)
				}
				_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
			})

			client := edgecloud.NewClient(nil)
			baseURL, _ := url.Parse(server.URL)
			client.BaseURL = baseURL
			client.Project = projectID
			client.Region = regionID

			pools, err := LBPoolGetByName(context.Background(), client, testName, testResourceID)
			assert.ErrorIs(t, err, tc.expectErr)
			assert.Nil(t, pools)
		})
	}
}

func TestWaitLoadbalancerProvisioningStatusActive(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	expectedResp := edgecloud.Loadbalancer{
		Name:               testResourceID,
		ProvisioningStatus: edgecloud.ProvisioningStatusActive,
	}
	URL := path.Join("/v1/loadbalancers", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	err := WaitLoadbalancerProvisioningStatusActive(context.Background(), client, testResourceID, nil)
	assert.NoError(t, err)
}

func TestWaitLoadbalancerProvisioningStatusActive_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	URL := path.Join("/v1/loadbalancers", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	var attempts uint = 2
	err := WaitLoadbalancerProvisioningStatusActive(context.Background(), client, testResourceID, &attempts)
	assert.Error(t, err)
}

func TestWaitLoadbalancerProvisioningStatusActive_CustomErrors(t *testing.T) {
	tests := []struct {
		name               string
		provisioningStatus edgecloud.ProvisioningStatus
		expectErr          error
	}{
		{
			name:               "ErrErrorState",
			provisioningStatus: edgecloud.ProvisioningStatusError,
			expectErr:          ErrErrorState,
		},
		{
			name:               "ErrNotActiveStatus",
			provisioningStatus: edgecloud.ProvisioningStatusPendingCreate,
			expectErr:          ErrNotActiveStatus,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mux := http.NewServeMux()
			server := httptest.NewServer(mux)
			defer server.Close()

			expectedResp := edgecloud.Loadbalancer{
				Name:               testResourceID,
				ProvisioningStatus: tc.provisioningStatus,
			}
			URL := path.Join("/v1/loadbalancers", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

			mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
				resp, err := json.Marshal(expectedResp)
				if err != nil {
					t.Errorf("failed to marshal response: %v", err)
				}
				_, _ = fmt.Fprint(w, string(resp))
			})

			client := edgecloud.NewClient(nil)
			baseURL, _ := url.Parse(server.URL)
			client.BaseURL = baseURL
			client.Project = projectID
			client.Region = regionID

			var attempts uint = 2
			err := WaitLoadbalancerProvisioningStatusActive(context.Background(), client, testResourceID, &attempts)
			assert.Error(t, err)
		})
	}
}

func TestFindPoolMemberByAddressPortAndSubnetID(t *testing.T) {
	testAddress := net.IP("192.168.1.1")
	testProtocolPort := 80
	testPool := edgecloud.Pool{
		Members: []edgecloud.PoolMember{
			{
				ID: testResourceID,
				PoolMemberCreateRequest: edgecloud.PoolMemberCreateRequest{
					Address:      testAddress,
					ProtocolPort: testProtocolPort,
					SubnetID:     testName,
				},
			},
		},
	}

	found := FindPoolMemberByAddressPortAndSubnetID(testPool, testAddress, testProtocolPort, testName)
	assert.True(t, found)
}

func TestFindPoolMemberByAddressPortAndSubnetID_MemberDoesNotExist(t *testing.T) {
	testProtocolPort := 80
	testPool := edgecloud.Pool{
		Members: []edgecloud.PoolMember{
			{
				ID: testResourceID,
				PoolMemberCreateRequest: edgecloud.PoolMemberCreateRequest{
					Address:      net.IP("192.168.1.1"),
					ProtocolPort: testProtocolPort,
					SubnetID:     testName,
				},
			},
		},
	}

	found := FindPoolMemberByAddressPortAndSubnetID(testPool, net.IP("192.168.1.2"), testProtocolPort, testName)
	assert.False(t, found)
}

func TestFindPoolMemberByAddressPortAndSubnetID_MemberAddressIsNil(t *testing.T) {
	testProtocolPort := 80
	testPool := edgecloud.Pool{
		Members: []edgecloud.PoolMember{
			{
				ID: testResourceID,
				PoolMemberCreateRequest: edgecloud.PoolMemberCreateRequest{
					ProtocolPort: testProtocolPort,
					SubnetID:     testName,
				},
			},
		},
	}
	found := FindPoolMemberByAddressPortAndSubnetID(testPool, net.IP("192.168.1.2"), testProtocolPort, testName)
	assert.False(t, found)
}

func TestLBSharedPoolList(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	pools := []edgecloud.Pool{
		{
			Listeners: []edgecloud.ID{{ID: testResourceID}},
		},
		{
			Listeners: []edgecloud.ID{{ID: testResourceID}},
		},
	}
	URL := path.Join("/v1/lbpools", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(pools)
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	sharedPools, err := LBSharedPoolList(context.Background(), client, testResourceID)
	assert.NoError(t, err)
	assert.Len(t, sharedPools, 2)
}

func TestLBSharedPoolList_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	URL := path.Join("/v1/lbpools", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	sharedPools, err := LBSharedPoolList(context.Background(), client, testResourceID)
	assert.Error(t, err)
	assert.Nil(t, sharedPools)
}

func TestDeletePoolByNameIfExist(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	pools := []edgecloud.Pool{{ID: testResourceID, Name: testName}}
	URLLBPoolsGet := path.Join("/v1/lbpools", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URLLBPoolsGet, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(pools)
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	taskDelete := &edgecloud.TaskResponse{Tasks: []string{testResourceID}}
	URLLBPoolsDelete := path.Join("/v1/lbpools", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URLLBPoolsDelete, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(taskDelete)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	task := &edgecloud.Task{ID: testResourceID, State: edgecloud.TaskStateFinished}
	URLTasks := path.Join("/v1/tasks", testResourceID)

	mux.HandleFunc(URLTasks, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(task)
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	err := DeletePoolByNameIfExist(context.Background(), client, testName, testResourceID)
	assert.NoError(t, err)
}

func TestDeletePoolByNameIfExist_LBPoolGetByName_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	URLLBPoolsGet := path.Join("/v1/lbpools", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URLLBPoolsGet, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	err := DeletePoolByNameIfExist(context.Background(), client, testName, testResourceID)
	assert.Error(t, err)
}

func TestDeletePoolByNameIfExist_LBPoolGetByName_ErrLoadbalancerPoolsNotFound(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	URL := path.Join("/v1/lbpools", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(nil)
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	err := DeletePoolByNameIfExist(context.Background(), client, testName, testResourceID)
	assert.NoError(t, err)
}

func TestDeletePoolByNameIfExist_PoolDelete_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	pools := []edgecloud.Pool{{ID: testResourceID, Name: testName}}
	URLLBPoolsGet := path.Join("/v1/lbpools", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URLLBPoolsGet, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(pools)
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	URLLBPoolsDelete := path.Join("/v1/lbpools", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URLLBPoolsDelete, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	err := DeletePoolByNameIfExist(context.Background(), client, testName, testResourceID)
	assert.Error(t, err)
}

func TestDeleteUnusedPools(t *testing.T) {
	oldPools := []edgecloud.Pool{{ID: testResourceID, Loadbalancers: []edgecloud.ID{{ID: testResourceID}}}}

	err := DeleteUnusedPools(context.Background(), nil, oldPools, []string{testResourceID}, nil)
	assert.NoError(t, err)
}

func TestDeleteUnusedPools_NotExist(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	oldPools := []edgecloud.Pool{{ID: testResourceID, Loadbalancers: []edgecloud.ID{{ID: testResourceID}}}}

	taskDelete := &edgecloud.TaskResponse{Tasks: []string{testResourceID}}
	URLLBPoolsDelete := path.Join("/v1/lbpools", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URLLBPoolsDelete, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(taskDelete)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	task := &edgecloud.Task{ID: testResourceID, State: edgecloud.TaskStateFinished}
	URLTasks := path.Join("/v1/tasks", testResourceID)

	mux.HandleFunc(URLTasks, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(task)
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	expectedResp := edgecloud.Loadbalancer{
		Name:               testResourceID,
		ProvisioningStatus: edgecloud.ProvisioningStatusActive,
	}
	URL := path.Join("/v1/loadbalancers", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	err := DeleteUnusedPools(context.Background(), client, oldPools, []string{}, nil)
	assert.NoError(t, err)
}

func TestDeleteUnusedPools_NotExist_PoolDelete_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	oldPools := []edgecloud.Pool{{ID: testResourceID, Loadbalancers: []edgecloud.ID{{ID: testResourceID}}}}

	URLLBPoolsDelete := path.Join("/v1/lbpools", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URLLBPoolsDelete, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	err := DeleteUnusedPools(context.Background(), client, oldPools, []string{}, nil)
	assert.Error(t, err)
}

func TestDeleteUnusedPools_NotExist_WaitForTaskComplete_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	oldPools := []edgecloud.Pool{{ID: testResourceID, Loadbalancers: []edgecloud.ID{{ID: testResourceID}}}}

	taskDelete := &edgecloud.TaskResponse{Tasks: []string{testResourceID}}
	URLLBPoolsDelete := path.Join("/v1/lbpools", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URLLBPoolsDelete, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(taskDelete)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	URLTasks := path.Join("/v1/tasks", testResourceID)

	mux.HandleFunc(URLTasks, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(&edgecloud.Task{ID: testResourceID, State: edgecloud.TaskStateError})
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	err := DeleteUnusedPools(context.Background(), client, oldPools, []string{}, nil)
	assert.Error(t, err)
}

func TestDeleteUnusedPools_WaitLoadbalancerProvisioningStatusActive_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	oldPools := []edgecloud.Pool{{ID: testResourceID, Loadbalancers: []edgecloud.ID{{ID: testResourceID}}}}

	taskDelete := &edgecloud.TaskResponse{Tasks: []string{testResourceID}}
	URLLBPoolsDelete := path.Join("/v1/lbpools", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URLLBPoolsDelete, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(taskDelete)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	task := &edgecloud.Task{ID: testResourceID, State: edgecloud.TaskStateFinished}
	URLTasks := path.Join("/v1/tasks", testResourceID)

	mux.HandleFunc(URLTasks, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(task)
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	expectedResp := edgecloud.Loadbalancer{
		Name:               testResourceID,
		ProvisioningStatus: edgecloud.ProvisioningStatusError,
	}
	URL := path.Join("/v1/loadbalancers", strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	var attempts uint = 2
	err := DeleteUnusedPools(context.Background(), client, oldPools, []string{}, &attempts)
	assert.Error(t, err)
}
