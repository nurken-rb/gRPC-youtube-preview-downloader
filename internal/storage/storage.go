package storage

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "previewer/internal/config"
)

func GetNewStorage() *Storage {
    return &Storage{}
}

type Storage struct {
    db *sql.DB
}

func (s *Storage) Open(cfg *config.Config) error {
    dsn := fmt.Sprintf("host=%s dbname=%s sslmode=%s user=%s password=%s port=%s",
        cfg.Database.Host, cfg.Database.DbName, cfg.Database.SSLMode,
        cfg.Database.User, cfg.Database.Password, cfg.Database.Port)
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return err
    }

    if err := db.Ping(); err != nil {
        return err
    }

    s.db = db
    return nil
}

func (s *Storage) InitDB() error {
    _, err := s.db.Query("CREATE TABLE IF NOT EXISTS  thumbnails (videoID varchar(120) primary key not null , videoURL varchar(255) not null unique, path varchar(120) not null unique)")
    return err
}

func (s *Storage) SaveThumbnail(videoID string, videoURL string, path string) error {
    _, err := s.db.Exec("INSERT INTO thumbnails (videoID, videoURL, path) values ($1, $2, $3)", videoID, videoURL, path)
    return err
}

func (s *Storage) DeleteThumbnailRow(videoID string) error {
    _, err := s.db.Exec("DELETE FROM thumbnails where videoID = $1", videoID)
    return err
}

func (s *Storage) GetPreviewPath(videoID string) (string, error) {
    var previewPath string
    if err := s.db.QueryRow("select path from thumbnails where videoID = $1", videoID).Scan(&previewPath); err != nil {
        return "", err
    }

    return previewPath, nil
}
