package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"rest-todo/internal/model"
	"time"
)

type Claims struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	jwt.StandardClaims
}

var Key = []byte("dgaggl;k3i40-k")

func ParseToken(tokenString string) (id int, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return Key, nil
	})

	if err != nil {
		fmt.Println(1)
		fmt.Println(err)
		return 0, err
	}

	claims, ok := token.Claims.(Claims)
	if !ok {
		fmt.Println(2)
		fmt.Println(err)
		return 0, err

	}

	return claims.Id, nil
}

func generateToken(u *model.User, expTime time.Time) (signedToken string, err error) {
	claims := &Claims{
		Id:   u.Id,
		Name: u.Name,
		StandardClaims: jwt.StandardClaims{
			Subject:   u.Name,
			ExpiresAt: expTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(Key) //error

	if err != nil {
		return "", err
	}
	return tokenString, nil
}
