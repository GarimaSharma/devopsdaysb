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

	router.HandleFunc("/hello", getDetails).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getDetails(w http.ResponseWriter, r *http.Request) {
	var addr = strings.Split(r.RemoteAddr, ":")[0]
	println(addr)
	if timesRequestBy(addr) < 5 {
		json.NewEncoder(w).Encode("hellooooo 👋")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
	}
}

func timesRequestBy(remoteAddr string) int {
	if val, ok := times[remoteAddr]; ok {
		times[remoteAddr] = val + 1
	} else {
		times[remoteAddr] = 1
	}
	return times[remoteAddr]
}
