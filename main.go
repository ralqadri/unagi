package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/ralqadri/unagi/cmd"
	"github.com/ralqadri/unagi/config"
)

// user-defined variables/parameters; configs/tokens whatever
// TODO: add option to add guildID in either config file or env var or flag
var (
	BotConfig 			*config.Config
	GuildID				string 				= ""
	RemoveCommands		bool				= true
)

// init: read the config file (grab bot tokens/prefixes/configs etc.)
func init() {
	var err error
	BotConfig, err = config.ReadConfig()
	if err != nil {
		log.Fatalf("Failed reading bot config file: %v", err)
	}
}

var s *discordgo.Session

// init: create new discord session
func init() {
	var err error
	s, err = discordgo.New("Bot " + BotConfig.Token)
	if err != nil {
		log.Fatalf("Failed creating new Discord session: %v", err)
	}
}

// list of slash commands
var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name: "ping",
			Description: "Send a ping message",
		},
		{
			Name: "echo",
			Description: "Echoes your message",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:			discordgo.ApplicationCommandOptionString,
					Name:			"message",
					Description: 	"Message to echo back",
					Required:		true,
				},
			},
		},
		{
			Name: "serverinfo",
			Description: "Get server info for cobalt's API",
		},
	}

	commandHandlers = map[string]func (s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping": cmd.PingHandler,
		"echo": cmd.EchoHandler,
		"serverinfo": cmd.ServerInfoHandler,
	}
)

// init: register slash commands/add command handlers
func init() {
	s.AddHandler(func (s *discordgo.Session, i *discordgo.InteractionCreate) {
		// get the command name; look up the command handler for it; execute it
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	}) 
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Bot is logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		// TODO: if i do change GuildID to flags (flag.String), refer this as a pointer
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to stop bot instance!")
	<-stop

	if RemoveCommands {
		log.Println("Removing commands...")

		for _, v := range registeredCommands {	
			err := s.ApplicationCommandDelete(s.State.User.ID, GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}