package main

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/gommon/log"
	"time"
)

type SignedDetails struct {
	Uid string
	jwt.StandardClaims
}

func GenerateAllTokens(uid string) (signedToken string, signedRefreshToken string, k *rsa.PrivateKey, err error) {
	claims := &SignedDetails{
		Uid: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS512, claims).SignedString(key)
	if err != nil {
		log.Panic(err)
		return
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims).SignedString(key)
	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, key, err
}

// ValidateToken validates the jwt token
//func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
//	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, errors.New("unexpected signing method")
//		}
//		return []byte(SECRET_KEY), nil
//	})
//
//	if err != nil {
//		msg = err.Error()
//		return
//	}
//
//	claims, ok := token.Claims.(*SignedDetails)
//	if !ok {
//		msg = fmt.Sprintf("the token is invalid")
//		msg = err.Error()
//		return
//	}
//	if claims.ExpiresAt < time.Now().Local().Unix() {
//		msg = fmt.Sprintf("token is expired")
//		msg = err.Error()
//		return
//	}
//
//	return claims, msg
//}
