package audit

import "context"

type loginMetaKey struct{}

type LoginMeta struct {
	RequestID string
	IP        string
	UserAgent string
}

func WithLoginMeta(ctx context.Context, meta LoginMeta) context.Context {
	return context.WithValue(ctx, loginMetaKey{}, meta)
}

func LoginMetaFromContext(ctx context.Context) LoginMeta {
	meta, _ := ctx.Value(loginMetaKey{}).(LoginMeta)
	return meta
}
