package authentication

import (
	"context"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
)

// ContextReader defines operations for reading authentication data from context
//
//go:generate mockgen -destination=./mocks/context_mock.go -package=authentication_mocks github.com/alexgrauroca/practice-food-delivery-platform/pkg/clients/authentication ContextReader
type ContextReader interface {
	GetSubject(ctx context.Context) (string, bool)
	RequireSubjectMatch(ctx context.Context, expectedSubject string) error
}

type contextReader struct {
	logger log.Logger
}

// NewContextReader creates a new ContextReader instance
func NewContextReader(logger log.Logger) ContextReader {
	return &contextReader{logger: logger}
}

// GetSubject retrieves the token subject from the given context.
// It returns the subject value and a boolean indicating whether the subject was found
// and successfully type asserted to string.
func (r *contextReader) GetSubject(ctx context.Context) (string, bool) {
	v := ctx.Value(subjectCtxKey)
	if v == nil {
		return "", false
	}
	subject, ok := v.(string)

	return subject, ok
}

// RequireSubjectMatch checks if the subject of the given context matches the expected value.
// If the subject is not found or does not match, it returns an error.
func (r *contextReader) RequireSubjectMatch(ctx context.Context, expectedSubject string) error {
	subject, ok := r.GetSubject(ctx)
	if !ok {
		r.logger.Warn("authentication context not found")
		return ErrInvalidToken
	}
	if subject != expectedSubject {
		r.logger.Warn(
			"subject mismatch with the token",
			log.Field{Key: "subject", Value: subject},
			log.Field{Key: "expectedSubject", Value: expectedSubject},
		)
		return ErrSubjectMismatch
	}
	return nil
}
