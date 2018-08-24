package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"strings"
)


func main() {
	router := mux.NewRouter()

	router.HandleFunc("/hello", getDetails).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getDetails(w http.ResponseWriter, r *http.Request) {
	var addr = strings.Split(r.RemoteAddr, ":")[0]
	println(addr)
	json.NewEncoder(w).Encode("hellooooo 👋")
}

