package main

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type mongoPattern struct {
	Uuid  string `bson:"uuid" json:"uuid"`
	Key   []byte `bson:"key" json:"key"`
	Token string `bson:"token" json:"token"`
}

type twoToken struct {
	TokenRe string `json:"token_re"`
	TokenAc string `json:"token_ac"`
}

func (app *Applecation) firstRoute(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["user-id"]
	if !ok || len(keys[0]) < 1 {
		// ToDo: return error of wrong/missing GUID
		return
	}
	uuid := keys[0]
	privateKey, _ := generatePrivatAndPublicKey()
	tokenAccess, tokenRefresh, _, err := generateAllTokens(uuid, privateKey)
	sendingToken(twoToken{tokenRefresh, tokenAccess}, w)
	re := &mongoPattern{
		Uuid:  uuid,
		Token: tokenRefresh,
		Key:   x509.MarshalPKCS1PrivateKey(privateKey),
	}

	_, err = app.UserAuth.InsertOne(context.TODO(), re)
	if err != nil {
		log.Fatal(err)
	}
}
func (app *Applecation) sekondRoute(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["user-id"]
	if !ok || len(keys[0]) < 1 {
		// ToDo: return error of wrong/missing GUID
		return
	}
	uuid := keys[0]
	tokenReAc := parseBodyToken(r)
	kl := app.findMongUuid(uuid)

	// check Refresh token
	if kl.Token != tokenReAc.TokenRe {
		return
	}
	privKey, _ := x509.ParsePKCS1PrivateKey(kl.Key)
	q, _ := ParseToken(kl.Token, privKey)

	// check Time
	if q.ExpiresAt < time.Now().Local().Add(time.Hour*time.Duration(1)).Unix() {
		return
	}

	key, _ := generatePrivatAndPublicKey()
	tokenAccess, tokenRefresh, _, _ := generateAllTokens(uuid, key)
	app.updateTokeninMongo(tokenRefresh, uuid)

	sendingToken(twoToken{tokenRefresh, tokenAccess}, w)
}
func (app *Applecation) findMongUuid(uuid string) mongoPattern {
	//result := database.FindMongo(app.userAuth)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var result mongoPattern
	err := app.UserAuth.FindOne(ctx, bson.M{"uid": uuid}).Decode(&result)
	defer cancel()
	if err == mongo.ErrNoDocuments {
		// Do something when no record was found
		fmt.Println("record does not exist")
	} else if err != nil {
		log.Fatal(err)
	}
	return result
}

func (app *Applecation) updateTokeninMongo(signedRefreshToken string, uuid string) {
	filter := bson.D{{"uid", uuid}}
	update := bson.D{
		{"$inc", bson.D{
			{"token_re", signedRefreshToken},
		}},
	}

	updateResult, _ := app.UserAuth.UpdateOne(context.TODO(), filter, update)
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
}

// send token to body user
func sendingToken(token twoToken, w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(token)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)

}

func parseBodyToken(r *http.Request) *twoToken {
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	token := &twoToken{}
	err := json.Unmarshal(body, &token)
	if err != nil {
		log.Fatal(err)
	}
	return token
}

func checkGuidInDataBase(guid []byte) bool {
	return true
}

//
//_, err = w.Write([]byte("Your key is: " + uuid))
//if err != nil {
//return
//}
