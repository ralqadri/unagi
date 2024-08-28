package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandleEchoCommand(s *discordgo.Session, m *discordgo.MessageCreate, prefix string, content string) {
	if strings.HasPrefix(content, prefix + "echo") {
		underlyingContent := strings.Trim(content, prefix + "echo")
		s.ChannelMessageSend(m.ChannelID, underlyingContent)
	}
}