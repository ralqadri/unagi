package cmd

import (
	"github.com/bwmarrin/discordgo"
)

func EchoHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: i.ApplicationCommandData().Options[0].StringValue(),
		},
	})
}