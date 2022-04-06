package main

import (
	"github.com/golang-jwt/jwt"
	"log"
	"time"
)

const (
	signingKey = "efwuef93u93bf9u23f9b39ufbu47bg75gb"
	userID     = "eif93ufn934ufub"
)

type JwtCustomClaims struct {
	Username string `json:"username"`
	Guid     string `json:"guid"`
	jwt.StandardClaims
}

const gu = "eif93jf93jf"

func getTokens() (map[string]string, error) {
	// Create token
	var tokenAccess = JwtCustomClaims{
		Guid: gu,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES512, tokenAccess)

	tokenAc, _ := token.SignedString([]byte(signingKey)) //внутри токен

	log.Println(tokenAc, "Hi")

	var tokenRefresh = JwtCustomClaims{
		Guid: gu,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(100 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token = jwt.NewWithClaims(jwt.SigningMethodES512, tokenRefresh)
	tokenRe, _ := token.SignedString([]byte(signingKey)) //внутри токен

	return map[string]string{
		"access_token":  tokenAc,
		"refresh_token": tokenRe,
	}, nil
	//token, _ := jwt.Parse(ta, func(token *jwt.Token)) (interface{}, error)
}
