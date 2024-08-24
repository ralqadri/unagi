package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type ServerInfo struct {
	Name		string 	`json:"name"`
	StartTime	string 	`json:"startTime"`
}


// https://stackoverflow.com/posts/68018927/revisions

func HandleServerInfoCommand(s *discordgo.Session, m *discordgo.MessageCreate, prefix string, content string) {
	if strings.HasPrefix(content, prefix + "info") {
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

		var serverInfo ServerInfo
		err = json.Unmarshal(body, &serverInfo)
		if err != nil {
			fmt.Println("error unmarshaling response body!: ", err)
			return
		}
		
		wrappedBodyMessage := fmt.Sprintf("```json\n%s\n```", string(body))
		
		s.ChannelMessageSend(m.ChannelID, wrappedBodyMessage)
		fmt.Println("body: ", string(body))
	}
}