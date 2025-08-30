package commands

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"sync"
	"time"
)

func dotenvToDuration(key string) time.Duration {
	str := os.Getenv(key)
	if str == "" {
		log.Fatal("Missing environment variable: " + key)
	}
	if str[len(str)-1] < 'a' || str[len(str)-1] > 'z' {
		str += "m"
	}
	duration, err := time.ParseDuration(str)
	if err != nil {
		fmt.Println("error parsing duration: ", err)
	}
	return duration
}

func startNewCycle(wg *sync.WaitGroup, lock chan struct{}) {
	defer wg.Done()
	duration := dotenvToDuration("AUTO_INTERVAL_NEW")
	color.Magenta("Starting new cycle every %s", duration.String())
	for range time.Tick(duration) {
		lock <- struct{}{} // acquire
		fmt.Println(time.Now().Format(time.RubyDate))
		New()
		<-lock // release
	}
}

func updateRunningCycles(wg *sync.WaitGroup, lock chan struct{}) {
	defer wg.Done()
	duration := dotenvToDuration("AUTO_INTERVAL_UPDATE")
	color.Magenta("Updating running cycles every %s", duration.String())
	for range time.Tick(duration) {
		lock <- struct{}{} // acquire
		fmt.Println(time.Now().Format(time.RubyDate))
		Update()
		<-lock // release
	}
}

func Auto() {
	color.Yellow("Starting Auto Mode - CTRL + C to exit")

	var wg sync.WaitGroup
	lock := make(chan struct{}, 1) // channel used as mutex

	wg.Add(2)
	go startNewCycle(&wg, lock)
	go updateRunningCycles(&wg, lock)

	wg.Wait()
}
