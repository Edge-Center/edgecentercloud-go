package util

import (
	"context"
	"errors"
	"net"

	edgecloud "github.com/Edge-Center/edgecentercloud-go/v2"
)

var (
	ErrLoadbalancersNotFound           = errors.New("no Loadbalancers were found for the specified search criteria")
	ErrLoadbalancerPoolsNotFound       = errors.New("no Loadbalancer pools were found for the specified search criteria")
	ErrLoadbalancerPoolsMemberNotFound = errors.New("no Loadbalancer pool member were found for the specified search criteria")
	ErrLoadbalancerListenerNotFound    = errors.New("no Loadbalancer listener were found for the specified search criteria")
	ErrMultipleResults                 = errors.New("multiple results where only one expected")
	ErrErrorState                      = errors.New("loadbalancer in Error state")
	ErrNotActiveStatus                 = errors.New("waiting for Active status")
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

func LBListenerGetByName(ctx context.Context, client *edgecloud.Client, name, loadBalancerID string) (*edgecloud.Listener, error) {
	var matchedLBListeners []edgecloud.Listener

	listenerListOptions := edgecloud.ListenerListOptions{
		LoadbalancerID: loadBalancerID,
	}
	lbListeners, _, err := client.Loadbalancers.ListenerList(ctx, &listenerListOptions)
	if err != nil {
		return nil, err
	}

	for _, lis := range lbListeners {
		if lis.Name == name &&
			lis.ProvisioningStatus != edgecloud.ProvisioningStatusDeleted &&
			lis.ProvisioningStatus != edgecloud.ProvisioningStatusPendingDelete {
			matchedLBListeners = append(matchedLBListeners, lis)
		}
	}

	switch len(matchedLBListeners) {
	case 1:
		return &matchedLBListeners[0], nil
	case 0:
		return nil, ErrLoadbalancerListenerNotFound
	default:
		return nil, ErrMultipleResults
	}
}

func LBPoolGetByName(ctx context.Context, client *edgecloud.Client, name, loadBalancerID string) (*edgecloud.Pool, error) {
	var matchedLBPools []edgecloud.Pool

	poolListOptions := edgecloud.PoolListOptions{
		LoadbalancerID: loadBalancerID,
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

func PoolMemberGetByID(ctx context.Context, client *edgecloud.Client, poolID, memberID string) (*edgecloud.PoolMember, error) {
	lbPool, _, err := client.Loadbalancers.PoolGet(ctx, poolID)
	if err != nil {
		return nil, err
	}

	for _, member := range lbPool.Members {
		if member.ID == memberID {
			return &member, nil
		}
	}

	return nil, ErrLoadbalancerPoolsMemberNotFound
}

func LBSharedPoolList(ctx context.Context, client *edgecloud.Client, loadBalancerID string) ([]edgecloud.Pool, error) {
	var sharedPools []edgecloud.Pool

	poolListOptions := edgecloud.PoolListOptions{
		LoadbalancerID: loadBalancerID,
	}
	lbPools, _, err := client.Loadbalancers.PoolList(ctx, &poolListOptions)
	if err != nil {
		return nil, err
	}

	for _, pool := range lbPools {
		if len(pool.Listeners) > 0 {
			sharedPools = append(sharedPools, pool)
		}
	}

	return sharedPools, nil
}

func WaitLoadbalancerProvisioningStatusActive(ctx context.Context, client *edgecloud.Client, loadBalancerID string, attempts *uint) error {
	return WithRetry(
		func() error {
			loadBalancer, _, err := client.Loadbalancers.Get(ctx, loadBalancerID)
			if err != nil {
				return err
			}

			switch loadBalancer.ProvisioningStatus { //nolint:exhaustive
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

func DeletePoolByNameIfExist(ctx context.Context, client *edgecloud.Client, name, loadBalancerID string) error {
	pool, err := LBPoolGetByName(ctx, client, name, loadBalancerID)
	if err != nil {
		if errors.Is(err, ErrLoadbalancerPoolsNotFound) {
			return nil
		}
		return err
	}

	task, _, err := client.Loadbalancers.PoolDelete(ctx, pool.ID)
	if err != nil {
		return err
	}

	return WaitForTaskComplete(ctx, client, task.Tasks[0])
}

func DeleteUnusedPools(ctx context.Context, client *edgecloud.Client, oldPools []edgecloud.Pool, newPoolsIDs []string, attempts *uint) error {
	for _, oldPool := range oldPools {
		lbID := oldPool.Loadbalancers[0].ID

		var exist bool
		for _, newPool := range newPoolsIDs {
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

			if err = WaitLoadbalancerProvisioningStatusActive(ctx, client, lbID, attempts); err != nil {
				return err
			}
		}
	}

	return nil
}
