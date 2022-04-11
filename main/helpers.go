package main

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt"
	"log"
	"math/rand"
	"time"
)

type Claims struct {
	jwt.StandardClaims
	Uid string `json:"uid"`
}

const charts string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%&*()=-,"

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

func generateRandomString(size int) string {
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)

	output := make([]byte, size)
	randomness := make([]byte, size)
	_, err := rnd.Read(randomness)
	if err != nil {
		return ""
	}
	l := uint8(len(charts))
	for pos := range output {
		random := uint8(randomness[pos])
		randomPos := random % l
		output[pos] = charts[randomPos]
	}
	return string(output)
}
