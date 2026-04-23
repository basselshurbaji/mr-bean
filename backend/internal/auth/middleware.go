package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

func Middleware(ts TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				writeUnauthorized(w)
				return
			}

			raw := strings.TrimPrefix(header, "Bearer ")
			claims, err := ts.ValidateAccessToken(raw)
			if err != nil {
				writeUnauthorized(w)
				return
			}

			ctx := principal.WithUserID(r.Context(), claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"}) //nolint:errcheck
}
