package util

import (
	"context"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

func FlavorIsExist(ctx context.Context, client *edgecloud.Client, name string) (bool, error) {
	flavors, _, err := client.Loadbalancers.FlavorList(ctx, nil)
	if err != nil {
		return false, err
	}

	for _, f := range flavors {
		if f.FlavorName == name {
			return true, nil
		}
	}

	return false, nil
}
