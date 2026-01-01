package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

type videoMetaData struct {
	Streams []stream `json:"streams"`
}

type stream struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

const landscape = float64(16) / float64(9)
const portrait = float64(9) / float64(16)
const landscapeStr = "16:9"
const portraitStr = "9:16"

func calculateAspectRatio(width, height int) string {

	aspectRatio := float64(width) / float64(height)

	if landscape-0.2 < float64(aspectRatio) && float64(aspectRatio) < landscape+0.2 {
		return landscapeStr
	} else if portrait-0.2 < float64(aspectRatio) && float64(aspectRatio) < portrait+0.2 {
		return portraitStr
	} else {
		return "other"
	}

}

func getVideoAspectRatio(filePath string) (string, error) {

	args := []string{"-v", "error", "-print_format", "json", "-show_streams", filePath}

	cmd := exec.Command("ffprobe", args...)

	var output = bytes.Buffer{}

	cmd.Stdout = &output

	if err := cmd.Run(); err != nil {
		log.Println("failed to run command for aspect ratio")
		return "", err
	}

	data := output.Bytes()

	var videoMetaData videoMetaData

	if err := json.Unmarshal(data, &videoMetaData); err != nil {
		log.Println("falied to unmarshal video meta data")
		return "", err
	}

	if len(videoMetaData.Streams) == 0 {
		log.Println("video metadata did not contain any streams")
		return "", fmt.Errorf("error: video meta data had no streams")
	}

	height := videoMetaData.Streams[0].Height
	width := videoMetaData.Streams[0].Width

	aspectRatio := calculateAspectRatio(width, height)

	if aspectRatio != "16:9" && aspectRatio != "9:16" {
		aspectRatio = "other"
	}

	return aspectRatio, nil

}
