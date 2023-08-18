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
	for {
		select {
		case msg := <-msgs:
			fmt.Println("message:", msg)
		case err = <-errs:
			fmt.Println("error when streaming:", err)
			twitterStream.Close()
			<-stop // waiting for cleaning up

			// restart stream when error happens like connection closed by server
			twitterStream = stream.NewTwitterStream()
			msgs, errs, stop, err = twitterStream.Stream()
			if err != nil {
				log.Fatalf("unable to start stream: %s", err.Error())
			}
		case <-stopSig: // user force program to stop
			twitterStream.Close()
			<-stop
			return
		}
	}
}
