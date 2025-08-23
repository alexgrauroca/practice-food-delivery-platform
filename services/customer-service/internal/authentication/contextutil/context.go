package contextutil

import "context"

const SubjectCtxKey = "token-subject"

func GetSubject(ctx context.Context) (string, bool) {
	v := ctx.Value(SubjectCtxKey)
	if v == nil {
		return "", false
	}
	subject, ok := v.(string)
	return subject, ok
}
