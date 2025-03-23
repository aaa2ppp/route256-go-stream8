package middleware

import (
	"context"
	"net/http"
)

type authTokenContextKey struct{}

// ContextWithLogger добавляет логгер в контекст.
func ContextWithAuthToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, authTokenContextKey{}, token)
}

// GetLoggerFromContext возвращает логгер из контекста или nil, если логгер не найден.
func GetAuthTokenFromContext(ctx context.Context) string {
	if v := ctx.Value(authTokenContextKey{}); v != nil {
		return v.(string)
	}
	return ""
}

func Auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Authtoken")
		if token != "" {
			ctx := ContextWithAuthToken(r.Context(), token)
			r = r.WithContext(ctx)
		}
		h.ServeHTTP(w, r)
	})
}
