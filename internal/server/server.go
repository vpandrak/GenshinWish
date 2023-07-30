package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"rest-todo/internal/Repository"
	"rest-todo/internal/auth"
	"strings"
)

func Serve(repo Repository.Repository, ctx context.Context) {
	auth := auth.Authicator{repo, ctx}
	mainMux := http.NewServeMux()
	accountHandler := http.HandlerFunc(account)
	mainMux.Handle("/", RequireAuth(accountHandler))
	mainMux.HandleFunc("/getunit", getUnit)
	mainMux.HandleFunc("/logout", auth.LogOut)
	mainMux.HandleFunc("/signin", auth.SignIn)
	mainMux.HandleFunc("/signup", auth.SignUp)

	err := http.ListenAndServe(":80", mainMux)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nonAuth := []string{"/signin", "/signup", "/getunit"}
		requestPath := r.URL.Path
		log.Println("register")
		for _, value := range nonAuth {
			if value == requestPath {
				return
			} else {
				if !isAuthorized(r) {
					// todo redirect to sign in
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unauthorized"))
					return
				}
			}
		}
		next.ServeHTTP(w, r)
		log.Println(isAuthorized(r))
	})
}

func isAuthorized(r *http.Request) bool {
	header := r.Header.Get("Authorization")
	if header == "" {
		return false
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		return false
	}
	fmt.Println(headerParts[1])
	userID, err := auth.ParseToken(headerParts[1])
	if err != nil {
		fmt.Println("erer")
		fmt.Println(err)
		return false
	}
	fmt.Println(userID)
	return true
}
