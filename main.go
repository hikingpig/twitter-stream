package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hikingpig/twitter-stream/stream"
)

func main() {
	twitterStream := stream.NewTwitterStream()
	msgs, errs, stop, err := twitterStream.Stream()
	if err != nil {
		log.Fatalf("unable to start stream: %s", err.Error())
	}
	stopSig := make(chan os.Signal, 1)
	signal.Notify(stopSig, syscall.SIGINT, syscall.SIGTERM)
	done := false
	for !done {
		select {
		case msg := <-msgs:
			fmt.Println("message:", msg)
		case err = <-errs:
			fmt.Println("error when streaming:", err)
			twitterStream.Close()
			done = true
		case <-stopSig:
			twitterStream.Close()
			done = true
		}
	}
	<-stop
}
