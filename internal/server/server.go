package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"rest-todo/internal/Repository"
	"rest-todo/internal/auth"
)

func Serve(repo Repository.Repository, ctx context.Context) {
	auth := auth.Authicator{repo, ctx}
	mainMux := http.NewServeMux()
	authMiddleWare := RequireAuth(mainMux)
	mainMux.Handle("/", authMiddleWare)
	mainMux.HandleFunc("/getunit", getUnit)
	mainMux.HandleFunc("/signin", auth.SignIn)
	mainMux.HandleFunc("/signup", auth.SignUp)

	err := http.ListenAndServe(":80", mainMux)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func RequireAuth(next *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nonAuth := []string{"/signin", "/signup"}
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
	})
}

func isAuthorized(r *http.Request) bool {
	return true
}
