package middlewares

import (
	"ckassa-callback/pkg/errors"
	"ckassa-callback/pkg/logger"
	"ckassa-callback/pkg/render"
	"context"
	"net/http"
)

var authUser = ctxKey("userId")

// WithUserInfo проверка авторизации пользователя.
func WithUserInfo(lg logger.Logger, next http.Handler, verifier UserVerifier) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		userId := req.Header.Get("X-Auth-User")
		if userId == "" {
			_ = render.JSON(wr, http.StatusUnauthorized, errors.Unauthorized("Basic auth header is not present"))
			return
		}

		verified := verifier.VerifySecret(req.Context(), userId)
		if !verified {
			wr.WriteHeader(http.StatusUnauthorized)
			_ = render.JSON(wr, http.StatusUnauthorized, errors.Unauthorized("Invalid username or secret"))
			return
		}

		req = req.WithContext(context.WithValue(req.Context(), authUser, userId))
		next.ServeHTTP(wr, req)
	})
}

// User extracts the username injected into the context by the auth middleware.
func User(req *http.Request) (string, bool) {
	val := req.Context().Value(authUser)
	if userName, ok := val.(string); ok {
		return userName, true
	}

	return "", false
}

type ctxKey string

// UserVerifier implementation is responsible for verifying the name-secret pair.
type UserVerifier interface {
	VerifySecret(ctx context.Context, userId string) bool
}

// UserVerifierFunc implements UserVerifier.
type UserVerifierFunc func(ctx context.Context, name, secret string) bool

// VerifySecret delegates call to the wrapped function.
func (uvf UserVerifierFunc) VerifySecret(ctx context.Context, name, secret string) bool {
	return uvf(ctx, name, secret)
}
