// Package acontext will contain all context keys shared within the applications.
package acontext

import (
	"context"
	"fmt"
)

// ContextKey is a custom type for representing context keys.
type ContextKey string

const (
	// ContextKeyUserID holds the context key for the user ID.
	ContextKeyUserID ContextKey = "userID"
	// ContextKeyRequestIDHeader holds the context key for the request ID header.
	ContextKeyRequestIDHeader ContextKey = "X-Request-ID"
	// ContextKeyCallerIDHeader holds the context key for the caller ID header.
	ContextKeyCallerIDHeader ContextKey = "X-Caller-ID"
	// ContextKeyForwardedForHeader holds the context key for the x-forwarded-for header.
	ContextKeyForwardedForHeader ContextKey = "X-Forwarded-For"
	// ContextKeyIPAddress holds the context key for an IP address.
	ContextKeyIPAddress ContextKey = "ipAddress"
	// ContextKeyTraceIDHeader holds the context key for the X-Cloud-Trace-Context header.
	ContextKeyTraceIDHeader ContextKey = "X-Cloud-Trace-Context"
)

// WithUserID creates a new context with the passed user ID.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ContextKeyUserID, userID)
}

// UserID attempts to retrieve the user ID from the context.
// It will return an error if no user ID is found.
func UserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(ContextKeyUserID).(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("userID is not in the context")
	}
	return userID, nil
}
