package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// AuthMiddleware is a middleware function to authenticate access tokens
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")
		parts := strings.Split(token, " ")

		// Check if the token is in the expected format
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
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
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Failed to read response:", err)
			return
		}

		var googleResponseData map[string]interface{}
		err = json.Unmarshal(body, &googleResponseData)
		if err != nil {
			log.Fatal("Failed to parse JSON response:", err)
			return
		}

		email, ok := googleResponseData["email"].(string)
		if ok {
			userCreateErr := createUserIfNotExists(email)
			if userCreateErr != nil {
				fmt.Fprintf(w, "Something went wrong")
				return
			}
		} else {
			fmt.Fprintf(w, "Something went wrong. Email not found")
			return
		}

		// Create a new context with the email value
		ctx := context.WithValue(r.Context(), "email", email)
		r = r.WithContext(ctx)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

func createUserIfNotExists(userEmail string) error {

	//stub out to user service
	return nil
}
