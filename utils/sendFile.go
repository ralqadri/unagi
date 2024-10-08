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
		log.Printf("error opening file!: %s", err)
		return
	}
	defer file.Close()

	fileInfo, err :=  os.Stat(filepath)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "error getting file info!: " + err.Error())
		log.Printf("error getting file info!: %s", err)
		return
	}

	if fileInfo.Size() > 26214400 {
		s.ChannelMessageSend(m.ChannelID, "file is too big to send! (max 25MB)")
		log.Printf("file is too big to send! (max 25MB): %s", filepath)
	} else {
		_, err = s.ChannelFileSend(m.ChannelID, filename, file)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "error sending file!: " + err.Error())
			log.Printf("error sending file!: %s", err)
			return
		} else {
			log.Printf("successfully sent file: %s\n", filepath)
		}
	}
}