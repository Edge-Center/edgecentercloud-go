package util

import (
	"context"
	"errors"

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
