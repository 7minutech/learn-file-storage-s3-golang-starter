package main

import (
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
)

func (cfg *apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error) {
	if video.VideoURL == nil {
		return video, nil
	}

	videoParts := strings.Split(*video.VideoURL, ",")
	if len(videoParts) != 2 {
		return video, nil
	}

	bucket := videoParts[0]
	key := videoParts[1]

	presignedURL, err := generatePresignedURL(cfg.s3Client, bucket, key, time.Minute*10)
	if err != nil {
		return video, err
	}

	video.VideoURL = aws.String(presignedURL)

	return video, nil
}
