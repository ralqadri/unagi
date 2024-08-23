package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/ralqadri/unagi/config"
)

var (
	BotId string
	goBot *discordgo.Session
	cfg *config.Config
)

func Start() {
	var err error
	cfg, err = config.ReadConfig()

	if err != nil {
		fmt.Println("failed starting & reading config!: ", err)
		return 
	}

	goBot, err = discordgo.New("bot " + cfg.Token)
	if err != nil {
		fmt.Println("failed creating discord bot session!: ", err)
		return
	}

	user, err := goBot.User("@me")
	if err != nil {
		fmt.Println("error obtaining current user: ", err)
	}

	BotId = user.ID

	goBot.AddHandler(messageHandler)

	err = goBot.Open()
	if err != nil {
		fmt.Println("error opening connection to discord: ", err)
		return
	}

	fmt.Println("bot is now connected!")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore messages from bot itself/self responses
	if m.Author.ID == BotId {
		return
	}

	prefix := cfg.BotPrefix
	if m.Content == prefix + "ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong!")
	}
}