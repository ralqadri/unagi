package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/ralqadri/unagi/commands"
	"github.com/ralqadri/unagi/config"
)

var (
	BotId string
	dg *discordgo.Session
	cfg *config.Config
)

func Start() (*discordgo.Session, error) {
	var err error
	cfg, err = config.ReadConfig()

	if err != nil {
		fmt.Println("failed starting & reading config!: ", err)
		return nil, err
	}

	// apparently "Bot " is a required prefix for this token type 
	dg, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		fmt.Println("failed creating discord bot session!: ", err)
		return nil, err
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection to discord: ", err)
		return nil, err
	}

	fmt.Println("bot is now connected! ctrl+c to exit")
	return dg, nil
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore messages from bot itself/self responses
	if m.Author.ID == s.State.User.ID {
		return
	}

	prefix := cfg.BotPrefix
	content := m.Content

	// command: ping
	commands.HandlePingCommand(s, m, prefix, content)

	// command: echo
	commands.HandleEchoCommand(s, m, prefix, content)

	// command: dl
	commands.HandleDownloadCommand(s, m, prefix, content)
}