package main

import (
	"fmt"
	"time"
)

func main() {

	var stopSpinner bool
	var c chan struct{} = make(chan struct{}) // event marker

	stopSpinner = false
	go spinner(1*time.Second, stopSpinner, c)

	// do something...
	stopSpinner = true

	<-c // wait spinner stop

	//do something else

}

func spinner(delay time.Duration, stopSpinner bool, c chan struct{}) {
	for !stopSpinner {
		for _, r := range `-\|/` {

			fmt.Printf("\r%c", r)
			time.Sleep(delay)

		}
	}

	c <- struct{}{}
}
