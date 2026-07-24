package util

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	edgecloud "github.com/Edge-Center/edgecentercloud-go/v2"
)

func TestDBaaSClusterGetByCreatorTaskID(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	clusters := []edgecloud.DBaaSCluster{
		{ID: testResourceID, Name: testName},
		{ID: testResourceID2, Name: testName},
	}
	listURL := path.Join(edgecloud.DBaaSClustersBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID))
	getURLFirst := path.Join(edgecloud.DBaaSClustersBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)
	getURLSecond := path.Join(edgecloud.DBaaSClustersBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID2)

	mux.HandleFunc(listURL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(clusters)
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	mux.HandleFunc(getURLFirst, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(&edgecloud.DBaaSCluster{
			ID:            testResourceID,
			Name:          testName,
			CreatorTaskID: testResourceID2,
		})
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	mux.HandleFunc(getURLSecond, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(&edgecloud.DBaaSCluster{
			ID:            testResourceID2,
			Name:          testName,
			CreatorTaskID: testResourceID,
		})
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

	cluster, err := DBaaSClusterGetByCreatorTaskID(context.Background(), client, testName, testResourceID)
	assert.NoError(t, err)
	assert.NotNil(t, cluster)
	assert.Equal(t, testResourceID2, cluster.ID)
}

func TestDBaaSClusterGetByCreatorTaskID_NotFound(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	clusters := []edgecloud.DBaaSCluster{{ID: testResourceID, Name: testName}}
	listURL := path.Join(edgecloud.DBaaSClustersBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID))
	getURL := path.Join(edgecloud.DBaaSClustersBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(listURL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(clusters)
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	mux.HandleFunc(getURL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(&edgecloud.DBaaSCluster{
			ID:            testResourceID,
			Name:          testName,
			CreatorTaskID: testResourceID2,
		})
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

	cluster, err := DBaaSClusterGetByCreatorTaskID(context.Background(), client, testName, testResourceID)
	assert.ErrorIs(t, err, ErrDBaaSClustersNotFound)
	assert.Nil(t, cluster)
}

func TestWaitDBaaSClusterByCreatorTaskID_Timeout(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	listURL := path.Join(edgecloud.DBaaSClustersBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID))
	mux.HandleFunc(listURL, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"results":[]}`)
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	cluster, err := WaitDBaaSClusterByCreatorTaskID(context.Background(), client, testName, testResourceID, time.Millisecond)
	assert.Equal(t, edgecloud.NewArgError("task error", errTaskWaitTimeout.Error()), err)
	assert.Nil(t, cluster)
}

func TestWaitDBaaSClusterStatusHealthy(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	getURL := path.Join(edgecloud.DBaaSClustersBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)
	mux.HandleFunc(getURL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(&edgecloud.DBaaSCluster{
			ID:     testResourceID,
			Status: DBaaSClusterHealthyStatus,
		})
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

	cluster, err := WaitDBaaSClusterStatusHealthy(context.Background(), client, testResourceID)
	assert.NoError(t, err)
	assert.NotNil(t, cluster)
	assert.Equal(t, testResourceID, cluster.ID)
}

func TestWaitDBaaSClusterStatusHealthy_NotReady(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	getURL := path.Join(edgecloud.DBaaSClustersBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)
	mux.HandleFunc(getURL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(&edgecloud.DBaaSCluster{
			ID:     testResourceID,
			Status: "PENDING_CREATE",
		})
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

	cluster, err := WaitDBaaSClusterStatusHealthy(context.Background(), client, testResourceID, 20*time.Millisecond)
	assert.ErrorIs(t, err, ErrDBaaSClusterNotReady)
	assert.Nil(t, cluster)
}

func TestCreateDBaaSClusterAndWait(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	reqBody := edgecloud.DBaaSClusterCreateRequest{
		Name: testName,
		DBMS: edgecloud.DBaaSDbmsType{
			Type:    "POSTGRESQL",
			Version: "17.5",
		},
		Flavor: "db-g2-standard-2-4-30",
		Volume: edgecloud.DBaaSVolume{
			Size: 30,
			Type: "db_standard",
		},
		Interface: edgecloud.DBaaSClusterInterface{
			NetworkID: testResourceID,
			SubnetID:  testResourceID2,
		},
	}

	createURL := path.Join(edgecloud.DBaaSClustersBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID))
	getURL := path.Join(edgecloud.DBaaSClustersBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID)

	mux.HandleFunc(createURL, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			got := &edgecloud.DBaaSClusterCreateRequest{}
			if err := json.NewDecoder(r.Body).Decode(got); err != nil {
				t.Fatalf("failed to decode request body: %v", err)
			}
			assert.Equal(t, reqBody, *got)

			resp, err := json.Marshal(&edgecloud.TaskResponse{Tasks: []string{testResourceID2}})
			if err != nil {
				t.Fatalf("failed to marshal JSON: %v", err)
			}
			_, _ = fmt.Fprint(w, string(resp))

			return
		}

		resp, err := json.Marshal([]edgecloud.DBaaSCluster{{ID: testResourceID, Name: testName}})
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	mux.HandleFunc(getURL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(&edgecloud.DBaaSCluster{
			ID:            testResourceID,
			Name:          testName,
			CreatorTaskID: testResourceID2,
			Status:        DBaaSClusterHealthyStatus,
		})
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

	cluster, err := CreateDBaaSClusterAndWait(context.Background(), client, reqBody)
	assert.NoError(t, err)
	assert.NotNil(t, cluster)
	assert.Equal(t, testResourceID, cluster.ID)
	assert.Equal(t, DBaaSClusterHealthyStatus, cluster.Status)
}

func TestCreateDBaaSClusterAndWait_EmptyTaskIDs(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	createURL := path.Join(edgecloud.DBaaSClustersBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID))
	mux.HandleFunc(createURL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(&edgecloud.TaskResponse{})
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

	cluster, err := CreateDBaaSClusterAndWait(context.Background(), client, edgecloud.DBaaSClusterCreateRequest{Name: testName})
	assert.EqualError(t, err, "DBaaS cluster create returned no task IDs")
	assert.Nil(t, cluster)
}

func TestWaitDBaaSBackupStatusFinished(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	backupID := testResourceID
	getURL := path.Join(edgecloud.DBaaSBackupsBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID), backupID)
	mux.HandleFunc(getURL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(&edgecloud.DBaaSBackup{
			ID:     backupID,
			Status: "FINISHED",
		})
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

	backup, err := WaitDBaaSBackupStatusFinished(context.Background(), client, backupID)
	assert.NoError(t, err)
	assert.NotNil(t, backup)
	assert.Equal(t, backupID, backup.ID)
}

func TestWaitDBaaSBackupStatusFinished_Error(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	backupID := testResourceID
	getURL := path.Join(edgecloud.DBaaSBackupsBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID), backupID)
	mux.HandleFunc(getURL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(&edgecloud.DBaaSBackup{
			ID:     backupID,
			Status: "ERROR",
		})
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

	backup, err := WaitDBaaSBackupStatusFinished(context.Background(), client, backupID, 20*time.Millisecond)
	assert.EqualError(t, err, "DBaaS backup entered ERROR status")
	assert.Nil(t, backup)
}

func TestWaitDBaaSBackupStatusFinished_Timeout(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	backupID := testResourceID
	getURL := path.Join(edgecloud.DBaaSBackupsBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID), backupID)
	mux.HandleFunc(getURL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(&edgecloud.DBaaSBackup{
			ID:     backupID,
			Status: "BUILDING",
		})
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

	backup, err := WaitDBaaSBackupStatusFinished(context.Background(), client, backupID, time.Millisecond)
	assert.ErrorIs(t, err, ErrDBaaSBackupNotReady)
	assert.Nil(t, backup)
}

func TestCreateDBaaSBackupAndWait(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	reqBody := edgecloud.DBaaSBackupCreateRequest{
		Name:      "my-backup",
		ClusterID: testResourceID,
	}

	createURL := path.Join(edgecloud.DBaaSBackupsBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID))
	getURL := path.Join(edgecloud.DBaaSBackupsBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID), testResourceID2)

	mux.HandleFunc(createURL, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			got := &edgecloud.DBaaSBackupCreateRequest{}
			if err := json.NewDecoder(r.Body).Decode(got); err != nil {
				t.Fatalf("failed to decode request body: %v", err)
			}
			assert.Equal(t, reqBody, *got)

			resp, err := json.Marshal(&edgecloud.TaskResponse{Tasks: []string{testResourceID2}})
			if err != nil {
				t.Fatalf("failed to marshal JSON: %v", err)
			}
			_, _ = fmt.Fprint(w, string(resp))

			return
		}

		// GET on createURL = BackupsList
		backups := []edgecloud.DBaaSBackup{
			{
				ID:            testResourceID,
				CreatorTaskID: testResourceID,
				Status:        "BUILDING",
			},
			{
				ID:            testResourceID2,
				CreatorTaskID: testResourceID2,
				Status:        "BUILDING",
			},
		}
		resp, err := json.Marshal(backups)
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	mux.HandleFunc(getURL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(&edgecloud.DBaaSBackup{
			ID:     testResourceID2,
			Status: "FINISHED",
		})
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

	backup, err := CreateDBaaSBackupAndWait(context.Background(), client, reqBody)
	assert.NoError(t, err)
	assert.NotNil(t, backup)
	assert.Equal(t, testResourceID2, backup.ID)
	assert.Equal(t, "FINISHED", backup.Status)
}

func TestCreateDBaaSBackupAndWait_EmptyTaskIDs(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	createURL := path.Join(edgecloud.DBaaSBackupsBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID))
	mux.HandleFunc(createURL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(&edgecloud.TaskResponse{})
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

	backup, err := CreateDBaaSBackupAndWait(context.Background(), client, edgecloud.DBaaSBackupCreateRequest{Name: "test"})
	assert.EqualError(t, err, "DBaaS backup create returned no task IDs")
	assert.Nil(t, backup)
}

func TestWaitDBaaSBackupByCreatorTaskID_Timeout(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	listURL := path.Join(edgecloud.DBaaSBackupsBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID))
	mux.HandleFunc(listURL, func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"results":[]}`)
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	backup, err := WaitDBaaSBackupByCreatorTaskID(context.Background(), client, testResourceID, time.Millisecond)
	assert.ErrorIs(t, err, ErrDBaaSBackupNotReady)
	assert.Nil(t, backup)
}

func TestWaitDBaaSBackupByCreatorTaskID_Found(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	listURL := path.Join(edgecloud.DBaaSBackupsBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID))
	mux.HandleFunc(listURL, func(w http.ResponseWriter, r *http.Request) {
		backups := []edgecloud.DBaaSBackup{
			{ID: testResourceID, CreatorTaskID: testResourceID2, Status: "BUILDING"},
			{ID: testResourceID2, CreatorTaskID: testResourceID, Status: "FINISHED"},
		}
		resp, err := json.Marshal(backups)
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

	backup, err := WaitDBaaSBackupByCreatorTaskID(context.Background(), client, testResourceID)
	assert.NoError(t, err)
	assert.NotNil(t, backup)
	assert.Equal(t, testResourceID2, backup.ID)
}
