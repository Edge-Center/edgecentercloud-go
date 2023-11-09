package util

import (
	"context"
	"errors"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

var ErrFloatingIPsNotFound = errors.New("no FloatingIPs were found for the specified search criteria")

func FloatingIPsListByPortID(ctx context.Context, client *edgecloud.Client, portID string) ([]edgecloud.FloatingIP, error) {
	var FloatingIPs []edgecloud.FloatingIP

	fips, _, err := client.Floatingips.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, fip := range fips {
		if fip.PortID == portID {
			FloatingIPs = append(FloatingIPs, fip)
		}
	}

	if len(FloatingIPs) == 0 {
		return nil, ErrFloatingIPsNotFound
	}

	return FloatingIPs, nil
}
