package middlewares

import (
	"net/http"

	"go.uber.org/zap"
)

// AuthMiddleware is a middleware function to authenticate access tokens
func AuthMiddleware(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the access token from the request header
		//token := r.Header.Get("Authorization")

		// // Validate the access token (you can implement your own logic here)
		// if token != "your-access-token" {
		// 	logger.Warn("Unauthorized request", zap.String("path", r.URL.Path))
		// 	w.WriteHeader(http.StatusUnauthorized)
		// 	fmt.Fprint(w, "Unauthorized")
		// 	return
		// }

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
