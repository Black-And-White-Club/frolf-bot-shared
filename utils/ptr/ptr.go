// Package ptr provides pointer conversion utilities for Go.
//
// These helpers are useful for:
//   - Partial update patterns where nil means "don't update"
//   - Building SQL update structs with optional fields
//   - Converting between value and pointer types cleanly
package ptr

import "time"

// Of returns a pointer to the provided value.
// Useful for creating pointers to literals.
func Of[T any](v T) *T {
	return &v
}

// IfNonEmpty returns a pointer to the string if non-empty, otherwise nil.
// Useful for partial update patterns where empty string means "don't update".
func IfNonEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// IfTrue returns a pointer to the bool if true, otherwise nil.
// Useful for partial update patterns where false means "don't update".
func IfTrue(b bool) *bool {
	if !b {
		return nil
	}
	return &b
}

// IfNotNil returns a pointer to a value if the input pointer is not nil.
// The transform function converts the dereferenced value to the target type.
func IfNotNil[T, R any](v *T, transform func(T) R) *R {
	if v == nil {
		return nil
	}
	result := transform(*v)
	return &result
}

// TimeToUnixNano converts a time pointer to unix nanoseconds pointer.
// Useful for storing time as int64 in databases.
func TimeToUnixNano(t *time.Time) *int64 {
	if t == nil {
		return nil
	}
	nano := t.UnixNano()
	return &nano
}

// UnixNanoToTime converts unix nanoseconds to a time pointer.
// Returns nil if nano is 0 (representing no time).
func UnixNanoToTime(nano int64) *time.Time {
	if nano == 0 {
		return nil
	}
	t := time.Unix(0, nano)
	return &t
}

// ValueOr returns the dereferenced value or a default if nil.
func ValueOr[T any](v *T, defaultVal T) T {
	if v == nil {
		return defaultVal
	}
	return *v
}

// Value returns the dereferenced value or zero value if nil.
func Value[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}
