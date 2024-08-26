package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
			log.Fatalf("error fetching cobalt api!: %s", err)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "error reading response body!: " + err.Error())
			log.Fatalf("error reading response body!: %s", err)
		}

		var serverInfo ServerInfo
		err = json.Unmarshal(body, &serverInfo)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "error unmarshaling response body!: " + err.Error())
			log.Fatalf("error unmarshaling response body!: %s", err)
		}
		fmt.Println("serverInfo: ", serverInfo)
		
		wrappedBodyMessage := fmt.Sprintf("```json\n%s\n```", string(body))
		
		s.ChannelMessageSend(m.ChannelID, wrappedBodyMessage)
		fmt.Println("body: ", string(body))
	}
}