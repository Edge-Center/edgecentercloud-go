package util

import (
	"context"
	"errors"
	"fmt"

	edgecloud "github.com/Edge-Center/edgecentercloud-go/v2"
)

var ErrSecGroupNotFound = errors.New("security group is not found")

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

func SecurityGroupListByIDs(ctx context.Context, client *edgecloud.Client, sgIDs []string) (sgs []edgecloud.SecurityGroup, err error) {
	allSgs, _, err := client.SecurityGroups.List(ctx, nil)
	if err != nil {
		return nil, err
	}

	allSGsMap := make(map[string]edgecloud.SecurityGroup)
	for _, sg := range allSgs {
		allSGsMap[sg.ID] = sg
	}

	for _, sgID := range sgIDs {
		if sg, ok := allSGsMap[sgID]; ok {
			sgs = append(sgs, sg)
		} else {
			return nil, fmt.Errorf("%w: id %s", ErrSecGroupNotFound, sgID)
		}
	}

	return sgs, nil
}
