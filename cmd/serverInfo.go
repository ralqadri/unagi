package cmd

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/lostdusty/gobalt"
)

func ServerInfoHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Recover from panic if anything fails; https://golang-id.org/blog/defer-panic-and-recover/
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovering from panic in `DownloadMediaHandler`: %v", r)
			content := "Unexpected error while processing your request! Maybe try again? (Panic on ServerInfoHandler)"
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &content,
			})
		}
	}()

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
