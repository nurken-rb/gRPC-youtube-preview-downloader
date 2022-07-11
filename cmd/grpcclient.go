package main

import (
    "context"
    "errors"
    "flag"
    "fmt"
    log "github.com/sirupsen/logrus"
    "google.golang.org/grpc/credentials/insecure"
    "os"
    "strings"
    "sync"

    "google.golang.org/grpc"
    "previewer/api"
)

func main() {
    conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("connection failed: %v", err)
    }
    defer conn.Close()

    c := api.NewPreviewerClient(conn)

    var async bool
    flag.BoolVar(&async, "async", false, "download multiple previews in async mode")
    flag.Parse()
    args := flag.Args()

    if async == true {
        args = args[1:]

        var wg sync.WaitGroup
        for _, url := range args {
            wg.Add(1)
            go func(url string) {
                previewImage := downloadPreview(c, url)
                filename, err := findVideoID(url)
                if err != nil {
                    log.Warnf("error during finding video id %v", err)
                    wg.Done()
                    return
                }
                if err := savePreviewAsFile(previewImage.Image, filename); err != nil {
                    log.Warnf("Error occured when saving preview %v", err)
                }
                wg.Done()
            }(url)
        }
        wg.Wait()
    } else {
        url := args[0]
        previewImage := downloadPreview(c, url)
        filename, err := findVideoID(url)
        if err != nil {
            log.Warnf("error during finding video id %v", err)
            return
        }
        if err := savePreviewAsFile(previewImage.Image, filename); err != nil {
            log.Warnf("Error occured when saving preview %v", err)
        }
    }
}

func downloadPreview(c api.PreviewerClient, videoURL string) *api.PreviewImage {
    previewUrl := api.PreviewUrl{
        VideoUrl: videoURL,
    }

    ctx := context.Background()
    previewImg, err := c.GetPreview(ctx, &previewUrl)

    if err != nil {
        log.Fatalf("request failed: %v", err)
    }

    return previewImg
}

func savePreviewAsFile(imageData []byte, filename string) error {
    fmt.Println("save with", filename)
    file, err := os.Create(fmt.Sprintf("./thumbnails_client/%v.jpg", filename))
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

func findVideoID(videoUrl string) (string, error) {
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
