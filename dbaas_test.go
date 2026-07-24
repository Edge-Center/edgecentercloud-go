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

const (
	dbaasClusterName  = "my-db-cluster"
	dbaasFlavor       = "db-g2-standard-2-4-30"
	dbaasNetworkID    = "f0d19cec-5c3f-4853-886e-304915960ff4"
	dbaasSubnetID     = "rtb19cec-5c3f-4853-886e-3045fk9g5kg9"
	dbaasClusterID    = "123e4567-e89b-12d3-a456-426614174000"
	dbaasUserName     = "user_name"
	dbaasDatabaseName = "user_analytics"
	testBackupID      = "e13bd88d-79c8-4871-a9c2-a4baca741d0f"
)

func TestDBaaSServiceOp_ClusterCreate(t *testing.T) {
	setup()
	defer teardown()

	request := DBaaSClusterCreateRequest{
		Name:             dbaasClusterName,
		Description:      "Database Cluster",
		HighAvailability: true,
		DBMS: DBaaSDbmsType{
			Type:    "POSTGRESQL",
			Version: "17.5",
		},
		Flavor: dbaasFlavor,
		Volume: DBaaSVolume{
			Size: 30,
			Type: "db_standard",
		},
		Interface: DBaaSClusterInterface{
			NetworkID: dbaasNetworkID,
			SubnetID:  dbaasSubnetID,
		},
	}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(DBaaSClustersBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := &DBaaSClusterCreateRequest{}
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, *reqBody)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.DBaaS.ClusterCreate(ctx, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestDBaaSServiceOp_ClustersList(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []DBaaSCluster{{ID: dbaasClusterID}}
	URL := path.Join(DBaaSClustersBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, err := json.Marshal(expectedResp)
		if err != nil {
			t.Errorf("failed to marshal response: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.DBaaS.ClustersList(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestDBaaSServiceOp_ClusterGet(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &DBaaSCluster{
		ID:   dbaasClusterID,
		Name: "my-db-cluster",
	}
	clusterID := dbaasClusterID
	URL := path.Join(DBaaSClustersBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID), clusterID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(expectedResp)
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.DBaaS.ClusterGet(ctx, clusterID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestDBaaSServiceOp_ClusterDelete(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	clusterID := dbaasClusterID
	URL := path.Join(DBaaSClustersBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID), clusterID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(expectedResp)
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.DBaaS.ClusterDelete(ctx, clusterID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestDBaaSServiceOp_ClusterUpdate(t *testing.T) {
	setup()
	defer teardown()

	request := DBaaSClusterUpdateRequest{
		Name: "updated-cluster",
	}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	clusterID := dbaasClusterID
	URL := path.Join(DBaaSClustersBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID), clusterID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		reqBody := &DBaaSClusterUpdateRequest{}
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, *reqBody)
		resp, _ := json.Marshal(expectedResp)
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.DBaaS.ClusterUpdate(ctx, clusterID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestDBaaSServiceOp_UsersList(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []DBaaSUser{
		{
			Name: dbaasUserName,
			Databases: []DBaaSUserDatabase{
				{Name: dbaasDatabaseName},
			},
		},
	}
	clusterID := dbaasClusterID
	URL := path.Join(DBaaSClustersBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID), clusterID, "users")

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(expectedResp)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.DBaaS.UsersList(ctx, clusterID, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestDBaaSServiceOp_UserCreate(t *testing.T) {
	setup()
	defer teardown()

	request := DBaaSUserCreateRequest{
		Name:     dbaasUserName,
		Password: "Qwerty!123456",
		Databases: []DBaaSUserDatabase{
			{Name: dbaasDatabaseName},
		},
	}
	expectedResp := &DBaaSUser{
		Name: dbaasUserName,
		Databases: []DBaaSUserDatabase{
			{Name: dbaasDatabaseName},
		},
	}
	clusterID := dbaasClusterID
	URL := path.Join(DBaaSClustersBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID), clusterID, "users")

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := &DBaaSUserCreateRequest{}
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, *reqBody)
		resp, _ := json.Marshal(expectedResp)
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.DBaaS.UserCreate(ctx, clusterID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestDBaaSServiceOp_UserGet(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &DBaaSUser{
		Name: dbaasUserName,
		Databases: []DBaaSUserDatabase{
			{Name: dbaasDatabaseName},
		},
	}
	clusterID := dbaasClusterID
	username := dbaasUserName
	URL := path.Join(DBaaSClustersBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID), clusterID, "users", username)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(expectedResp)
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.DBaaS.UserGet(ctx, clusterID, username)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestDBaaSServiceOp_UserUpdate(t *testing.T) {
	setup()
	defer teardown()

	request := DBaaSUserUpdateRequest{
		Password: "Qwerty!123456",
	}
	clusterID := dbaasClusterID
	username := dbaasUserName
	URL := path.Join(DBaaSClustersBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID), clusterID, "users", username)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		reqBody := &DBaaSUserUpdateRequest{}
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, *reqBody)
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.DBaaS.UserUpdate(ctx, clusterID, username, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 204)
}

func TestDBaaSServiceOp_UserDelete(t *testing.T) {
	setup()
	defer teardown()

	clusterID := dbaasClusterID
	username := dbaasUserName
	URL := path.Join(DBaaSClustersBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID), clusterID, "users", username)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.DBaaS.UserDelete(ctx, clusterID, username)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 204)
}

func TestDBaaSServiceOp_UserGrantAccess(t *testing.T) {
	setup()
	defer teardown()

	clusterID := dbaasClusterID
	username := dbaasUserName
	database := dbaasDatabaseName
	URL := path.Join(DBaaSClustersBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID), clusterID, "users", username, "databases", database)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.DBaaS.UserGrantAccess(ctx, clusterID, username, database)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 204)
}

func TestDBaaSServiceOp_UserRevokeAccess(t *testing.T) {
	setup()
	defer teardown()

	clusterID := dbaasClusterID
	username := dbaasUserName
	database := dbaasDatabaseName
	URL := path.Join(DBaaSClustersBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID), clusterID, "users", username, "databases", database)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.DBaaS.UserRevokeAccess(ctx, clusterID, username, database)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 204)
}

func TestDBaaSServiceOp_DatabasesList(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []DBaaSDatabase{{Name: dbaasDatabaseName}}
	clusterID := dbaasClusterID
	URL := path.Join(DBaaSClustersBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID), clusterID, "databases")

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(expectedResp)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.DBaaS.DatabasesList(ctx, clusterID, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestDBaaSServiceOp_DatabaseCreate(t *testing.T) {
	setup()
	defer teardown()

	request := DBaaSDatabaseCreateRequest{
		Name:     "my-database",
		Encoding: "UTF8",
		Locale:   "en_US.utf8",
	}
	expectedResp := &DBaaSDatabase{Name: "my-database"}
	clusterID := dbaasClusterID
	URL := path.Join(DBaaSClustersBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID), clusterID, "databases")

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := &DBaaSDatabaseCreateRequest{}
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, *reqBody)
		resp, _ := json.Marshal(expectedResp)
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.DBaaS.DatabaseCreate(ctx, clusterID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestDBaaSServiceOp_DatabaseDelete(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := &DBaaSDatabase{Name: dbaasDatabaseName}
	clusterID := dbaasClusterID
	databaseName := dbaasDatabaseName
	URL := path.Join(DBaaSClustersBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID), clusterID, "databases", databaseName)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(expectedResp)
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.DBaaS.DatabaseDelete(ctx, clusterID, databaseName)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestDBaaSServiceOp_DbmsList(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []DBaaSDbms{
		{
			ID:      1,
			Type:    "POSTGRESQL",
			Version: "17.5",
		},
	}
	URL := path.Join(DBaaSDbmsBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(expectedResp)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.DBaaS.DbmsList(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestDBaaSServiceOp_BackupCreate(t *testing.T) {
	setup()
	defer teardown()

	request := DBaaSBackupCreateRequest{
		Name:      "my-backup",
		ClusterID: dbaasClusterID,
	}
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(DBaaSBackupsBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		reqBody := &DBaaSBackupCreateRequest{}
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, *reqBody)
		resp, _ := json.Marshal(expectedResp)
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.DBaaS.BackupCreate(ctx, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestDBaaSServiceOp_BackupsList(t *testing.T) {
	setup()
	defer teardown()

	expectedResp := []DBaaSBackup{{ID: testBackupID, Name: "my-backup"}}
	URL := path.Join(DBaaSBackupsBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(expectedResp)
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	respActual, resp, err := client.DBaaS.BackupsList(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestDBaaSServiceOp_BackupGet(t *testing.T) {
	setup()
	defer teardown()

	backupID := testBackupID
	expectedResp := &DBaaSBackup{ID: backupID, Name: "my-backup"}
	URL := path.Join(DBaaSBackupsBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID), backupID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		resp, _ := json.Marshal(expectedResp)
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.DBaaS.BackupGet(ctx, backupID, false)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestDBaaSServiceOp_BackupUpdate(t *testing.T) {
	setup()
	defer teardown()

	newName := "new-backup-name"
	request := DBaaSBackupUpdateRequest{
		Name: &newName,
	}
	backupID := testBackupID
	expectedResp := &DBaaSBackup{ID: backupID, Name: "new-backup-name"}
	URL := path.Join(DBaaSBackupsBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID), backupID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		reqBody := &DBaaSBackupUpdateRequest{}
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, request, *reqBody)
		resp, _ := json.Marshal(expectedResp)
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.DBaaS.BackupUpdate(ctx, backupID, request)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}

func TestDBaaSServiceOp_BackupDelete(t *testing.T) {
	setup()
	defer teardown()

	backupID := testBackupID
	expectedResp := &TaskResponse{Tasks: []string{taskID}}
	URL := path.Join(DBaaSBackupsBasePathV3, strconv.Itoa(projectID), strconv.Itoa(regionID), backupID)

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		resp, _ := json.Marshal(expectedResp)
		_, _ = fmt.Fprint(w, string(resp))
	})

	respActual, resp, err := client.DBaaS.BackupDelete(ctx, backupID)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	require.Equal(t, respActual, expectedResp)
}
