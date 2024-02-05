package util

import (
	edgecloud "github.com/Edge-Center/edgecentercloud-go/v2"
)

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
