package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var times map[string]int

func main() {
	router := mux.NewRouter()
	times = make(map[string]int)

	router.HandleFunc("/", sayHello).Methods("GET")
	router.HandleFunc("/healthcheck", getHealth).Methods("GET")
	router.HandleFunc("/getDragonCount", getDetails).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Healthy")
}

func getDetails(w http.ResponseWriter, r *http.Request) {
	var addr = strings.Split(r.RemoteAddr, ":")[0]
	println(addr)
	count := timesRequestBy(addr)
	if count < 5 {
		json.NewEncoder(w).Encode("You dragon count is ")
		json.NewEncoder(w).Encode(count)
	} else if timesRequestBy(addr) == 50 {
		println("Throtling requests from %s. Accepting from now on", addr)
		times[addr] = 1
	} else {
		println("too many requests from %s. Rejecting it from now on", addr)
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("429 - too many requests"))
	}
}
func sayHello(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("hello 👋, Welcome to World of Dragons!")
}

func timesRequestBy(remoteAddr string) int {
	if val, ok := times[remoteAddr]; ok {
		times[remoteAddr] = val + 1
	} else {
		times[remoteAddr] = 1
	}
	return times[remoteAddr]
}
