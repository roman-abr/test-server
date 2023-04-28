package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var secret = []byte("secret")

type AccessToken struct {
	AccessToken string `json:"access_token"`
}

func GenerateToken(id primitive.ObjectID) AccessToken {
	payload := jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	signed, err := token.SignedString(secret)
	if err != nil {
		// return
	}
	return AccessToken{AccessToken: signed}
}

func CheckToken(token string) (any, bool) {
	claims := jwt.MapClaims{}
	result, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, false
	}
	id, err := claims.GetSubject()
	if err != nil {
		return nil, false
	}
	return id, result.Valid
}
