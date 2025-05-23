package middleware

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	kinolkHostEnv             = "KINOLK_FRONTEND_HOST"
	kinolkAllowedMethodsEnv   = "KINOLK_METHODS"
	kinolkAllowCredentialsEnv = "KINOLK_ALLOW_CRED"
	kinolkAllowedHeadersEnv   = "KINOLK_ALLOW_HEADERS"
)

// Middleware for enabling needed CORS
func MiddlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.Ctx(r.Context())

		w.Header().Set("Access-Control-Allow-Origin", viper.GetString(kinolkHostEnv))
		w.Header().Set("Access-Control-Allow-Methods", viper.GetString(kinolkAllowedMethodsEnv))
		w.Header().Set("Access-Control-Allow-Credentials", viper.GetString(kinolkAllowCredentialsEnv))
		w.Header().Set("Access-Control-Allow-Headers", viper.GetString(kinolkAllowedHeadersEnv))

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)

			logger.Info().Msg(fmt.Sprintf("options asked from %s", r.RequestURI))
			return
		}

		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

// Middleware for preventing any panic in server, so it won't instantly crash
func PreventPanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			logger := log.Ctx(r.Context())

			if err := recover(); err != nil {

				logger.Error().Msgf("Catched by middleware: panic happend: %v", err)

				http.Error(w, "Internal server error", 500)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
