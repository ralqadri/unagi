package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/lostdusty/gobalt"
	"github.com/ralqadri/unagi/utils"
)

const FilesizeLimit int64 = 26214400 // 25MB

func DownloadMediaHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Recover from panic if anything fails; https://golang-id.org/blog/defer-panic-and-recover/
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovering from panic in `DownloadMediaHandler`: %v", r)
			content := "Unexpected error while processing your request! Maybe try again? (Panic on DownloadMediaHandler)"
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &content,
			})
		}
	}()

	downloadMedia := gobalt.CreateDefaultSettings()
	downloadMedia.Url = i.ApplicationCommandData().Options[0].StringValue()

	var returnedUrl string

	user := i.Interaction.User
	if user == nil {
		user = i.Interaction.Member.User
	}

	if user != nil {
		log.Printf("User %v (%v) requested to download %v", user.Username, user.ID, downloadMedia.Url)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Starting the download...",
		},
	})

	destination, err := gobalt.Run(downloadMedia)
	if err != nil {
		log.Printf("Failed to download media! URL: %v // %v", downloadMedia.Url, err)

		content := fmt.Sprintf("Failed to download media! \nError: `%v`", err)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &content,
		})
		return
	}

	if destination.Status == "error" || destination.Status == "rate-limit" {
		log.Printf("Failed to download media! URL: %v // %v", downloadMedia.Url, destination.Status)

		content := fmt.Sprintf("Failed to download media! \nCobalt returned status: `%v` \nError: %v", destination.Status, destination.Text)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &content,
		})
		return
	}

	// TODO: Refactor this to use SelectMenu component instead... but it's a fucking pain in the ass
	// This implementation is a compromise for now
	if destination.Status == "picker" {
		carouselIndex := i.ApplicationCommandData().Options[1].IntValue() - 1

		if int(carouselIndex) > len(destination.URLs) || int(carouselIndex) < 0 {
			log.Printf("Index out of bounds! Index picked is %v", carouselIndex)

			content := fmt.Sprintf("There's no media #%v, silly!", carouselIndex+1)
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &content,
			})

			return
		}

		// returnedUrl = destination.URLs[carouselIndex]
	}

	// TODO: Change this shit to followup messages
	if destination.Status == "success" || destination.Status == "stream" || destination.Status == "redirect" {
		returnedUrl = destination.URL

		log.Printf("Trying to download %v", returnedUrl)
		filename, filepath, err := utils.DownloadFile(returnedUrl)
		if err != nil {
			log.Printf("Failed to download file: %v", err)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: fmt.Sprintf("Failed to download file: %v", err),
				},
			})
			return
		}

		file, err := os.Open(filepath)
		if err != nil {
			log.Printf("Failed to open file %v: %v", filepath, err)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: fmt.Sprintf("Failed to open file %v: %v", filepath, err),
				},
			})
			return
		}
		defer file.Close()

		fileInfo, err := os.Stat(filepath)
		if err != nil {
			log.Printf("Failed to get file %v info: %v", filepath, err)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: fmt.Sprintf("Failed to get file %v info: %v", filepath, err),
				},
			})
			return
		}

		// TODO: Change to a constant
		if fileInfo.Size() > FilesizeLimit {
			log.Printf("File `%v` is too big to send! (Max file size is 25MB)", filename)

			content := fmt.Sprintf("File `%v` is too big to send! (Max file size is 25MB)", filename)
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &content,
			})

			return
		} else {
			log.Printf("File %v is ready!", filepath)

			content := fmt.Sprintf("File `%v` is ready!", filename)
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &content,
				Files: []*discordgo.File{
					{
						Name:        fileInfo.Name(),
						ContentType: "application/octet-stream",
						Reader:      file,
					},
				},
			})

			if user != nil {
				log.Printf("File %v should be sent to user %v (%v) now!", filepath, user.Username, user.ID)
			}
		}
	}

}
