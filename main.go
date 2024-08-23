package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ralqadri/unagi/bot"
)

func main() {
	dg, err := bot.Start()
	if err != nil {
		fmt.Println("error starting bot!: ", err)
		return
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	dg.Close()

	fmt.Println("bot is shutting down...")
}