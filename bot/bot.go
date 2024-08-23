package bot

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/ralqadri/unagi/config"
)

var (
	BotId string
	dg *discordgo.Session
	cfg *config.Config
)

func Start() {
	var err error
	cfg, err = config.ReadConfig()

	if err != nil {
		fmt.Println("failed starting & reading config!: ", err)
		return 
	}

	// apparently "Bot " is a required prefix for this token type 
	dg, err = discordgo.New("Bot " + cfg.Token)
	if err != nil {
		fmt.Println("failed creating discord bot session!: ", err)
		return
	}

	// user, err := dg.User("@me")
	// if err != nil {
	// 	fmt.Println("error obtaining current user: ", err)
	// }

	// BotId = user.ID

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection to discord: ", err)
		return
	}

	fmt.Println("bot is now connected! ctrl+c to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore messages from bot itself/self responses
	if m.Author.ID == s.State.User.ID {
		return
	}

	prefix := cfg.BotPrefix
	if m.Content == prefix + "ping" {
		s.ChannelMessageSend(m.ChannelID, "pong!")
	}
}