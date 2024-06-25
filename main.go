package main

import (
	"app/middleware"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func main() {
	opt := option.WithCredentialsFile("credentials.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error initializing app: %v\n", err)
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Error creating client: %v\n", err)
	}

	premiumHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// params, err := url.ParseQuery(r.URL.RawQuery)
		// if err != nil {
		// 	http.Error(w, "Bad request", http.StatusBadRequest)
		// 	return
		// }

		// idToken := r.Header.Get("auth-token")
		// idToken := params.Get("auth-token")
		// if idToken == "" {
		// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		// 	return
		// }

		token := r.Context().Value("verifiedToken")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// token, err := client.VerifyIDToken(context.Background(), idToken)
		// if err != nil {
		// 	// log.Fatalf("Error verifying token: %v\n", err)
		// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		// 	log.Printf("Error verifying token: %v\n", err)
		// 	return
		// }

		tmpl, err := template.ParseFiles("static/premium.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error parsing templates: %v\n", err)
			return
		}

		data := map[string]interface{}{
			"Token": token,
		}

		tmpl.Execute(w, data)

		// w.Header().Set("content-type", "application/json")
		// response := map[string]interface{}{
		// 	"message": "Verified token",
		// 	"token":   token,
		// }
		// json.NewEncoder(w).Encode(response)

		// fmt.Fprintf(w, "Verified token: %v\n", token)
	})

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)

	http.Handle("/premium", middleware.AuthMiddleWare(client)(premiumHandler))

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error creating server: %v\n", err)
	}
}
