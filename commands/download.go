package commands

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type RequestBody struct {
	Url		string	`json:"url"`
}

// https://www.practical-go-lessons.com/post/go-how-to-send-post-http-requests-with-a-json-body-cbhvuqa220ds70kp2lkg

func HandleDownloadCommand(s *discordgo.Session, m *discordgo.MessageCreate, prefix string, content string) {
	if strings.HasPrefix(content, prefix + "dl") {
		
		// getting the download link from the message content
		downloadLink := strings.Trim(content, prefix + "dl")
		log.Printf("starting to download: %s", downloadLink)

		// preparing the json for the request body
		reqBody := RequestBody{
			Url: downloadLink,
		}

		marshalled, err := json.Marshal(reqBody)
		if err != nil {
			log.Fatalf("impossible to marshall request body!: %s", err)
		}

		// preparing request body for the POST request
		req, err := http.NewRequest("POST", "https://api.cobalt.tools/api/json", bytes.NewReader(marshalled))
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "error creating request!: " + err.Error())
			log.Fatalf("can't build request!: %s", err)
		}

		req.Header = http.Header{
			"Accept":  {"application/json"},
			"Content-Type": {"application/json"},
		}

		client := http.Client{Timeout: 15 * time.Second}

		res, err := client.Do(req)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "error fetching cobalt api!: " + err.Error())
			log.Fatalf("error fetching cobalt api!: %s", err)
		}
		log.Printf("status code: %d", res.StatusCode)

		// close body to free resources; defer will execute this at the end of this current func
		defer res.Body.Close()

		// read body
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "error reading response body!: " + err.Error())
			log.Fatalf("error reading response body!: %s", err)
		}
		s.ChannelMessageSend(m.ChannelID, string(resBody))
		log.Printf("resBody: %s", string(resBody))
	}
}