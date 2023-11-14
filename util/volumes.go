package util

import (
	"context"
	"errors"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

var (
	ErrVolumesNotFound    = errors.New("no Volumes were found for the specified search criteria")
	ErrVolumesNotAttached = errors.New("volume failed to be attached within the allocated time")
)

func VolumesListByName(ctx context.Context, client *edgecloud.Client, name string) ([]edgecloud.Volume, error) {
	var volumes []edgecloud.Volume

	volList, _, err := client.Volumes.List(ctx, nil)
	if err != nil {
		return nil, err
	}

	for _, volume := range volList {
		if volume.Name == name {
			volumes = append(volumes, volume)
		}
	}

	if len(volumes) == 0 {
		return nil, ErrVolumesNotFound
	}

	return volumes, nil
}

func WaitDiskAttachedToInstance(ctx context.Context, client *edgecloud.Client, volumeID, instanceID string) error {
	return client.Retryer.Run(ctx, func(ctx context.Context) error {
		volume, _, err := client.Volumes.Get(ctx, volumeID)
		if err != nil {
			return err
		}

		for _, attachment := range volume.Attachments {
			if instanceID == attachment.ServerID {
				return nil
			}
		}
		
		return ErrVolumesNotAttached
	})
}
