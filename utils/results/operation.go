// Package results provides generic result types for service layer operations.
package results

// OperationResult represents the outcome of a service operation.
//
// This type implements a clear distinction between:
//   - Domain outcomes (Success/Failure): Normal business logic results that should
//     be published as events. The handler acks the message regardless.
//   - Infrastructure errors (returned error): Technical failures that should trigger
//     a message retry (DB down, network errors, etc.)
//
// Usage patterns:
//   - Success only: Normal success case, publish success event
//   - Failure only: Domain failure (validation, not found), publish failure event
//   - Neither: Unexpected; should return an error instead
//   - Both: Not used; Success takes precedence
type OperationResult struct {
	// Success holds the success event payload when the operation completes normally.
	// The handler will publish this to the success topic.
	Success any

	// Failure holds the failure event payload when a domain failure occurs.
	// This is NOT a Go error - it's a normal domain outcome.
	// The handler will publish this to the failure topic.
	Failure any
}

// IsSuccess returns true if the result contains a success payload.
func (r OperationResult) IsSuccess() bool {
	return r.Success != nil
}

// IsFailure returns true if the result contains a failure payload.
func (r OperationResult) IsFailure() bool {
	return r.Failure != nil
}

// SuccessResult creates a result with a success payload.
func SuccessResult(payload any) OperationResult {
	return OperationResult{Success: payload}
}

// FailureResult creates a result with a failure payload.
func FailureResult(payload any) OperationResult {
	return OperationResult{Failure: payload}
}

// MapToHandlerResults maps an OperationResult to handler Result slice for the handler wrapper.
// successTopic and failureTopic are the topics to publish success/failure events to.
func (r OperationResult) MapToHandlerResults(successTopic, failureTopic string) []HandlerResult {
	if r.IsSuccess() {
		return []HandlerResult{{
			Topic:   successTopic,
			Payload: r.Success,
		}}
	}
	if r.IsFailure() {
		return []HandlerResult{{
			Topic:   failureTopic,
			Payload: r.Failure,
		}}
	}
	return nil
}

// HandlerResult represents a domain event outcome for the handler layer.
// This mirrors the Result type in handlerwrapper but is defined here
// to avoid circular imports.
type HandlerResult struct {
	Topic    string
	Payload  any
	Metadata map[string]string
}
