package edgecloud

type VolumeType string

const (
	Standard  VolumeType = "standard"
	SsdHiIops VolumeType = "ssd_hiiops"
	Cold      VolumeType = "cold"
	Ultra     VolumeType = "ultra"
)
