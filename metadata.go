package edgecloud

type MetadataDetailed struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	ReadOnly bool   `json:"read_only"`
}

type Metadata map[string]interface{}

// MetadataCreateRequest represent a metadata create struct.
type MetadataCreateRequest struct {
	Metadata
}
