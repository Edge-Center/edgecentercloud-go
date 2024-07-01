package util

import (
	"context"
	"errors"
	"fmt"

	edgecloud "github.com/Edge-Center/edgecentercloud-go/v2"
)

const InstanceShutoffStatus = "SHUTOFF"

var (
	ErrInstanceNotShutOff        = errors.New("the instance is not shut off")
	ErrInstanceInterfaceNotFound = errors.New("instance interface not found")
	ErrInstancePortNotFound      = errors.New("instance port not found")
)

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

func InstanceNetworkInterfaceByID(ctx context.Context, client *edgecloud.Client, instanceID string, portID string) (*edgecloud.InstancePortInterface, error) {
	instanceIfaceList, _, err := client.Instances.InterfaceList(ctx, instanceID)
	if err != nil {
		return nil, err
	}
	for _, iface := range instanceIfaceList {
		if iface.PortID == portID {
			return &iface, nil
		}
	}

	return nil, fmt.Errorf("%w :there is no interface port with id %s in instance with id %s", ErrInstanceInterfaceNotFound, portID, instanceID)
}

func InstanceNetworkPortByID(ctx context.Context, client *edgecloud.Client, instanceID string, portID string) (*edgecloud.InstancePort, error) {
	instancePortList, _, err := client.Instances.PortsList(ctx, instanceID)
	if err != nil {
		return nil, err
	}
	for _, port := range instancePortList {
		if port.ID == portID {
			return &port, nil
		}
	}

	return nil, fmt.Errorf("%w :there is no port with id %s in instance with id %s", ErrInstancePortNotFound, portID, instanceID)
}
