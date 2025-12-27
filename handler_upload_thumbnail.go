package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	const maxMemory = 10 << 20

	if err := r.ParseMultipartForm(maxMemory); err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse multipart form", err)
		return
	}

	file, header, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse form file", err)
		return
	}

	defer file.Close()

	mediaType := header.Header.Get("Content-Type")

	if mediaType == "" {
		respondWithError(w, http.StatusBadRequest, "Missing Content-Type", nil)
		return
	}

	mediaType, _, err = mime.ParseMediaType(mediaType)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to parse media type", err)
		return
	}

	if mediaType != "image/png" && mediaType != "image/jpeg" {
		respondWithError(w, http.StatusBadRequest, "Invalid Content-Type", nil)
		return
	}

	dbVideo, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get video", err)
		return
	}

	if userID != dbVideo.UserID {
		respondWithError(w, http.StatusUnauthorized, "user is not creator of video", nil)
		return
	}

	prefix := "image/"
	if !strings.HasPrefix(mediaType, prefix) {
		respondWithError(w, http.StatusBadRequest, "Unsupported Content-Type", nil)
		return
	}
	fileExtension := mediaType[len(prefix):]

	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not read bytes into key", err)
		return
	}
	videoKey := base64.RawURLEncoding.EncodeToString(key)

	videoPath := videoKey + "." + fileExtension
	videoFilePath := filepath.Join(cfg.assetsRoot, videoPath)

	createdFile, err := os.Create(videoFilePath)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create video file", err)
		return
	}

	defer createdFile.Close()

	if _, err := io.Copy(createdFile, file); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to copy multipart File to destination File", err)
		return
	}

	var thumbnailURL = fmt.Sprintf("http://localhost:%s/assets/%s.%s", cfg.port, videoKey, fileExtension)

	dbVideo.ThumbnailURL = &thumbnailURL

	if err := cfg.db.UpdateVideo(dbVideo); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update video", err)
		return
	}

	respondWithJSON(w, http.StatusOK, dbVideo)
}
