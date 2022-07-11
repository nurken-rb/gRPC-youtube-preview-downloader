package thumbnail

import (
    "errors"
    "fmt"
    log "github.com/sirupsen/logrus"
    "io"
    "net/http"
    "os"
    "previewer/internal/storage"
    "strings"
)

func NewThumbnail() *thumbnail {
    return &thumbnail{thumbnailDir: "./thumbnails_server"}
}

type thumbnail struct {
    videoURL     string
    videoID      string
    thumbnailDir string
    filename     string
}

const viUrl = "https://img.youtube.com/vi/"
const previewResolution = "/sddefault.jpg"

func (t *thumbnail) DownloadThumbnail(videoURL string, storage *storage.Storage) ([]byte, error) {
    videoID, err := t.findVideoID(videoURL)
    if err != nil {
        return make([]byte, 0, 0), err
    }
    t.videoURL = videoURL
    t.videoID = videoID
    t.filename = fmt.Sprintf("%v.jpg", t.videoID)
    isDownloadedAlready := t.checkIfAlreadyDownloaded()

    if isDownloadedAlready {
        previewPath, err := storage.GetPreviewPath(t.videoID)
        if err != nil {
            return make([]byte, 0, 0), err
        }

        file, err := os.Open(previewPath)
        if err != nil {
            return make([]byte, 0, 0), err
        }
        defer file.Close()

        fileInfo, _ := file.Stat()
        fileSize := fileInfo.Size()
        body := make([]byte, fileSize)
        if _, err = file.Read(body); err != nil {
            return make([]byte, 0, 0), err
        }

        log.Infof("preview was taken from db cache")
        return body, nil
    }

    resp, err := http.Get(viUrl + videoID + previewResolution)
    if err != nil {
        return make([]byte, 0, 0), err
    }

    body, err := io.ReadAll(resp.Body)

    if err != nil {
        log.Warnf("reading response after downloading preview is failed with error %v", err)
        return make([]byte, 0, 0), err
    }

    if err = t.savePreviewAsFile(body); err != nil {
        log.Warnf("saving preview is failed with error %v", err)
        return make([]byte, 0, 0), err
    }

    previewPath := fmt.Sprintf("%v/%v", t.thumbnailDir, t.filename)
    if err = storage.SaveThumbnail(t.videoID, t.videoURL, previewPath); err != nil {
        log.Warnf("cannot save thumbnail in storage %v", err)
    }

    return body, nil
}

func (t *thumbnail) findVideoID(videoUrl string) (string, error) {

    equalityIndex := strings.Index(videoUrl, "=")
    ampIndex := strings.Index(videoUrl, "&")

    var videoId string
    if equalityIndex != -1 {
        if ampIndex != -1 {
            videoId = videoUrl[equalityIndex+1 : ampIndex]
        } else {
            return "", errors.New("invalid video url, cannot extract video id")
        }
    } else {
        return "", errors.New("invalid video url, cannot extract video id")
    }

    return videoId, nil
}

func (t *thumbnail) savePreviewAsFile(imageData []byte) error {
    file, err := os.Create(fmt.Sprintf("./%v/%v", t.thumbnailDir, t.filename))
    if err != nil {
        log.Warnf("error during saving while %v", err)
        return err
    }
    defer file.Close()

    if _, err = file.Write(imageData); err != nil {
        log.Warnf("error during writing response body to file %v", err)
        return err
    }

    return nil
}

func (t *thumbnail) checkIfAlreadyDownloaded() bool {
    path := fmt.Sprintf("%v/%v", t.thumbnailDir, t.filename)
    isFileExists, err := checkIfExists(path)

    if err != nil {
        return false
    }

    return isFileExists
}

func checkIfExists(name string) (bool, error) {
    _, err := os.Stat(name)
    if err == nil {
        return true, nil
    }
    if errors.Is(err, os.ErrNotExist) {
        return false, nil
    }
    return false, err
}
