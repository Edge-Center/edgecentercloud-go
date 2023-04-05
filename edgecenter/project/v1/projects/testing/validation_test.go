package testing

import (
	"testing"

	"github.com/stretchr/testify/require"

	edgecloud "github.com/Edge-Center/edgecentercloud-go"
	"github.com/Edge-Center/edgecentercloud-go/edgecenter/project/v1/projects"
)

func TestUpdateOptsValidation(t *testing.T) {
	opts := projects.UpdateOpts{}
	err := edgecloud.TranslateValidationError(opts.Validate())
	require.Error(t, err)
	opts = projects.UpdateOpts{
		Name: "test",
	}
	err = edgecloud.TranslateValidationError(opts.Validate())
	require.NoError(t, err)
}
