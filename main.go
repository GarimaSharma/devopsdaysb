package main

import (
	  "encoding/json"
    "log"
		// "fmt"
    "net/http"
    "github.com/gorilla/mux"
)

// our main function
func main() {
    router := mux.NewRouter()

		router.HandleFunc("/hello", getDetails).Methods("GET")
		log.Fatal(http.ListenAndServe(":8000", router))
}

func getDetails(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("hellooooo 👋")
	json.NewEncoder(w).Encode(r.RemoteAddr)
}

