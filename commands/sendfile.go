// TODO: make this command into a tool/helper function instead

package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func HandleSendFileCommand(s *discordgo.Session, m *discordgo.MessageCreate, prefix string, content string, filepath string) {
	fmt.Printf("trying to send file: %s\n", filepath)
	file, err := os.Open(filepath);
	
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "error opening file!: " + err.Error())
		log.Fatalf("error opening file!: %s", err)
	}
	defer file.Close()
	
	// TODO: the shit about grabbing extension for the filepath
	_, err = s.ChannelFileSend(m.ChannelID, filepath, file)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "error sending file!: " + err.Error())
		log.Fatalf("error sending file!: %s", err)
	}
	fmt.Printf("successfully sent file: %s\n", filepath)
}