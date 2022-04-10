package main

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

func generateAllTokens(uid string, keyPrivate *rsa.PrivateKey) (signedAccessToken string, signedRefreshToken string, fTime time.Time, err error) {
	finishTime := time.Now().Local().Add(time.Hour * time.Duration(24))
	tokenAccess := createAccessToken(uid, keyPrivate)
	tokenRefresh := createRefreshToken(uid, keyPrivate)
	return tokenAccess, tokenRefresh, finishTime, err
}

func ParseToken(tokenString string, signingKey *rsa.PrivateKey) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		} // jwt.SigningMethodHMAC{} NO
		return signingKey.Public(), nil
	})
	if err != nil {
		fmt.Println(err)
	}

	// type-assert `Claims` into a variable of the appropriate type
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		fmt.Println(claims.Id)
	}
	return token.Claims.(*Claims), nil
}

func refreshOperation(Uuid string) {

}
