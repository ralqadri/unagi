package utils

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
)

func DownloadFile(url string) (string, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	filename := GetFileName(resp)
	filepath := fmt.Sprintf("%s/%s", "./downloads", filename)

	out, err := os.Create(filepath)
	if err != nil {
		return "", "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return filename, filepath, err
}

func GetFileName(resp *http.Response) string {
	// get filename from Content-Disposition header
	filename := resp.Header.Get("Content-Disposition")
	if filename != "" {
		_, params, err := mime.ParseMediaType(filename)
		log.Printf("GetFileName() // Params: %v", params)
		if err == nil {
			filename = params["filename"]
			return filename
		}
	}

	filename = path.Base(resp.Request.URL.Path)

	filename = SanitizeFileName(filename)

	if filename == "" {
		filename = "downloadedFile"
	}

	return filename
}
