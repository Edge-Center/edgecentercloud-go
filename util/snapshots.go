package util

import (
	"context"
	"errors"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

const SnapshotReadyStatus = "available"

var (
	ErrSnapshotsNotFound = errors.New("no Snapshots were found for the specified search criteria")
	ErrSnapshotNotReady  = errors.New("snapshot failed to be ready within the allocated time")
)

func SnapshotsListByStatusAndVolumeID(ctx context.Context, client *edgecloud.Client, status, volumeID string) ([]edgecloud.Snapshot, error) {
	var snapshots []edgecloud.Snapshot

	snapList, _, err := client.Snapshots.List(ctx, &edgecloud.SnapshotListOptions{VolumeID: volumeID})
	if err != nil {
		return nil, err
	}

	for _, snap := range snapList {
		if snap.Status == status {
			snapshots = append(snapshots, snap)
		}
	}

	if len(snapshots) == 0 {
		return nil, ErrSnapshotsNotFound
	}

	return snapshots, nil
}

func SnapshotsListByNameAndVolumeID(ctx context.Context, client *edgecloud.Client, name, volumeID string) ([]edgecloud.Snapshot, error) {
	var snapshots []edgecloud.Snapshot

	snapList, _, err := client.Snapshots.List(ctx, &edgecloud.SnapshotListOptions{VolumeID: volumeID})
	if err != nil {
		return nil, err
	}

	for _, snap := range snapList {
		if snap.Name == name {
			snapshots = append(snapshots, snap)
		}
	}

	if len(snapshots) == 0 {
		return nil, ErrSnapshotsNotFound
	}

	return snapshots, nil
}

func WaitSnapshotStatusReady(ctx context.Context, client *edgecloud.Client, snapshotID string, attempts *uint) error {
	return WithRetry(
		func() error {
			snapshot, _, err := client.Snapshots.Get(ctx, snapshotID)
			if err != nil {
				return err
			}

			if snapshot.Status == SnapshotReadyStatus {
				return nil
			}

			return ErrSnapshotNotReady
		},
		attempts,
	)
}
