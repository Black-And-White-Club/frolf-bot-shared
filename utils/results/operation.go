// Package results provides generic result types for service layer operations.
package results

// OperationResult represents the outcome of a service operation using Generics.
// S (Success): The domain type returned on success (e.g., []LeaderboardEntry).
// F (Failure): The domain type returned on expected failures (e.g., ValidationErrors).
//
// This type implements a clear distinction between:
//   - Domain outcomes (Success/Failure): Normal business logic results that should
//     be published as events. The handler acks the message regardless.
//   - Infrastructure errors (returned error): Technical failures that should trigger
//     a message retry (DB down, network errors, etc.)
type OperationResult[S any, F any] struct {
	// Success holds the success payload.
	Success *S

	// Failure holds the domain failure payload.
	Failure *F
}

// IsSuccess returns true if the result contains a success payload.
func (r OperationResult[S, F]) IsSuccess() bool {
	return r.Success != nil
}

// IsFailure returns true if the result contains a failure payload.
func (r OperationResult[S, F]) IsFailure() bool {
	return r.Failure != nil
}

// SuccessResult creates a type-safe result with a success payload.
func SuccessResult[S any, F any](payload S) OperationResult[S, F] {
	return OperationResult[S, F]{Success: &payload}
}

// FailureResult creates a type-safe result with a failure payload.
func FailureResult[S any, F any](payload F) OperationResult[S, F] {
	return OperationResult[S, F]{Failure: &payload}
}

// Map transforms the internal types.
// Note: We use OperationResult[any, any] as the return to "erase"
// the specific domain types after they've been converted to Event Payloads.
func (r OperationResult[S, F]) Map(
	onSuccess func(S) any,
	onFailure func(F) any,
) OperationResult[any, any] {
	if r.IsSuccess() {
		val := onSuccess(*r.Success)
		return OperationResult[any, any]{Success: &val}
	}
	if r.IsFailure() {
		val := onFailure(*r.Failure)
		return OperationResult[any, any]{Failure: &val}
	}
	return OperationResult[any, any]{}
}

// MapToHandlerResults can now be a standalone function or
// a method on the "erased" type.
func (r OperationResult[S, F]) ToHandlerResults(successTopic, failureTopic string) []HandlerResult {
	var results []HandlerResult
	if r.IsSuccess() {
		results = append(results, HandlerResult{
			Topic:   successTopic,
			Payload: *r.Success,
		})
	} else if r.IsFailure() {
		results = append(results, HandlerResult{
			Topic:   failureTopic,
			Payload: *r.Failure,
		})
	}
	return results
}

// HandlerResult represents a domain event outcome for the handler layer.
type HandlerResult struct {
	Topic    string
	Payload  any
	Metadata map[string]string
}
