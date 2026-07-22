package edgecloud

import (
	"context"
	"fmt"
	"net/http"
)

const (
	DBaaSClustersBasePathV3 = "/dbaas/v3/clusters"
	DBaaSDbmsBasePathV3     = "/dbaas/v3/dbms"
	DBaaSBackupsBasePathV3  = "/dbaas/v3/backups"
)

type DBaaSService interface {
	DBaaSClusters
	DBaaSUsers
	DBaaSDatabases
	DBaaSBackups
	DBaaSDbmsService
}

type DBaaSClusters interface {
	ClusterCreate(context.Context, DBaaSClusterCreateRequest) (*TaskResponse, *Response, error)
	ClustersList(context.Context, *DBaaSClusterListOptions) ([]DBaaSCluster, *Response, error)
	ClusterGet(context.Context, string) (*DBaaSCluster, *Response, error)
	ClusterDelete(context.Context, string) (*TaskResponse, *Response, error)
	ClusterUpdate(context.Context, string, DBaaSClusterUpdateRequest) (*TaskResponse, *Response, error)
}

type DBaaSUsers interface {
	UsersList(ctx context.Context, clusterID string, opts *DBaaSUserListOptions) ([]DBaaSUser, *Response, error)
	UserCreate(ctx context.Context, clusterID string, reqBody DBaaSUserCreateRequest) (*DBaaSUser, *Response, error)
	UserGet(ctx context.Context, clusterID string, username string) (*DBaaSUser, *Response, error)
	UserUpdate(ctx context.Context, clusterID string, username string, reqBody DBaaSUserUpdateRequest) (*Response, error)
	UserDelete(ctx context.Context, clusterID string, username string) (*Response, error)
	UserGrantAccess(ctx context.Context, clusterID string, username string, database string) (*Response, error)
	UserRevokeAccess(ctx context.Context, clusterID string, username string, database string) (*Response, error)
}

type DBaaSDatabases interface {
	DatabasesList(ctx context.Context, clusterID string, opts *DBaaSDatabaseListOptions) ([]DBaaSDatabase, *Response, error)
	DatabaseCreate(ctx context.Context, clusterID string, reqBody DBaaSDatabaseCreateRequest) (*DBaaSDatabase, *Response, error)
	DatabaseDelete(ctx context.Context, clusterID string, databaseName string) (*DBaaSDatabase, *Response, error)
}

type DBaaSBackups interface {
	BackupCreate(ctx context.Context, reqBody DBaaSBackupCreateRequest) (*TaskResponse, *Response, error)
	BackupsList(ctx context.Context, opts *DBaaSBackupListOptions) ([]DBaaSBackup, *Response, error)
	BackupsListPage(ctx context.Context, opts *DBaaSBackupListOptions) (*DBaaSBackupsPage, *Response, error)
	BackupGet(ctx context.Context, backupID string, includePrices bool) (*DBaaSBackup, *Response, error)
	BackupUpdate(ctx context.Context, backupID string, reqBody DBaaSBackupUpdateRequest) (*DBaaSBackup, *Response, error)
	BackupDelete(ctx context.Context, backupID string) (*TaskResponse, *Response, error)
}

type DBaaSDbmsService interface {
	DbmsList(ctx context.Context, opts *DBaaSDbmsListOptions) ([]DBaaSDbms, *Response, error)
}

type DBaaSServiceOp struct {
	client *Client
}

var _ DBaaSService = &DBaaSServiceOp{}

type DBaaSCluster struct {
	ID          string         `json:"id"`
	ProjectID   int            `json:"project_id"`
	RegionID    int            `json:"region_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	DBMS        *DBaaSDbmsType `json:"dbms"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`

	// Get-specific (отсутствуют в List)
	TaskID           string                 `json:"task_id,omitempty"`
	CreatorTaskID    string                 `json:"creator_task_id,omitempty"`
	HighAvailability bool                   `json:"high_availability,omitempty"`
	Flavor           string                 `json:"flavor,omitempty"`
	Volume           *DBaaSVolume           `json:"volume,omitempty"`
	Interface        *DBaaSClusterInterface `json:"interface,omitempty"`
	Connection       *DBaaSConnection       `json:"connection,omitempty"`
}

type DBaaSClusterCreateRequest struct {
	Name             string                `json:"name"`
	Description      string                `json:"description,omitempty"`
	HighAvailability bool                  `json:"high_availability"`
	DBMS             DBaaSDbmsType         `json:"dbms"`
	Flavor           string                `json:"flavor"`
	Volume           DBaaSVolume           `json:"volume"`
	Interface        DBaaSClusterInterface `json:"interface"`
}

type DBaaSClusterUpdateRequest struct {
	Name        string       `json:"name,omitempty"`
	Description string       `json:"description,omitempty"`
	Flavor      string       `json:"flavor,omitempty"`
	Volume      *DBaaSVolume `json:"volume,omitempty"`
}

type DBaaSDbmsType struct {
	Type    string `json:"type"`
	Version string `json:"version"`
}

type DBaaSVolume struct {
	Size int        `json:"size"`
	Type VolumeType `json:"type"`
}

type DBaaSClusterInterface struct {
	NetworkID string `json:"network_id"`
	SubnetID  string `json:"subnet_id"`
}

type DBaaSClusterListOptions struct {
	Limit  int `url:"limit,omitempty"`
	Offset int `url:"offset,omitempty"`
}

type dbaasClustersRoot struct {
	Count    int
	Clusters []DBaaSCluster `json:"results"`
}

type DBaaSConnection struct {
	Method string `json:"method"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
}

type DBaaSUserListOptions struct {
	Limit  int `url:"limit,omitempty"`
	Offset int `url:"offset,omitempty"`
}

type dbaasUserRoot struct {
	Count int
	Users []DBaaSUser `json:"results"`
}

type DBaaSUserDatabase struct {
	Name string `json:"name"`
}

type DBaaSUser struct {
	Name      string              `json:"name"`
	Databases []DBaaSUserDatabase `json:"databases"`
}

type DBaaSUserCreateRequest struct {
	Name      string              `json:"name"`
	Password  string              `json:"password"`
	Databases []DBaaSUserDatabase `json:"databases,omitempty"`
}

type DBaaSUserUpdateRequest struct {
	Password string `json:"password"`
}

type DBaaSDatabaseListOptions struct {
	Limit  int `url:"limit,omitempty"`
	Offset int `url:"offset,omitempty"`
}

type dbaasDatabaseRoot struct {
	Count     int
	Databases []DBaaSDatabase `json:"results"`
}

type DBaaSDatabase struct {
	Name string `json:"name"`
}

type DBaaSDatabaseCreateRequest struct {
	Name     string `json:"name"`
	Encoding string `json:"encoding,omitempty"`
	Locale   string `json:"locale,omitempty"`
}

type DBaaSDbmsListOptions struct {
	Limit  int `url:"limit,omitempty"`
	Offset int `url:"offset,omitempty"`
}

type dbaasDbmsRoot struct {
	Count int
	Dbms  []DBaaSDbms `json:"results"`
}

type DBaaSDbms struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Version string `json:"version"`
}

type DBaaSBackup struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	BackupType    string         `json:"backup_type"`
	ClusterID     string         `json:"cluster_id"`
	ParentID      string         `json:"parent_id"`
	Status        string         `json:"status"`
	Size          float64        `json:"size"`
	IsService     bool           `json:"is_service"`
	HasChild      bool           `json:"has_child,omitempty"`
	DBMS          *DBaaSDbmsType `json:"dbms"`
	CreatedAt     string         `json:"created_at"`
	UpdatedAt     string         `json:"updated_at"`
	FinishedAt    string         `json:"finished_at"`
	TaskID        string         `json:"task_id"`
	CreatorTaskID string         `json:"creator_task_id"`
}

type DBaaSBackupCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ClusterID   string `json:"cluster_id"`
	ParentID    string `json:"parent_id,omitempty"`
}

type DBaaSBackupUpdateRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type DBaaSBackupListOptions struct {
	ClusterID     string `url:"cluster_id,omitempty"`
	BackupType    string `url:"backup_type,omitempty"`
	IsService     *bool  `url:"is_service,omitempty"`
	CreatedFrom   string `url:"created_from,omitempty"`
	CreatedTo     string `url:"created_to,omitempty"`
	DbmsID        int    `url:"dbms_id,omitempty"`
	Search        string `url:"search,omitempty"`
	IncludePrices bool   `url:"include_prices,omitempty"`
	Limit         int    `url:"limit,omitempty"`
	Offset        int    `url:"offset,omitempty"`
}

type DBaaSBackupsPage struct {
	Count   int           `json:"count"`
	Results []DBaaSBackup `json:"results"`
}

func (s *DBaaSServiceOp) ClusterCreate(ctx context.Context, reqBody DBaaSClusterCreateRequest) (*TaskResponse, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(DBaaSClustersBasePathV3)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	tasks := new(TaskResponse)
	resp, err := s.client.Do(ctx, req, tasks)
	if err != nil {
		return nil, resp, err
	}

	return tasks, resp, err
}

func (s *DBaaSServiceOp) ClustersList(ctx context.Context, opts *DBaaSClusterListOptions) ([]DBaaSCluster, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(DBaaSClustersBasePathV3)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(dbaasClustersRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Clusters, resp, err
}

func (s *DBaaSServiceOp) ClusterGet(ctx context.Context, clusterID string) (*DBaaSCluster, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(DBaaSClustersBasePathV3), clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	cluster := new(DBaaSCluster)
	resp, err := s.client.Do(ctx, req, cluster)
	if err != nil {
		return nil, resp, err
	}

	return cluster, resp, err
}

func (s *DBaaSServiceOp) ClusterDelete(ctx context.Context, clusterID string) (*TaskResponse, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(DBaaSClustersBasePathV3), clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	tasks := new(TaskResponse)
	resp, err := s.client.Do(ctx, req, tasks)
	if err != nil {
		return nil, resp, err
	}

	return tasks, resp, err
}

func (s *DBaaSServiceOp) ClusterUpdate(ctx context.Context, clusterID string, reqBody DBaaSClusterUpdateRequest) (*TaskResponse, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(DBaaSClustersBasePathV3), clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	tasks := new(TaskResponse)
	resp, err := s.client.Do(ctx, req, tasks)
	if err != nil {
		return nil, resp, err
	}

	return tasks, resp, err
}

func (s *DBaaSServiceOp) UsersList(ctx context.Context, clusterID string, opts *DBaaSUserListOptions) ([]DBaaSUser, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/users", s.client.addProjectRegionPath(DBaaSClustersBasePathV3), clusterID)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(dbaasUserRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Users, resp, err
}

func (s *DBaaSServiceOp) UserCreate(ctx context.Context, clusterID string, reqBody DBaaSUserCreateRequest) (*DBaaSUser, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/users", s.client.addProjectRegionPath(DBaaSClustersBasePathV3), clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	user := new(DBaaSUser)
	resp, err := s.client.Do(ctx, req, user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, err
}

func (s *DBaaSServiceOp) UserGet(ctx context.Context, clusterID, username string) (*DBaaSUser, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/users/%s", s.client.addProjectRegionPath(DBaaSClustersBasePathV3), clusterID, username)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	user := new(DBaaSUser)
	resp, err := s.client.Do(ctx, req, user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, err
}

func (s *DBaaSServiceOp) UserUpdate(ctx context.Context, clusterID, username string, reqBody DBaaSUserUpdateRequest) (*Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	path := fmt.Sprintf("%s/%s/users/%s", s.client.addProjectRegionPath(DBaaSClustersBasePathV3), clusterID, username)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}

func (s *DBaaSServiceOp) UserDelete(ctx context.Context, clusterID, username string) (*Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	path := fmt.Sprintf("%s/%s/users/%s", s.client.addProjectRegionPath(DBaaSClustersBasePathV3), clusterID, username)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}

func (s *DBaaSServiceOp) UserGrantAccess(ctx context.Context, clusterID, username, database string) (*Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	path := fmt.Sprintf("%s/%s/users/%s/databases/%s", s.client.addProjectRegionPath(DBaaSClustersBasePathV3), clusterID, username, database)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}

func (s *DBaaSServiceOp) UserRevokeAccess(ctx context.Context, clusterID, username, database string) (*Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return resp, err
	}

	path := fmt.Sprintf("%s/%s/users/%s/databases/%s", s.client.addProjectRegionPath(DBaaSClustersBasePathV3), clusterID, username, database)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}

func (s *DBaaSServiceOp) DatabasesList(ctx context.Context, clusterID string, opts *DBaaSDatabaseListOptions) ([]DBaaSDatabase, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/databases", s.client.addProjectRegionPath(DBaaSClustersBasePathV3), clusterID)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(dbaasDatabaseRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Databases, resp, err
}

func (s *DBaaSServiceOp) DatabaseCreate(ctx context.Context, clusterID string, reqBody DBaaSDatabaseCreateRequest) (*DBaaSDatabase, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/databases", s.client.addProjectRegionPath(DBaaSClustersBasePathV3), clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	database := new(DBaaSDatabase)
	resp, err := s.client.Do(ctx, req, database)
	if err != nil {
		return nil, resp, err
	}

	return database, resp, err
}

func (s *DBaaSServiceOp) DatabaseDelete(ctx context.Context, clusterID, databaseName string) (*DBaaSDatabase, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s/databases/%s", s.client.addProjectRegionPath(DBaaSClustersBasePathV3), clusterID, databaseName)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	database := new(DBaaSDatabase)
	resp, err := s.client.Do(ctx, req, database)
	if err != nil {
		return nil, resp, err
	}

	return database, resp, err
}

func (s *DBaaSServiceOp) DbmsList(ctx context.Context, opts *DBaaSDbmsListOptions) ([]DBaaSDbms, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(DBaaSDbmsBasePathV3)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(dbaasDbmsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Dbms, resp, err
}

func (s *DBaaSServiceOp) BackupCreate(ctx context.Context, reqBody DBaaSBackupCreateRequest) (*TaskResponse, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(DBaaSBackupsBasePathV3)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	tasks := new(TaskResponse)
	resp, err := s.client.Do(ctx, req, tasks)
	if err != nil {
		return nil, resp, err
	}

	return tasks, resp, err
}

func (s *DBaaSServiceOp) BackupsList(ctx context.Context, opts *DBaaSBackupListOptions) ([]DBaaSBackup, *Response, error) {
	page, resp, err := s.BackupsListPage(ctx, opts)
	if err != nil {
		return nil, resp, err
	}

	return page.Results, resp, nil
}

func (s *DBaaSServiceOp) BackupsListPage(ctx context.Context, opts *DBaaSBackupListOptions) (*DBaaSBackupsPage, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := s.client.addProjectRegionPath(DBaaSBackupsBasePathV3)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	page := new(DBaaSBackupsPage)
	resp, err := s.client.Do(ctx, req, page)
	if err != nil {
		return nil, resp, err
	}

	return page, resp, nil
}

func (s *DBaaSServiceOp) BackupGet(ctx context.Context, backupID string, includePrices bool) (*DBaaSBackup, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(DBaaSBackupsBasePathV3), backupID)
	if includePrices {
		path += "?include_prices=true"
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	backup := new(DBaaSBackup)
	resp, err := s.client.Do(ctx, req, backup)
	if err != nil {
		return nil, resp, err
	}

	return backup, resp, err
}

func (s *DBaaSServiceOp) BackupUpdate(ctx context.Context, backupID string, reqBody DBaaSBackupUpdateRequest) (*DBaaSBackup, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(DBaaSBackupsBasePathV3), backupID)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	backup := new(DBaaSBackup)
	resp, err := s.client.Do(ctx, req, backup)
	if err != nil {
		return nil, resp, err
	}

	return backup, resp, err
}

func (s *DBaaSServiceOp) BackupDelete(ctx context.Context, backupID string) (*TaskResponse, *Response, error) {
	if resp, err := s.client.Validate(); err != nil {
		return nil, resp, err
	}

	path := fmt.Sprintf("%s/%s", s.client.addProjectRegionPath(DBaaSBackupsBasePathV3), backupID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	tasks := new(TaskResponse)
	resp, err := s.client.Do(ctx, req, tasks)
	if err != nil {
		return nil, resp, err
	}

	return tasks, resp, err
}
