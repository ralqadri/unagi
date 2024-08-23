package main

import (
	"github.com/ralqadri/unagi/bot"
)

func main() {
	bot.Start()

	<-make(chan struct{})
}