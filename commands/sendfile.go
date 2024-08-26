// TODO: make this command into a tool/helper function instead

package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

// TODO: failsafes; add a timeout and maybe also a filesize limit
// filesize limit ref: https://discord.com/developers/docs/reference#uploading-files
func HandleSendFileCommand(s *discordgo.Session, m *discordgo.MessageCreate, prefix string, content string, filepath string, filename string) {
	fmt.Printf("trying to send file: %s\n", filepath)
	file, err := os.Open(filepath);
	
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "error opening file!: " + err.Error())
		log.Fatalf("error opening file!: %s", err)
	}
	defer file.Close()
	
	// TODO: the shit about grabbing extension for the filepath
	_, err = s.ChannelFileSend(m.ChannelID, filename, file)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "error sending file!: " + err.Error())
		log.Fatalf("error sending file!: %s", err)
	}
	fmt.Printf("successfully sent file: %s\n", filepath)
}