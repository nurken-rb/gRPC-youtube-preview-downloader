package thumbnail

import (
    "testing"
)

func TestFindVideoID(t *testing.T) {
    tmb := NewThumbnail()
    urlList := []string{"https://www.youtube.com/watch?v=ZAMicELbvPc&ab_channel=%D0%96%D0%98%D0%97%D0%9D%D0%AC%D0%A4%D0%A3%D0%A2%D0%91%D0%9E%D0%9B",
        "https://www.youtube.com/watch?v=s-7pyIxz8Qg&ab_channel=MovieclipsClassicTrailers"}
    videoIDList := []string{"ZAMicELbvPc", "s-7pyIxz8Qg"}

    for ind, url := range urlList {
        id, err := tmb.findVideoID(url)
        if err != nil {
            t.Errorf("Must be nil return! err = %v", err)
        }

        if id != videoIDList[ind] {
            t.Errorf("Found video ID is wrong. ID = %s, URL = %s", id, url)
        }
    }
}

func TestFindVideoIDError(t *testing.T) {
    tmb := NewThumbnail()
    urlList := []string{"invalid-url"}

    for _, url := range urlList {
        _, err := tmb.findVideoID(url)
        if err == nil {
            t.Errorf("Must be non nil return! err = %v", err)
        }
    }
}
