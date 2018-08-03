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
	if timesRequestBy(r.RemoteAddr) < 5 {
		json.NewEncoder(w).Encode("hellooooo 👋")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
	}

	//json.NewEncoder(w).Encode(timesRequestBy(r.RemoteAddr))

	//	var SkipRouter = errors.New("skip this router")
}

func timesRequestBy(remoteAddr string) int {
	if val, ok := times[remoteAddr]; ok {
		times[remoteAddr] = val + 1
	} else {
		times[remoteAddr] = 1
	}
	return times[remoteAddr]
}
