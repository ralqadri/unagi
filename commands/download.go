package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
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

// func sanitizeFileName(filename string) string {
// 	ext := path.Ext(filename)
// 	if ext == "" {
// 		return filename
// 	}

// 	// TODO: continue removing query parameters

// 	// return ogName + ext // return ogName with the extension only
// }

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

		// log.Printf("resBody: %s", resBody)

		if resBody.Status == "success" || resBody.Status == "redirect" || resBody.Status == "stream" {
			// grab url here later and fetch it
			// maybe i should probably move this to another file but whatever
			// s.ChannelMessageSend(m.ChannelID, resBody.Url)

			outRes, err := http.Get(resBody.Url)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "error fetching file!: " + err.Error())
				log.Fatalf("error fetching file!: %s", err)
			}
			defer outRes.Body.Close()

			// TODO: this shit
			// filename := path.Base(resBody.Url)
			// if filename == "" {
			// 	filename = "file.mp4"
			// }
			
			// TODO: this shit
			// filename = sanitizeFileName(filename)

			filename := "file.mp4"

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

			log.Printf("file downloaded!: (%s) // filepath: %s", filename, filepath)
			defer out.Close()

			HandleSendFileCommand(s, m, prefix, content, filepath)

		} else {
			s.ChannelMessageSend(m.ChannelID, resBody.Text)
		}
	}
}