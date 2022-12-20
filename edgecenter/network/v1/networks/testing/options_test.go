package testing

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Edge-Center/edgecentercloud-go/edgecenter/network/v1/networks"
)

func TestCreateOpts(t *testing.T) {
	options := networks.CreateOpts{
		Name: Network1.Name,
		Mtu:  1450,
		Type: "vxlan",
	}
	_, err := options.ToNetworkCreateMap()
	require.NoError(t, err)

	options = networks.CreateOpts{
		Name:         Network1.Name,
		Mtu:          1501,
		CreateRouter: true,
	}
	_, err = options.ToNetworkCreateMap()
	require.Error(t, err)

}
