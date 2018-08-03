package main

import (
	"encoding/json"
	"log"
	// "fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var times map[string]int

func main() {
	router := mux.NewRouter()
	times = make(map[string]int)

	router.HandleFunc("/hello", getDetails).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func getDetails(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("hellooooo 👋")
	json.NewEncoder(w).Encode(r.RemoteAddr)
	json.NewEncoder(w).Encode(timesRequestBy(r.RemoteAddr))
}

func timesRequestBy(remoteAddr string) int {
	if val, ok := times[remoteAddr]; ok {
		times[remoteAddr] = val + 1
	} else {
		times[remoteAddr] = 1
	}
	return times[remoteAddr]
}
