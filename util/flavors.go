package util

import (
	"context"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

func LoadbalancerFlavorIsExist(ctx context.Context, client *edgecloud.Client, flavorName string) (bool, error) {
	flavors, _, err := client.Loadbalancers.FlavorList(ctx, nil)
	if err != nil {
		return false, err
	}

	for _, f := range flavors {
		if f.FlavorName == flavorName {
			return true, nil
		}
	}

	return false, nil
}

func FlavorIsExist(ctx context.Context, client *edgecloud.Client, flavorName string) (bool, error) {
	flavors, _, err := client.Flavors.List(ctx, nil)
	if err != nil {
		return false, err
	}

	for _, f := range flavors {
		if f.FlavorName == flavorName {
			return true, nil
		}
	}

	return false, nil
}

func FlavorIsAvailable(ctx context.Context, client *edgecloud.Client, flavorName string, checkFlavorVolumeRequest *edgecloud.InstanceCheckFlavorVolumeRequest) (bool, error) {
	flavors, _, err := client.Instances.AvailableFlavors(ctx, checkFlavorVolumeRequest, nil)
	if err != nil {
		return false, err
	}

	for _, f := range flavors {
		if f.FlavorName == flavorName {
			return true, nil
		}
	}

	return false, nil
}
