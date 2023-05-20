package token

import (
	"errors"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	ErrInvalidTokenFmt = errors.New("invalid JWT-token format")
	ErrExctractClaims  = errors.New("cannot extract JWT-token claims")
	exampleTokenSecret = []byte("secret secret ecret key =)")
)

func GetToken(userID string, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": jwt.MapClaims{
			"id":       userID,
			"username": username,
		},
	})
	tokenString, err := token.SignedString(exampleTokenSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetClaims(r *http.Request) (string, string, error) {
	fullToken := r.Header.Get("Authorization")
	tokenParts := strings.SplitN(fullToken, " ", 2)
	if len(tokenParts) < 2 {
		return "", "", ErrInvalidTokenFmt
	}
	claims, ok := extractClaims(tokenParts[1])
	if !ok {
		return "", "", ErrExctractClaims
	}

	u, ok := claims["user"]
	if !ok {
		return "", "", ErrExctractClaims
	}

	uid, ok := u.(map[string]interface{})["id"]
	if !ok {
		return "", "", ErrExctractClaims
	}

	uname, ok := u.(map[string]interface{})["username"]
	if !ok {
		return "", "", ErrExctractClaims
	}
	return uid.(string), uname.(string), nil
}

func extractClaims(tokenStr string) (jwt.MapClaims, bool) {
	hmacSecretString := exampleTokenSecret
	hmacSecret := hmacSecretString
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return hmacSecret, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		return nil, false
	}
}
