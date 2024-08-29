package cmd

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/lostdusty/gobalt"
)

func ServerInfoHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	server, err := gobalt.CobaltServerInfo(gobalt.CobaltApi)

	if err != nil {
		log.Printf("Failed to get server info: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: fmt.Sprintf("Failed to get server info: %v", err),
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprintf(`
				> Version : %v
				> Commit: %v
				> Branch: %v
				> Name: %v
				> URL: %v
				> CORS: %v
				> Start time: %v
				`, server.Version, server.Commit, server.Branch, server.Name, server.URL, server.Cors, server.StartTime),
		},
	})
}
