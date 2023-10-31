package edgecloud

const (
	metadataPath     = "metadata"
	metadataItemPath = "metadata_item"
)

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

type MetadataItemOptions struct {
	Key string `url:"key,omitempty" validate:"omitempty"`
}

// MetadataRoot represents a Metadata root.
type MetadataRoot struct {
	Count    int
	Metadata []MetadataDetailed `json:"results"`
}
