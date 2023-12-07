package edgecloud

type Name struct {
	Name string `json:"name"`
}

type ID struct {
	ID string `json:"id"`
}

type IDName struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// idNameRoot represents a IDName root.
type idNameRoot struct {
	Count   int
	IDNames []IDName `json:"results"`
}
