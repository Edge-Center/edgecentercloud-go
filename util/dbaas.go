package util

import (
	"context"
	"errors"
	"fmt"
	"time"

	edgecloud "github.com/Edge-Center/edgecentercloud-go/v2"
)

const (
	DBaaSClusterHealthyStatus = "HEALTHY"
	DBaaSBackupFinishedStatus = "FINISHED"
	dbaasClusterPollInterval  = 5 * time.Second
	dbaasBackupPollInterval   = 5 * time.Second
	dbaasBackupPageSize       = 100
)

var (
	ErrDBaaSClustersNotFound = errors.New("no DBaaS clusters were found for the specified search criteria")
	ErrDBaaSClusterNotReady  = errors.New("DBaaS cluster failed to become healthy within the allocated time")
	ErrDBaaSBackupNotReady   = errors.New("DBaaS backup failed to become ready within the allocated time")
)

func CreateDBaaSClusterAndWait(ctx context.Context, client *edgecloud.Client, req edgecloud.DBaaSClusterCreateRequest, timeouts ...time.Duration) (*edgecloud.DBaaSCluster, error) {
	task, _, err := client.DBaaS.ClusterCreate(ctx, req)
	if err != nil {
		return nil, err
	}

	if task == nil || len(task.Tasks) == 0 {
		return nil, fmt.Errorf("DBaaS cluster create returned no task IDs")
	}

	cluster, err := WaitDBaaSClusterByCreatorTaskID(ctx, client, req.Name, task.Tasks[0], timeouts...)
	if err != nil {
		return nil, err
	}

	cluster, err = WaitDBaaSClusterStatusHealthy(ctx, client, cluster.ID, timeouts...)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

func DBaaSClusterGetByCreatorTaskID(ctx context.Context, client *edgecloud.Client, name string, creatorTaskID string) (*edgecloud.DBaaSCluster, error) {
	clusters, _, err := client.DBaaS.ClustersList(ctx, nil)
	if err != nil {
		return nil, err
	}

	for _, item := range clusters {
		if item.Name != name {
			continue
		}

		cluster, _, err := client.DBaaS.ClusterGet(ctx, item.ID)
		if err != nil {
			return nil, err
		}

		if cluster.CreatorTaskID == creatorTaskID {
			return cluster, nil
		}
	}

	return nil, ErrDBaaSClustersNotFound
}

func WaitDBaaSClusterByCreatorTaskID(ctx context.Context, client *edgecloud.Client, name string, creatorTaskID string, timeouts ...time.Duration) (*edgecloud.DBaaSCluster, error) {
	timeout := defaultTimeout
	if len(timeouts) > 0 {
		timeout = timeouts[0]
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(dbaasClusterPollInterval)
	defer ticker.Stop()

	for {
		cluster, err := DBaaSClusterGetByCreatorTaskID(ctx, client, name, creatorTaskID)
		if err == nil {
			return cluster, nil
		}

		if !errors.Is(err, ErrDBaaSClustersNotFound) {
			return nil, err
		}

		select {
		case <-ctx.Done():
			return nil, edgecloud.NewArgError("task error", errTaskWaitTimeout.Error())
		case <-ticker.C:
		}
	}
}

func WaitDBaaSClusterStatusHealthy(ctx context.Context, client *edgecloud.Client, clusterID string, timeouts ...time.Duration) (*edgecloud.DBaaSCluster, error) {
	timeout := defaultTimeout
	if len(timeouts) > 0 {
		timeout = timeouts[0]
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(dbaasClusterPollInterval)
	defer ticker.Stop()

	for {
		cluster, _, err := client.DBaaS.ClusterGet(ctx, clusterID)
		if err != nil {
			return nil, err
		}

		if cluster.Status == DBaaSClusterHealthyStatus {
			return cluster, nil
		}

		select {
		case <-ctx.Done():
			return nil, ErrDBaaSClusterNotReady
		case <-ticker.C:
		}
	}
}

func CreateDBaaSBackupAndWait(ctx context.Context, client *edgecloud.Client, req edgecloud.DBaaSBackupCreateRequest, timeouts ...time.Duration) (*edgecloud.DBaaSBackup, error) {
	ctx, cancel := contextWithTimeout(ctx, timeouts...)
	defer cancel()

	task, _, err := client.DBaaS.BackupCreate(ctx, req)
	if err != nil {
		return nil, err
	}

	if task == nil || len(task.Tasks) == 0 {
		return nil, fmt.Errorf("DBaaS backup create returned no task IDs")
	}

	backup, err := waitDBaaSBackupByCreatorTaskID(ctx, client, task.Tasks[0])
	if err != nil {
		return nil, err
	}

	return waitDBaaSBackupStatusFinished(ctx, client, backup.ID)
}

func WaitDBaaSBackupByCreatorTaskID(ctx context.Context, client *edgecloud.Client, creatorTaskID string, timeouts ...time.Duration) (*edgecloud.DBaaSBackup, error) {
	ctx, cancel := contextWithTimeout(ctx, timeouts...)
	defer cancel()

	return waitDBaaSBackupByCreatorTaskID(ctx, client, creatorTaskID)
}

func waitDBaaSBackupByCreatorTaskID(ctx context.Context, client *edgecloud.Client, creatorTaskID string) (*edgecloud.DBaaSBackup, error) {
	ticker := time.NewTicker(dbaasBackupPollInterval)
	defer ticker.Stop()

	for {
		for offset := 0; ; {
			page, _, err := client.DBaaS.BackupsListPage(ctx, &edgecloud.DBaaSBackupListOptions{
				Limit:  dbaasBackupPageSize,
				Offset: offset,
			})
			if err != nil {
				if ctx.Err() != nil {
					return nil, ErrDBaaSBackupNotReady
				}
				return nil, err
			}

			for i := range page.Results {
				if page.Results[i].CreatorTaskID == creatorTaskID {
					return &page.Results[i], nil
				}
			}

			if len(page.Results) == 0 || len(page.Results) < dbaasBackupPageSize || (page.Count > 0 && offset+len(page.Results) >= page.Count) {
				break
			}
			offset += len(page.Results)
		}

		select {
		case <-ctx.Done():
			return nil, ErrDBaaSBackupNotReady
		case <-ticker.C:
		}
	}
}

func WaitDBaaSBackupStatusFinished(ctx context.Context, client *edgecloud.Client, backupID string, timeouts ...time.Duration) (*edgecloud.DBaaSBackup, error) {
	ctx, cancel := contextWithTimeout(ctx, timeouts...)
	defer cancel()

	return waitDBaaSBackupStatusFinished(ctx, client, backupID)
}

func waitDBaaSBackupStatusFinished(ctx context.Context, client *edgecloud.Client, backupID string) (*edgecloud.DBaaSBackup, error) {
	ticker := time.NewTicker(dbaasBackupPollInterval)
	defer ticker.Stop()

	for {
		backup, _, err := client.DBaaS.BackupGet(ctx, backupID, false)
		if err != nil {
			if ctx.Err() != nil {
				return nil, ErrDBaaSBackupNotReady
			}
			return nil, err
		}

		if backup.Status == DBaaSBackupFinishedStatus {
			return backup, nil
		}

		if backup.Status == "ERROR" {
			return nil, fmt.Errorf("DBaaS backup entered ERROR status")
		}

		select {
		case <-ctx.Done():
			return nil, ErrDBaaSBackupNotReady
		case <-ticker.C:
		}
	}
}

func contextWithTimeout(ctx context.Context, timeouts ...time.Duration) (context.Context, context.CancelFunc) {
	if len(timeouts) > 0 {
		return context.WithTimeout(ctx, timeouts[0])
	}
	if _, ok := ctx.Deadline(); ok {
		return ctx, func() {}
	}

	return context.WithTimeout(ctx, defaultTimeout)
}
