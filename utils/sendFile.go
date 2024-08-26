package utils

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

// TODO: failsafes; add a timeout and maybe also a filesize limit
// filesize limit ref: https://discord.com/developers/docs/reference#uploading-files
func SendFileToChannel(s *discordgo.Session, m *discordgo.MessageCreate, prefix string, content string, filepath string, filename string) {
	log.Printf("trying to send file: %s\n", filepath)
	file, err := os.Open(filepath);
	
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "error opening file!: " + err.Error())
		log.Fatalf("error opening file!: %s", err)
	}
	defer file.Close()
	
	_, err = s.ChannelFileSend(m.ChannelID, filename, file)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "error sending file!: " + err.Error())
		log.Fatalf("error sending file!: %s", err)
	} else {
		log.Printf("successfully sent file: %s\n", filepath)
	}

	err = os.Remove(filepath)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "error deleting file!: " + err.Error())
		log.Fatalf("error deleting file!: %s", err)
	} else {
		log.Printf("successfully deleted file: %s\n", filepath)	
	}
}