package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"os"
)

var times map[string]int

func main() {
	f, err := os.OpenFile("/var/log/dragon-api.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	//defer to close when you're done with it, not because you think it's idiomatic!
	defer f.Close()

	//set output of logs to f
	log.SetOutput(f)

	router := mux.NewRouter()
	times = make(map[string]int)

	router.HandleFunc("/", sayHello).Methods("GET")
	router.HandleFunc("/healthcheck", getHealth).Methods("GET")
	router.HandleFunc("/getdragoncount", getDragonCount).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Healthy")
	log.Println("health check sent as healthy")
}

func getDragonCount(w http.ResponseWriter, r *http.Request) {
	var addr = strings.Split(r.RemoteAddr, ":")[0]
	println(addr)
	count := timesRequestBy(addr)
	if count < 10 {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Your dragon count is ")
		json.NewEncoder(w).Encode(count)
		log.Println("dragon count request responded.")
	} else if timesRequestBy(addr) == 20 {
		w.WriteHeader(http.StatusOK)
		log.Println("Throttling requests from and accepting from now on for", addr)
		times[addr] = 1
	} else {
		w.WriteHeader(http.StatusTooManyRequests)
		log.Println("Rejecting reuests because it is too many from", addr)
		w.Write([]byte("429 - too many requests"))
	}
}
func sayHello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("hello 👋, Welcome to World of Dragons!")
	log.Println("Hello request responded")
}

func timesRequestBy(remoteAddr string) int {
	if val, ok := times[remoteAddr]; ok {
		times[remoteAddr] = val + 1
	} else {
		times[remoteAddr] = 1
	}
	return times[remoteAddr]
}
