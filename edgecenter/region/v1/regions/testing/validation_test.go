package testing

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/region/v1/regions"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/region/v1/types"
)

func TestUpdateOptsValidation(t *testing.T) {
	opts := regions.UpdateOpts{}
	err := edgecloud.ValidateStruct(opts)
	require.Error(t, err)

	opts = regions.UpdateOpts{
		State: types.RegionStateDeleted,
	}
	err = edgecloud.ValidateStruct(opts)
	require.NoError(t, err)

	opts = regions.UpdateOpts{
		State: types.RegionStateActive,
	}
	err = edgecloud.ValidateStruct(opts)
	require.NoError(t, err)

	opts = regions.UpdateOpts{
		DisplayName: "test",
	}
	err = edgecloud.ValidateStruct(opts)
	require.NoError(t, err)

	opts = regions.UpdateOpts{
		SpiceProxyURL: edgecloud.MustParseURL("http://test.com"),
	}
	err = edgecloud.ValidateStruct(opts)
	require.NoError(t, err)

	opts = regions.UpdateOpts{
		EndpointType: types.EndpointTypePublic,
	}
	err = edgecloud.ValidateStruct(opts)
	require.NoError(t, err)

	opts = regions.UpdateOpts{
		ExternalNetworkID: uuid.NewV4().String(),
	}
	err = edgecloud.ValidateStruct(opts)
	require.NoError(t, err)

}
