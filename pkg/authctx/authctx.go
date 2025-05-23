package authctx

import "context"

type ctxKeyJWT string

const jwtTokenKey ctxKeyJWT = "jwtToken"

func WithJWT(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, jwtTokenKey, token)
}

func GetJWT(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(jwtTokenKey).(string)
	return token, ok
}
