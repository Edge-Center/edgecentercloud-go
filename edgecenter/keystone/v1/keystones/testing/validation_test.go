package testing

import (
	"testing"

	"github.com/stretchr/testify/require"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/keystone/v1/keystones"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/keystone/v1/types"
)

func TestUpdateOptsValidation(t *testing.T) {
	opts := keystones.UpdateOpts{}
	err := edgecloud.TranslateValidationError(opts.Validate())
	require.Error(t, err)
	opts = keystones.UpdateOpts{
		State: types.KeystoneStateDeleted,
	}
	err = edgecloud.TranslateValidationError(opts.Validate())
	require.NoError(t, err)
}
