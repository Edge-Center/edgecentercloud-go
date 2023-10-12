package edgecloud

// Flavor represents an EdgecenterCloud Flavor.
type Flavor struct {
	FlavorID            string            `json:"flavor_id"`
	FlavorName          string            `json:"flavor_name"`
	VCPUS               int               `json:"vcpus"`
	RAM                 int               `json:"ram"`
	HardwareDescription map[string]string `json:"hardware_description"`
	Disabled            bool              `json:"disabled"`
	ResourceClass       string            `json:"resource_class"`
}
