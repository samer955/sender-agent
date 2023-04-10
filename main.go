package main

import (
	"fmt"
	"github.com/samer955/sender-agent/service"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	sender := service.NewSender()
	sender.Start()

	//Run the program till its stopped (forced)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")

}
