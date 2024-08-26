package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ralqadri/unagi/utils"
)

type RequestBody struct {
	Url		string	`json:"url"`
}

type ResponseBody struct {
	Status 		string	`json:"status"`
	Url			string	`json:"url"`
	Text		string	`json:"text"`
}

// https://www.practical-go-lessons.com/post/go-how-to-send-post-http-requests-with-a-json-body-cbhvuqa220ds70kp2lkg

func HandleDownloadCommand(s *discordgo.Session, m *discordgo.MessageCreate, prefix string, content string) {
	if strings.HasPrefix(content, prefix + "dl") {
		downloadLink := strings.Trim(content, prefix + "dl")
		log.Printf("\nstarting to download: %s", downloadLink)

		// preparing the json for the request body
		reqBody := RequestBody{
			Url: downloadLink,
		}

		marshaled, err := json.Marshal(reqBody)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "impossible to marshal request body!: " + err.Error())
			log.Printf("impossible to marshal request body!: %s", err)
			return
		}

		// preparing request body for the POST request
		req, err := http.NewRequest("POST", "https://api.cobalt.tools/api/json", bytes.NewReader(marshaled))
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "error creating request!: " + err.Error())
			log.Printf("can't build request!: %s", err)
			return
		}

		req.Header = http.Header{
			"Accept":  {"application/json"},
			"Content-Type": {"application/json"},
		}

		client := http.Client{Timeout: 15 * time.Second}

		res, err := client.Do(req)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "error fetching cobalt api!: " + err.Error())
			log.Printf("error fetching cobalt api!: %s", err)
			return
		}
		defer res.Body.Close()

		// read body
		body, err := io.ReadAll(res.Body)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "error reading response body!: " + err.Error())
			log.Printf("error reading response body!: %s", err)
			return
		}

		log.Printf("resBody: %s", string(body))

		var resBody ResponseBody
		err = json.Unmarshal(body, &resBody)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "error unmarshaling response body!: " + err.Error())
			log.Printf("error unmarshaling response body!: %s", err)
			return
		}

		if resBody.Status == "success" || resBody.Status == "redirect" || resBody.Status == "stream" {
			// TODO: file download failsafes
			outRes, err := http.Get(resBody.Url)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "error fetching file!: " + err.Error())
				log.Printf("error fetching file!: %s", err)
				return
			}
			defer outRes.Body.Close()

			filename := ""
			if resBody.Status == "stream" {
				// TODO: a better way to fetch the title of the stream (and potentially for stuff that are also "success" and "redirect" responses)
				filename = "stream.mp4"
			} else {
				filename = utils.SanitizeFileName(path.Base(resBody.Url))
				log.Printf("returned filename: %s\n", filename)
			}

			filepath := fmt.Sprintf("%s/%s", "./downloads", filename)
			out, err := os.Create(filepath)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "error creating file!: " + err.Error())
				log.Printf("error creating file!: %s", err)
				return
			}

			_, err = io.Copy(out, outRes.Body)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "error copying file!: " + err.Error())
				log.Printf("error copying file!: %s", err)
				return
			}
			log.Printf("file downloaded!: %s // filepath: %s", filename, filepath)
			defer out.Close()

			utils.SendFileToChannel(s, m, prefix, content, filepath, filename)
			utils.CleanUpFile(filepath)
		} else {
			log.Printf("cobalt failed to process the link: %s", resBody.Text)
			s.ChannelMessageSend(m.ChannelID, resBody.Text)
			return
		}

	}
}