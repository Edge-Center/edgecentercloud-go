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
	dbaasClusterPollInterval  = 5 * time.Second
)

var (
	ErrDBaaSClustersNotFound = errors.New("no DBaaS clusters were found for the specified search criteria")
	ErrDBaaSClusterNotReady  = errors.New("DBaaS cluster failed to become healthy within the allocated time")
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
