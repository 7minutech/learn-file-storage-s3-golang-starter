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

func calculateAspectRatio(width, height int) string {
	gcf := greatestCommonFactor(width, height)

	return fmt.Sprintf("%d:%d", width/gcf, height/gcf)
}

func greatestCommonFactor(a, b int) int {
	factor := a
	if a > b {
		factor = b
	}
	for a%factor != 0 || b%factor != 0 {
		if a%factor == 0 || b%factor == 0 {
		}
		factor -= 1
	}
	return factor

}

func getVideoAspectRatio(filePath string) (string, error) {

	args := []string{"-v", "error", "-print-format", "json", "-show-streams", filePath}

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
