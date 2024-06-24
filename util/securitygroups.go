package util

import (
	"context"
	"errors"
	"slices"

	edgecloud "github.com/Edge-Center/edgecentercloud-go/v2"
)

var ErrDefaultSGNotFound = errors.New("default security group is not found")

type SecurityGroupRuleProtocol edgecloud.SecurityGroupRuleProtocol

func (s SecurityGroupRuleProtocol) List() []edgecloud.SecurityGroupRuleProtocol {
	return []edgecloud.SecurityGroupRuleProtocol{
		edgecloud.SGRuleProtocolANY,
		edgecloud.SGRuleProtocolAH,
		edgecloud.SGRuleProtocolDCCP,
		edgecloud.SGRuleProtocolEGP,
		edgecloud.SGRuleProtocolESP,
		edgecloud.SGRuleProtocolGRE,
		edgecloud.SGRuleProtocolICMP,
		edgecloud.SGRuleProtocolIGMP,
		edgecloud.SGRuleProtocolIPIP,
		edgecloud.SGRuleProtocolOSPF,
		edgecloud.SGRuleProtocolPGM,
		edgecloud.SGRuleProtocolRSVP,
		edgecloud.SGRuleProtocolSCTP,
		edgecloud.SGRuleProtocolTCP,
		edgecloud.SGRuleProtocolUDP,
		edgecloud.SGRuleProtocolUDPLITE,
		edgecloud.SGRuleProtocolVRRP,
	}
}

func (s SecurityGroupRuleProtocol) StringList() []string {
	protocols := s.List()
	strings := make([]string, 0, len(protocols))
	for _, protocol := range protocols {
		strings = append(strings, string(protocol))
	}

	return strings
}

func FindDefaultSG(ctx context.Context, client *edgecloud.Client) (*edgecloud.SecurityGroup, error) {
	allSGs, _, err := client.SecurityGroups.List(ctx, nil)
	if err != nil {
		return nil, err
	}
	defaultSGIndex := slices.IndexFunc(allSGs, func(group edgecloud.SecurityGroup) bool {
		return group.Name == "default"
	})
	if defaultSGIndex >= 0 {
		return &allSGs[defaultSGIndex], nil
	}

	return nil, ErrDefaultSGNotFound
}
