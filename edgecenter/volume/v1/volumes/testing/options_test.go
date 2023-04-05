package testing

import (
	"testing"

	"github.com/fatih/structs"
	"github.com/stretchr/testify/require"

	"github.com/Edge-Center/edgecentercloud-go/edgecenter/volume/v1/volumes"
)

func TestCreateOpts(t *testing.T) {
	options := volumes.CreateOpts{
		Source:               volumes.NewVolume,
		Name:                 "test",
		TypeName:             volumes.Standard,
		ImageID:              "",
		SnapshotID:           "",
		InstanceIDToAttachTo: "",
	}

	_, err := options.ToVolumeCreateMap()
	require.NoError(t, err)
	require.Contains(t, structs.Names(options), "Size")

	options = volumes.CreateOpts{
		Source:               volumes.NewVolume,
		Name:                 "test",
		Size:                 10,
		TypeName:             volumes.Standard,
		ImageID:              "",
		SnapshotID:           "",
		InstanceIDToAttachTo: "",
	}

	_, err = options.ToVolumeCreateMap()
	require.NoError(t, err)

	options = volumes.CreateOpts{
		Source:               volumes.Snapshot,
		Name:                 "test",
		Size:                 10,
		TypeName:             volumes.Standard,
		ImageID:              "",
		SnapshotID:           "",
		InstanceIDToAttachTo: "",
	}

	_, err = options.ToVolumeCreateMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "SnapshotID")

	options = volumes.CreateOpts{
		Source:               volumes.Snapshot,
		Name:                 "test",
		Size:                 0,
		TypeName:             volumes.Standard,
		ImageID:              "",
		SnapshotID:           "726ecfcc-7fd0-4e30-a86e-7892524aa483",
		InstanceIDToAttachTo: "",
	}

	_, err = options.ToVolumeCreateMap()
	require.NoError(t, err)

	options = volumes.CreateOpts{
		Source:               volumes.Image,
		Name:                 "test",
		Size:                 0,
		TypeName:             volumes.Standard,
		ImageID:              "",
		SnapshotID:           "726ecfcc-7fd0-4e30-a86e-7892524aa483",
		InstanceIDToAttachTo: "",
	}

	_, err = options.ToVolumeCreateMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "ImageID")

	options = volumes.CreateOpts{
		Source:               volumes.Image,
		Name:                 "test",
		Size:                 10,
		TypeName:             volumes.Standard,
		ImageID:              "726ecfcc-7fd0-4e30-a86e-7892524aa483",
		SnapshotID:           "",
		InstanceIDToAttachTo: "",
	}

	_, err = options.ToVolumeCreateMap()
	require.NoError(t, err)

	options = volumes.CreateOpts{
		Source:               volumes.Image,
		Name:                 "test",
		Size:                 10,
		TypeName:             volumes.Standard,
		ImageID:              "726ecfcc-7fd0-4e30-a86e-7892524aa483",
		SnapshotID:           "726ecfcc-7fd0-4e30-a86e-7892524aa483",
		InstanceIDToAttachTo: "",
	}

	_, err = options.ToVolumeCreateMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "Snapshot")

	options = volumes.CreateOpts{
		Source: volumes.NewVolume,
		Name:   "test",
		Size:   10,
	}

	_, err = options.ToVolumeCreateMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "TypeName is a required field")

	options = volumes.CreateOpts{
		Source:   volumes.NewVolume,
		TypeName: "",
		Name:     "test",
		Size:     10,
	}

	_, err = options.ToVolumeCreateMap()
	require.Error(t, err)
	require.Contains(t, err.Error(), "TypeName is a required field")
}
