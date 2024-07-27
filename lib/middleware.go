package lib

import "net/http"

type AuthMiddleware struct {
	cfg *Config
}

func NewAuthMiddleware(cfg *Config) *AuthMiddleware {
	return &AuthMiddleware{
		cfg: cfg,
	}
}

func (a *AuthMiddleware) ValidateCurrentUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := ValidateCurrentUser(a.cfg.ApiConfig, r)
		if err != nil {
			WriteResponse(w, err, nil)
			return
		}

		r = SetUserContext(r, user)
		next.ServeHTTP(w, r)
	})
}
