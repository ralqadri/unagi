package cmd

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/lostdusty/gobalt"
)

func DownloadMediaHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	downloadMedia := gobalt.CreateDefaultSettings()
	downloadMedia.Url = i.ApplicationCommandData().Options[0].StringValue()

	var returnedContent string

	destination, err := gobalt.Run(downloadMedia)
	if err != nil {
		log.Printf("Failed to download media! URL: %v // %v", downloadMedia.Url, err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: fmt.Sprintf("Failed to download media! URL: %v // %v", downloadMedia.Url, err),
			},
		})
		return
	}

	if destination.Status == "error" || destination.Status == "rate-limit" {
		log.Printf("Failed to download media! URL: %v // %v", downloadMedia.Url, destination.Status)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: fmt.Sprintf("Failed to download media! Cobalt returned status: `%v`", destination.Status),
			},
		})
		return
	}

	// TODO: Refactor this to use SelectMenu component instead... but it's a fucking pain in the ass
	// This implementation is a compromise for now
	if destination.Status == "picker" {
		carouselIndex := i.ApplicationCommandData().Options[1].IntValue() - 1

		if int(carouselIndex) > len(destination.URLs) || int(carouselIndex) < 0 {
			log.Printf("Index out of bounds! Index picked is %v", carouselIndex)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: fmt.Sprintf("There's no media #%v, silly!", carouselIndex+1),
				},
			})
			return
		}

		returnedContent = destination.URLs[carouselIndex]
	}

	if destination.Status == "success" || destination.Status == "stream" || destination.Status == "redirect" {
		returnedContent = destination.URL
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: returnedContent,
		},
	})

}
