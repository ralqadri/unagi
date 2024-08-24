package commands

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	cobaltApiUrl = "api.cobalt.tools"
)

// https://stackoverflow.com/posts/68018927/revisions

func HandleDownloadCommand(s *discordgo.Session, m *discordgo.MessageCreate, prefix string, content string) {
	if strings.HasPrefix(content, prefix + "dl") {
		client := http.Client{}
		req, err := http.NewRequest("GET", "https://api.cobalt.tools/api/serverInfo", nil)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "error creating request!: " + err.Error())
			fmt.Println("error creating request!: ", err)
		}

		req.Header = http.Header{
			"Accept":  {"application/json"},
			"Content-Type": {"application/json"},
		}

		res, err := client.Do(req)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "error fetching cobalt api!: " + err.Error())
			fmt.Println("error fetching cobalt api!: ", err)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "error reading response body!: " + err.Error())
			fmt.Println("error reading response body!: ", err)
			return
		}
		
		wrappedBody := fmt.Sprintf("```json\n%s\n```", string(body))
		
		s.ChannelMessageSend(m.ChannelID, wrappedBody)
		fmt.Println("body: ", string(body))
	}
}