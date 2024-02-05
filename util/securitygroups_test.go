package util

import (
	"testing"

	"github.com/stretchr/testify/assert"

	edgecloud "github.com/Edge-Center/edgecentercloud-go/v2"
)

func TestSecurityGroupsList(t *testing.T) {
	result := []edgecloud.SecurityGroupRuleProtocol{
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

	SGPs := SecurityGroupRuleProtocol("").List()

	assert.Equal(t, result, SGPs)
}

func TestSecurityGroupsStringList(t *testing.T) {
	result := []string{
		"any",
		"ah",
		"dccp",
		"egp",
		"esp",
		"gre",
		"icmp",
		"igmp",
		"ipip",
		"ospf",
		"pgm",
		"rsvp",
		"sctp",
		"tcp",
		"udp",
		"udplite",
		"vrrp",
	}

	SGPString := SecurityGroupRuleProtocol("").StringList()

	assert.Equal(t, result, SGPString)
}
