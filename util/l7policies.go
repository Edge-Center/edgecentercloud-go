package util

import (
	"context"
	"errors"
	"fmt"

	edgecloud "github.com/Edge-Center/edgecentercloud-go/v2"
)

var ErrL7PoliciesNotFound = errors.New("no l7Policies were found for the specified search criteria")

func L7PoliciesListByListenerID(ctx context.Context, client *edgecloud.Client, listenerID string) ([]edgecloud.L7Policy, error) {
	var L7Polices []edgecloud.L7Policy

	l7Policies, _, err := client.L7Policies.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, l7Policy := range l7Policies {
		if l7Policy.ListenerID == listenerID {
			L7Polices = append(L7Polices, l7Policy)
		}
	}

	if len(L7Polices) == 0 {
		return nil, ErrL7PoliciesNotFound
	}

	return L7Polices, nil
}

func GetLbL7PolicyFromName(ctx context.Context, client *edgecloud.Client, name string) (*edgecloud.L7Policy, error) {
	allPolicies, _, err := client.L7Policies.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error from getting list of l7 policies: %w", err)
	}
	var resultPolicies []edgecloud.L7Policy
	for _, policy := range allPolicies {
		if policy.Name == name {
			resultPolicies = append(resultPolicies, policy)
		}
	}

	var l7Policy *edgecloud.L7Policy
	switch len(resultPolicies) {
	case 1:
		l7Policy = &resultPolicies[0]
	case 0:
		return nil, fmt.Errorf("%w: resource \"l7policy\"; name \"%s\"", edgecloud.ErrResourceDoesntExist, name)
	default:
		return nil, fmt.Errorf("%w: resource \"l7policy\"; name \"%s\"", edgecloud.ErrMultipleResourcesWithTheSameName, name)
	}

	return l7Policy, nil
}
