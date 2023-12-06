package util

import (
	"context"
	"errors"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

var ErrServerGroupNotFound = errors.New("no server group was found for the specified search criteria")

func ServerGroupGetByInstance(ctx context.Context, client *edgecloud.Client, instanceID string) (*edgecloud.ServerGroup, error) {
	sgs, _, err := client.ServerGroups.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, sg := range sgs {
		for _, instance := range sg.Instances {
			if instance.InstanceID == instanceID {
				return &sg, nil
			}
		}
	}

	return nil, ErrServerGroupNotFound
}
