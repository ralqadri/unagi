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
		
		// getting the download link from the message content
		downloadLink := strings.Trim(content, prefix + "dl")
		log.Printf("\nstarting to download: %s", downloadLink)

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
		// log.Printf("status code: %d", res.StatusCode)

		// close body to free resources; defer will execute this at the end of this current func
		defer res.Body.Close()

		// read body
		body, err := io.ReadAll(res.Body)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "error reading response body!: " + err.Error())
			log.Fatalf("error reading response body!: %s", err)
		}

		log.Printf("resBody: %s", string(body))

		var resBody ResponseBody
		err = json.Unmarshal(body, &resBody)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "error unmarshaling response body!: " + err.Error())
			log.Fatalf("error unmarshaling response body!: %s", err)
		}

		if resBody.Status == "success" || resBody.Status == "redirect" || resBody.Status == "stream" {
			// TODO: file download failsafes
			outRes, err := http.Get(resBody.Url)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "error fetching file!: " + err.Error())
				log.Fatalf("error fetching file!: %s", err)
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
				log.Fatalf("error creating file!: %s", err)
			}

			_, err = io.Copy(out, outRes.Body)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "error copying file!: " + err.Error())
				log.Fatalf("error copying file!: %s", err)
			}

			log.Printf("file downloaded!: %s // filepath: %s", filename, filepath)
			
			fileInfo, err :=  os.Stat(filepath)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "error getting file info!: " + err.Error())
				log.Fatalf("error getting file info!: %s", err)
			}

			defer out.Close()

			if fileInfo.Size() > 26214400 {
				s.ChannelMessageSend(m.ChannelID, "file is too big to send! (max 25MB)")
				log.Fatalf("file is too big to send! (max 25MB): %s", filepath)
			} else {
				utils.SendFileToChannel(s, m, prefix, content, filepath, filename)
			}
			utils.CleanUpFile(filepath)

		} else {
			s.ChannelMessageSend(m.ChannelID, resBody.Text)
		}
	}
}