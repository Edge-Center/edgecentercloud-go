package util

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

func TestFindDefaultSG(t *testing.T) {
	ctx := context.Background()
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	sgs := []edgecloud.SecurityGroup{{ID: testResourceID, Name: "default"}, {ID: testResourceID2, Name: testName}}
	URL := path.Join("/v1/securitygroups/", strconv.Itoa(projectID), strconv.Itoa(regionID))

	mux.HandleFunc(URL, func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(sgs)
		if err != nil {
			t.Fatalf("failed to marshal JSON: %v", err)
		}
		_, _ = fmt.Fprintf(w, `{"results":%s}`, string(resp))
	})

	client := edgecloud.NewClient(nil)
	baseURL, _ := url.Parse(server.URL)
	client.BaseURL = baseURL
	client.Project = projectID
	client.Region = regionID

	defaultSG, err := FindDefaultSG(ctx, client)
	require.NoError(t, err)
	assert.Equal(t, testResourceID, defaultSG.ID)
}
