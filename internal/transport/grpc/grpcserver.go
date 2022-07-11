package grpcserver

import (
    "context"
    log "github.com/sirupsen/logrus"
    "google.golang.org/grpc"
    "previewer/api"
    "previewer/internal/storage"
    "previewer/internal/thumbnail"
)

type GRPCServer struct {
    api.UnimplementedPreviewerServer
    srv     *grpc.Server
    Storage *storage.Storage
}

func (s *GRPCServer) GetPreview(_ context.Context, req *api.PreviewUrl) (*api.PreviewImage, error) {
    log.Infof("get request for preview with url %v", req.VideoUrl)
    tmb := thumbnail.NewThumbnail()

    body, err := tmb.DownloadThumbnail(req.VideoUrl, s.Storage)
    if err != nil {
        log.Warnf("downloading preview is failed with error %v", err)
        return &api.PreviewImage{Image: []byte{}}, err
    }

    return &api.PreviewImage{Image: body}, nil
}
