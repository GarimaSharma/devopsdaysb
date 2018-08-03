package main

import (
	"fmt"
	"time"
)

func periodicFunc(tick time.Time) {
	fmt.Println("Tick at: ", tick)
}
func main() {
	for t := range time.NewTicker(2 * time.Second).C {
		periodicFunc(t)
	}
}
