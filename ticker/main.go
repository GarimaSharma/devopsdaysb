package main

import (
	"fmt"
	"github.com/mileusna/crontab"
	"time"
)

func periodicFunc(tick time.Time) {
	fmt.Println("Tick at: ", tick)
}
func main() {
	ctab := crontab.New()

	ctab.MustAddJob("* * * * *", myFunc2)

	for t := range time.NewTicker(2 * time.Second).C {
		periodicFunc(t)
	}
}
func myFunc2(){
	println("hell yeah")
}