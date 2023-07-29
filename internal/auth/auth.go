package auth

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"rest-todo/internal/Repository"
	"rest-todo/internal/model"
)

type Authicator struct {
	Repository.Repository
	context.Context
}

const haltSalt = "dfasfklsadkfl;"

type Authorizer interface {
	SignInToken(ctx context.Context, user *model.User) (string, error)
}

func SignInToken(ctx context.Context, user *model.User) (string, error) {
	pwd := sha1.New()
	pwd.Write([]byte(user.Password))
	pwd.Write([]byte(haltSalt))
	user.Password = fmt.Sprintf("%x", pwd.Sum(nil))
	return "df", nil
}

func (auth Authicator) SignIn(w http.ResponseWriter, r *http.Request) {
	var u model.User
	err := json.NewEncoder(w).Encode(&u)
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
		res := fmt.Sprint("welcome %s", account.Name)
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

func hasher(password string) string {
	hasher := sha1.New()
	hasher.Write([]byte(password))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}
