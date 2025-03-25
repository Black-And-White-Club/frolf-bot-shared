package events

// ResultPayload is a common payload structure for indicating success/failure.
type ResultPayload struct {
	CommonMetadata
	Status      string `json:"result"`
	ErrorDetail string `json:"reason,omitempty"`
}
