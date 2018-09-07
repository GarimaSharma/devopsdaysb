package main

import (
	"log"
	"time"
	"unsafe"

	"github.com/mileusna/crontab"
	"os"
)

func periodicFunc(tick time.Time) {
}
func main() {
	f, err := os.OpenFile("/var/log/scheduler.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	//defer to close when you're done with it, not because you think it's idiomatic!
	defer f.Close()

	//set output of logs to f
	log.SetOutput(f)

	ctab := crontab.New()

	 ctab.MustAddJob("0 11 * * *", myFunc2)

	for t := range time.NewTicker(10 * time.Minute).C {
		periodicFunc(t)
	}

}
//This is chaos monkey trying to increase the memory of system at 3 pm everyday
func myFunc2() {
	var buffer [100 * 1024 * 1024]string
	log.Println("The size of the buffer is: %d bytes\n", unsafe.Sizeof(buffer))
	time.Sleep(5 * time.Minute)
}