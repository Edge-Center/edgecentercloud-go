package util

import (
	"context"
	"errors"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

const InstanceShutoffStatus = "SHUTOFF"

var ErrInstanceNotShutOff = errors.New("the instance is not shut off")

func WaitForInstanceShutoff(ctx context.Context, client *edgecloud.Client, instanceID string, attempts *uint) error {
	_, _, err := client.Instances.InstanceStop(ctx, instanceID)
	if err != nil {
		return err
	}

	return WithRetry(
		func() error {
			instance, _, err := client.Instances.Get(ctx, instanceID)
			if err != nil {
				return err
			}

			if instance.Status == InstanceShutoffStatus {
				return nil
			}

			return ErrInstanceNotShutOff
		},
		attempts,
	)
}
