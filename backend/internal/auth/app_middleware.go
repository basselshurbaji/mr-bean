package auth

import (
	"net/http"
	"strings"

	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

func AppMiddleware(svc AppTokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				writeUnauthorized(w)
				return
			}

			raw := strings.TrimPrefix(header, "Bearer ")
			userID, err := svc.Validate(r.Context(), raw)
			if err != nil {
				writeUnauthorized(w)
				return
			}

			ctx := principal.WithUserID(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}