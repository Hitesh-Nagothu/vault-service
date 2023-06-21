package middlewares

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// AuthMiddleware is a middleware function to authenticate access tokens
func AuthMiddleware(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")
		parts := strings.Split(token, " ")

		// Check if the token is in the expected format
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			logger.Warn("Unauthorized request", zap.String("path", r.URL.Path))
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized")
			return
		}

		accessToken := parts[1]

		url := fmt.Sprintf("https://www.googleapis.com/oauth2/v3/tokeninfo?access_token=%s", accessToken)

		tokenValidityReq, tokenValidityErr := http.NewRequest("GET", url, nil)
		if tokenValidityErr != nil {
			log.Fatal("failed to create token validity request")
			return
		}

		client := http.Client{}
		resp, err := client.Do(tokenValidityReq)
		if err != nil {
			log.Fatal("Request failed:", err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Failed to read response:", err)
		}

		fmt.Println("Response:", string(body))

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
