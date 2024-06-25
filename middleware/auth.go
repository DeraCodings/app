package middleware

import (
	"context"
	"log"
	"net/http"
	"net/url"

	"firebase.google.com/go/auth"
)

func AuthMiddleWare(client *auth.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// parse Query parameters
			params, err := url.ParseQuery(r.URL.RawQuery)
			if err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				// log.Fatalf("")
				return
			}

			idToken := params.Get("auth-token")
			if idToken == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// verify ID token
			verifiedToken, err := client.VerifyIDToken(context.Background(), idToken)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				log.Printf("Error verifying token: %v\n", err)
			}

			// Add token to the request context
			ctx := context.WithValue(r.Context(), "verifiedToken", verifiedToken)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
