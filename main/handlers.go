package main

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type mongoPattern struct {
	Uid     string
	TokenRe string
	KeyPriv rsa.PrivateKey
}

func (app *applecation) handleGuid(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["user-id"]
	if !ok || len(keys[0]) < 1 {
		// ToDo: return error of wrong/missing GUID
		return
	}
	uuid := keys[0]

	//result := database.FindMongo(app.userAuth)
	//fmt.Println(result)
	_, re, keySec, _ := GenerateAllTokens("dwiqwi9nfunf")
	str := base64.StdEncoding.EncodeToString([]byte(re))

	userProfil := &mongoPattern{
		Uid:     uuid,
		TokenRe: str,
		KeyPriv: *keySec,
	}
	//CA, mes := ValidateToken(Ac)

	insertResult, err := app.userAuth.InsertOne(context.TODO(), userProfil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	_, err = w.Write([]byte("Your key is: " + uuid))
	if err != nil {
		return
	}
}

func id() string {
	return uuid.New().String()
}

func checkGuidInDataBase(guid []byte) bool {
	return true
}
