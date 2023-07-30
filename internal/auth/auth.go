package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"rest-todo/internal/Repository"
	"rest-todo/internal/model"
	"time"
)

type Authicator struct {
	Repository.Repository
	context.Context
}

type Authorizer interface {
	SignInToken(ctx context.Context, user *model.User) (string, error)
}

func (auth Authicator) SignIn(w http.ResponseWriter, r *http.Request) {
	var u model.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	account, err := auth.GetByName(auth.Context, u.Name)
	if errors.Is(err, Repository.ErrNotExist) {
		res := fmt.Sprint("account doesnt exist")
		if _, err := io.WriteString(w, res); err != nil {
			log.Println("cannot write in sign in")
		}
		return
	}

	u.Password = hasher(u.Password)

	if u.Password == account.Password {
		res := fmt.Sprint("welcome ", account.Name)
		expirationTime := time.Now().Add(5 * time.Minute)
		token, err := generateToken(&u, expirationTime)

		if err != nil {
			log.Println(err)
		}

		var bearer = "Bearer " + token
		if r.Method != http.MethodPost {
			w.Header().Set("Authorization", bearer)
		}

		if _, err := io.WriteString(w, res); err != nil {
			log.Println("cannot write in sign in")
		}

	} else {
		res := fmt.Sprint("invalid login or password")
		if _, err := io.WriteString(w, res); err != nil {
			log.Println("cannot write in sign in")
		}
	}

}

func (auth Authicator) SignUp(w http.ResponseWriter, r *http.Request) {
	var u model.User
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u.Password = hasher(u.Password)
	_, err = auth.Create(auth.Context, u)

	if errors.Is(err, Repository.ErrDuplicate) {
		fmt.Print("record: %+v already exists\n", u)
	} else if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Println("some error")
	}
	fmt.Println(u.Name)
}

func (auth Authicator) LogOut(w http.ResponseWriter, r *http.Request) {
	w.Header().Del("Authorization")
}
