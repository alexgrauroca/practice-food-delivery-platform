package authentication

import "context"

// ContextReader defines operations for reading authentication data from context
//
//go:generate mockgen -destination=./mocks/context_mock.go -package=authentication_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/authentication ContextReader
type ContextReader interface {
	GetSubject(ctx context.Context) (string, bool)
}

type contextReader struct{}

// NewContextReader creates a new ContextReader instance
func NewContextReader() ContextReader {
	return &contextReader{}
}

// GetSubject retrieves the token subject from the given context.
// It returns the subject value and a boolean indicating whether the subject was found
// and successfully type asserted to string.
func (ctxReader *contextReader) GetSubject(ctx context.Context) (string, bool) {
	v := ctx.Value(subjectCtxKey)
	if v == nil {
		return "", false
	}
	subject, ok := v.(string)

	return subject, ok
}
