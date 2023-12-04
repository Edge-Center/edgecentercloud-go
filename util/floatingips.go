package util

import (
	"context"
	"errors"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

var (
	ErrFloatingIPsNotFound = errors.New("no FloatingIPs were found for the specified search criteria")
	ErrFloatingIPNotFound  = errors.New("no FloatingIP was found for the specified search criteria")
)

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

func FloatingIPDetailedByIPAddress(ctx context.Context, client *edgecloud.Client, floatingIPAddress string) (*edgecloud.FloatingIP, error) {
	fips, _, err := client.Floatingips.List(ctx)
	if err != nil {
		return nil, err
	}

	var floatingIP edgecloud.FloatingIP
	var found bool
	for _, fip := range fips {
		if fip.FloatingIPAddress == floatingIPAddress {
			floatingIP = fip
			found = true
			break
		}
	}

	if !found {
		return nil, ErrFloatingIPNotFound
	}

	return &floatingIP, nil
}

func FloatingIPDetailedByID(ctx context.Context, client *edgecloud.Client, id string) (*edgecloud.FloatingIP, error) {
	fips, _, err := client.Floatingips.List(ctx)
	if err != nil {
		return nil, err
	}

	var floatingIP edgecloud.FloatingIP
	var found bool
	for _, fip := range fips {
		if fip.ID == id {
			floatingIP = fip
			found = true
			break
		}
	}

	if !found {
		return nil, ErrFloatingIPNotFound
	}

	return &floatingIP, nil
}
