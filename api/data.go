package api

type Beast struct {
	BeastName   string            `json:"BeastName"`
	Type        string            `json:"Type"`
	CR          string            `json:"CR"`
	Attributes  map[string]string `json:"Attributes"`
	Description string            `json:"Description"`
}
