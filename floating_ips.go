package edgecloud

type FloatingIP struct {
	ID                string                 `json:"id"`
	CreatedAt         string                 `json:"created_at"`
	UpdatedAt         string                 `json:"updated_at"`
	Status            string                 `json:"status"`
	FixedIPAddress    string                 `json:"fixed_ip_address"`
	FloatingIPAddress string                 `json:"floating_ip_address"`
	DNSDomain         string                 `json:"dns_domain"`
	DNSName           string                 `json:"dns_name"`
	RouterID          string                 `json:"router_id"`
	SubnetID          string                 `json:"subnet_id"`
	CreatorTaskID     string                 `json:"creator_task_id"`
	TaskID            string                 `json:"task_id"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

type FloatingIPSource string

const (
	NewFloatingIP      FloatingIPSource = "new"
	ExistingFloatingIP FloatingIPSource = "existing"
)

type InterfaceFloatingIP struct {
	Source             FloatingIPSource `json:"source" validate:"required,enum"`
	ExistingFloatingID string           `json:"existing_floating_id" validate:"rfe=Source:existing,sfe=Source:new,omitempty,uuid"`
}
