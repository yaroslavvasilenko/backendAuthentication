package main

import (
	"log"
	"net/http"
)

// We have some guid already = we send it in our request:
// localhost:8080/signin/?user-id=82e177bd14364bfea2425f63888e15f1

func main() {

	// "Signin" and "Welcome" are the handlers that we will implement
	http.HandleFunc("/", handleGuid)

	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func checkGuidInDataBase(guid []byte) bool {
	return true
}

func handleGuid(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["user-id"]
	if !ok || len(keys[0]) < 1 {
		// ToDo: return error of wrong/missing GUID
		return
	}
	key := keys[0]
	_, err := w.Write([]byte("Your key is: " + key))
	if err != nil {
		return
	}
}
