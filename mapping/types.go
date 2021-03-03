package mapping

type IncomingRoute struct {
	Host  string                  `json:"host"`
	Start string                  `json:"start"`
	Steps map[string]IncomingStep `json:"steps"`
}

type IncomingStep struct {
	Service    string            `json:"service"`
	Config     interface{}       `json:"config,omitempty"`
	References map[string]string `json:"references,omitempty"`
}
