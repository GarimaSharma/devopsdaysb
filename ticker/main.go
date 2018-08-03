package main

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/mileusna/crontab"
)

func periodicFunc(tick time.Time) {
}
func main() {
	ctab := crontab.New()

	ctab.MustAddJob("* * * * *", myFunc2)

	for t := range time.NewTicker(2 * time.Second).C {
		periodicFunc(t)
	}
}
func myFunc2() {
	var buffer [100 * 1024 * 1024]string
	fmt.Printf("The size of the buffer is: %d bytes\n", unsafe.Sizeof(buffer))
	time.Sleep(300 * time.Second)
}