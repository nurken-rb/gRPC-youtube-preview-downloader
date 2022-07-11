package main

import (
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
	"previewer/api"
	"previewer/internal/config"
	storepkg "previewer/internal/storage"
	grpcserver "previewer/internal/transport/grpc"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("cannot get config to start app %v", err)
	}

	storage := storepkg.GetNewStorage()
	if err = storage.Open(&cfg); err != nil {
		log.Fatalf("cannot connect to database %v", err)
	}

	if err = storage.InitDB(); err != nil {
		log.Fatalf("cannot init database %v", err)
	}

	s := grpc.NewServer()
	srv := &grpcserver.GRPCServer{Storage: storage}

	api.RegisterPreviewerServer(s, srv)

	log.Infof("Application is up")

	go clearCacheDir(cfg.App.CacheDirPath, storage)

	lis, err := net.Listen("tcp", cfg.App.BindAddr)
	if err != nil {
		log.Fatal(err)
	}
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}

}

func clearCacheDir(previewCacheDir string, storage *storepkg.Storage) {
	for {
		names, err := ioutil.ReadDir(previewCacheDir)
		if err != nil {
			log.Warnf("reading files from cache dir ended with err %v", err)
		}

		for _, entry := range names {
			err = os.RemoveAll(path.Join([]string{previewCacheDir, entry.Name()}...))
			if err != nil {
				log.Warnf("error occured while deleting files for updating cache %v", err)
			}

			filename := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
			if err = storage.DeleteThumbnailRow(filename); err != nil {
				log.Warnf("error occured while deleting thumbnail row from db %v", err)
			}
		}

		time.Sleep(1 * time.Hour)
		log.Info("cache dir cleared")
	}
}
