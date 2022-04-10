package main

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/golang-jwt/jwt"
	"log"
	"time"
)

type Claims struct {
	jwt.StandardClaims
	Uid string `json:"uid"`
}

func createAccessToken(uid string, key *rsa.PrivateKey) string {
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			Id:        uid,
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(1)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS512, claims).SignedString(key) // This is Token
	if err != nil {
		log.Panicln(err)
	}
	return token
}

func createRefreshToken(uid string, key *rsa.PrivateKey) string {
	refreshClaims := &Claims{
		StandardClaims: jwt.StandardClaims{
			Id:        uid,
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims).SignedString(key) // This is Token
	if err != nil {
		log.Panicln(err)
	}
	return refreshToken
}
func generatePrivatAndPublicKey() (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}
	publicKey := privateKey.PublicKey

	return privateKey, &publicKey
}
