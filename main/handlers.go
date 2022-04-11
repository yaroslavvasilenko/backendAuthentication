package main

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type mongoPattern struct {
	Token string `bson:"token" json:"token"`
	Uuid  string `bson:"uuid" json:"uuid"`
}

type tokenPairId struct {
	TokenRe string `json:"token_re"`
	TokenAc string `json:"token_ac"`
	Uuid    string `json:"uuid"`
}

//  creates 2 linked tokens(SHA512), transmits them in base64 format
//  for a user with a uuid,
//  sends a refresh token and user uuid to the database
func (app *Application) firstRoute(w http.ResponseWriter, r *http.Request) {
	// Get the user-id GUID from request
	keys, ok := r.URL.Query()["user-id"]
	if !ok || len(keys[0]) < 1 {
		log.Println("not found uuid")
		return
	}

	uuid := keys[0]
	privateKey := app.Secret
	tokenAccess, tokenRefresh := generateAllTokens(uuid, privateKey)
	sendingToken(&tokenPairId{tokenRefresh, tokenAccess, uuid}, w)
	mongoP := &mongoPattern{
		Uuid:  uuid,
		Token: tokenRefresh,
	}

	_, err := app.UserAuth.InsertOne(context.TODO(), mongoP)
	if err != nil {
		log.Fatal(err)
	}
}

//  searches for a user by id, checks refresh token
//  creates 2 linked tokens(SHA512), transmits them in base64 format
//  for a user with a uuid,
//  update Refresh token in database
//  if the token is not valid, the function stops executing
func (app Application) secondRoute(w http.ResponseWriter, r *http.Request) {
	tokenReAcId := parseBodyToken(r)

	refreshTokenBd := app.findMongUuid(tokenReAcId.Uuid)

	// check Refresh token
	if refreshTokenBd.Token != tokenReAcId.TokenRe {
		return
	}
	valid, tokenRe := decoding(tokenReAcId)
	if !valid {
		return
	}
	privateKey := app.Secret
	// check valid token
	userRefreshToken, valid := ParseToken(tokenRe, privateKey)

	if !valid {
		return
	}
	_, validAc := ParseToken(tokenReAcId.TokenAc, privateKey)
	if !validAc {
		return
	}

	// check Time
	if userRefreshToken.ExpiresAt < time.Now().Local().Unix() {
		return
	}

	tokenAccess, tokenRefresh := generateAllTokens(tokenReAcId.Uuid, privateKey)
	app.updateTokeninMongo(tokenRefresh, tokenReAcId.Uuid)

	sendingToken(&tokenPairId{tokenRefresh, tokenAccess, tokenReAcId.Uuid}, w)
}

// search user data in BD
func (app Application) findMongUuid(uuid string) mongoPattern {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var result mongoPattern
	err := app.UserAuth.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&result)
	defer cancel()
	if err == mongo.ErrNoDocuments {
		// Do something when no record was found
		log.Println("record does not exist")
	} else if err != nil {
		log.Fatal(err)
	}
	return result
}

// Update user information
func (app Application) updateTokeninMongo(signedRefreshToken string, uuid string) {
	filter := bson.D{{"uuid", uuid}}
	update := bson.D{
		{"$set", bson.D{
			{"token", signedRefreshToken},
		}},
	}

	updateResult, _ := app.UserAuth.UpdateOne(context.TODO(), filter, update)
	log.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
}

// send token to body user
func sendingToken(token *tokenPairId, w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(token)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)

}

func parseBodyToken(r *http.Request) *tokenPairId {
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	token := &tokenPairId{}
	err := json.Unmarshal(body, &token)
	if err != nil {
		log.Fatal(err)
	}
	return token
}
