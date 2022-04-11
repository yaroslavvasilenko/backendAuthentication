package main

import (
	"crypto/rsa"
	"errors"
	"github.com/golang-jwt/jwt"
	"log"
	"strings"
)

//  creates linked refresh and access tokens
func generateAllTokens(uid string, keyPrivate *rsa.PrivateKey) (signedAccessToken string, signedRefreshToken string) {
	tokenAccess := createAccessToken(uid, keyPrivate)
	tokenRefresh := createRefreshToken(uid, keyPrivate)
	tokenRefresh = encoding(tokenAccess, tokenRefresh)
	return tokenAccess, tokenRefresh
}

// ParseToken checks whether the token is valid and returns its structure
// We don't return any errors because any produced error indicates that the token is invalid
func ParseToken(tokenString string, signingKey *rsa.PrivateKey) (*Claims, bool) {
	f := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return signingKey.Public(), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, f)

	if err != nil {
		log.Println(err)
		return nil, false
	}

	// type-assert `Claims` into a variable of the appropriate type
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		log.Println(claims.Id)
		return token.Claims.(*Claims), true
	}
	return token.Claims.(*Claims), false
}

// token binding and additional protection
func encoding(tokenAccess string, tokenRefresh string) string {
	q := strings.Split(tokenAccess, ".")
	w := strings.Split(tokenRefresh, ".")
	randStr := generateRandomString(12)
	rightJwt := q[2][len(q[2])-6:]
	w[2] = randStr + w[2] + rightJwt
	return strings.Join(w, ".")
}

//token decoupling and validation check
func decoding(pair *tokenPairId) (bool, string) {
	tokenRe, tokenAc := pair.TokenRe, pair.TokenAc
	flag := false
	tokAc := strings.Split(tokenAc, ".")
	tokenRefresh := strings.Split(tokenRe, ".")

	rightJwt := tokAc[2][len(tokAc[2])-6:]
	connJwt := tokenRefresh[2][len(tokenRefresh[2])-6:]
	tokenRefresh[2] = tokenRefresh[2][12 : len(tokenRefresh[2])-6]
	if rightJwt == connJwt {
		flag = true
	}
	return flag, strings.Join(tokenRefresh, ".")
}
