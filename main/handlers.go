package main

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"
)

type mongoPattern struct {
	Uid     string
	TokenRe string
	KeyPriv rsa.PrivateKey
}

func (app *Applecation) firstRoute(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["user-id"]
	if !ok || len(keys[0]) < 1 {
		// ToDo: return error of wrong/missing GUID
		return
	}
	uuid := keys[0]

	privateKey, _ := generatePrivatAndPublicKey()
	_, tokenRefresh, finishTime, _ := generateAllTokens(uuid, privateKey)
	//result := database.FindMongo(app.userAuth)

	re := &mongoPattern{
		Uid:     uuid,
		TokenRe: tokenRefresh,
		KeyPriv: *privateKey,
	}

	insertResult, err := app.UserAuth.InsertOne(context.TODO(), re)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenRefresh,
		Expires: finishTime,
	})

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenRefresh,
		Expires: finishTime,
	})

}
func (app *Applecation) sekondRoute(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the JWT string from the cookie
	tknStr := c.Value
	ParseToken(tknStr)

	// Initialize a new instance of `Claims`
	//refreshOperation()
	w.Write([]byte("Отображение заметки..."))
}

//CA, mes := ValidateToken(Ac)

//http.SetCookie(w, &http.Cookie{
//	Name:    "token",
//	Value:   tokenRefresh,
//	Expires: expirationTime,
//})
//
//func id() string {
//	return uuid.New().String()
//}

func checkGuidInDataBase(guid []byte) bool {
	return true
}

//
//_, err = w.Write([]byte("Your key is: " + uuid))
//if err != nil {
//return
//}
