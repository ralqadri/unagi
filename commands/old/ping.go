package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandlePingCommand(s *discordgo.Session, m *discordgo.MessageCreate, prefix string, content string) {
	if strings.HasPrefix(content, prefix + "ping") {
		s.ChannelMessageSend(m.ChannelID, "pong!")
	}
}