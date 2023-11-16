package util

import (
	"context"
	"errors"
	"net"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
)

var (
	ErrLoadbalancersNotFound     = errors.New("no Loadbalancers were found for the specified search criteria")
	ErrLoadbalancerPoolsNotFound = errors.New("no Loadbalancer pools were found for the specified search criteria")
	ErrMultipleResults           = errors.New("multiple results where only one expected")
	ErrErrorState                = errors.New("loadbalancer in Error state")
	ErrNotActiveStatus           = errors.New("waiting for Active status")
)

func LoadbalancerGetByName(ctx context.Context, client *edgecloud.Client, name string) (*edgecloud.Loadbalancer, error) {
	var matchedLBs []edgecloud.Loadbalancer

	lbs, _, err := client.Loadbalancers.List(ctx, nil)
	if err != nil {
		return nil, err
	}

	for _, lb := range lbs {
		if lb.Name == name &&
			lb.ProvisioningStatus != edgecloud.ProvisioningStatusDeleted &&
			lb.ProvisioningStatus != edgecloud.ProvisioningStatusPendingDelete {
			matchedLBs = append(matchedLBs, lb)
		}
	}

	switch len(matchedLBs) {
	case 1:
		return &matchedLBs[0], nil
	case 0:
		return nil, ErrLoadbalancersNotFound
	default:
		return nil, ErrMultipleResults
	}
}

func LBPoolGetByName(ctx context.Context, client *edgecloud.Client, name, loadBalancerID string) (*edgecloud.Pool, error) {
	var matchedLBPools []edgecloud.Pool

	poolListOptions := edgecloud.PoolListOptions{
		LoadBalancerID: loadBalancerID,
		Details:        true,
	}
	lbPools, _, err := client.Loadbalancers.PoolList(ctx, &poolListOptions)
	if err != nil {
		return nil, err
	}

	for _, pool := range lbPools {
		if pool.Name == name &&
			pool.ProvisioningStatus != edgecloud.ProvisioningStatusDeleted &&
			pool.ProvisioningStatus != edgecloud.ProvisioningStatusPendingDelete {
			matchedLBPools = append(matchedLBPools, pool)
		}
	}

	switch len(matchedLBPools) {
	case 1:
		return &matchedLBPools[0], nil
	case 0:
		return nil, ErrLoadbalancerPoolsNotFound
	default:
		return nil, ErrMultipleResults
	}
}

func WaitLoadBalancerProvisioningStatusActive(ctx context.Context, client *edgecloud.Client, loadBalancerID string, attempts *uint) error {
	return WithRetry(
		func() error {
			loadBalancer, _, err := client.Loadbalancers.Get(ctx, loadBalancerID)
			if err != nil {
				return err
			}

			switch loadBalancer.ProvisioningStatus { // nolint: exhaustive
			case edgecloud.ProvisioningStatusActive:
				return err
			case edgecloud.ProvisioningStatusError:
				return ErrErrorState
			default:
				return ErrNotActiveStatus
			}
		},
		attempts,
	)
}

func FindPoolMemberByAddressPortAndSubnetID(pool edgecloud.Pool, addr net.IP, protocolPort int, subnetID string) (found bool) {
	for _, member := range pool.Members {
		if member.Address == nil {
			continue
		}

		if net.IP.Equal(member.Address, addr) && member.ProtocolPort == protocolPort && member.SubnetID == subnetID {
			found = true
			break
		}
	}

	return found
}

func DeleteUnusedPools(ctx context.Context, client *edgecloud.Client, oldPools []edgecloud.Pool, newPools []string, attempts *uint) error {
	for _, oldPool := range oldPools {
		lbID := oldPool.LoadBalancers[0].ID

		var exist bool
		for _, newPool := range newPools {
			if oldPool.ID == newPool {
				exist = true
				break
			}
		}

		if !exist {
			task, _, err := client.Loadbalancers.PoolDelete(ctx, oldPool.ID)
			if err != nil {
				return err
			}
			if err = WaitForTaskComplete(ctx, client, task.Tasks[0]); err != nil {
				return err
			}

			if err = WaitLoadBalancerProvisioningStatusActive(ctx, client, lbID, attempts); err != nil {
				return err
			}
		}
	}

	return nil
}
