package upload

import (
	"github.com/globbie/knowdy-media/internal/app/domain/mediafile"

	"time"
)

type MediaFileInfo struct {
	Name      string
	Owner     string
	Size      uint64
	MimeType  string
	CreatedAt time.Time
}

type Interface interface {
	CreateFile(mimeType string, fileName string, fileSize uint64, ownerId string, path string) (MediaFileInfo, error)
}

type UseCases struct {
	MediaStorage mediafile.Interface
}

func (uc *UseCases) CreateFile(mimeType string, fileName string, fileSize uint64, ownerId string, path string) (MediaFileInfo, error) {
	t := time.Now()
	_, err := uc.MediaStorage.CreateFile(mimeType, fileName, fileSize, ownerId, path, t)
	if err != nil {
		return MediaFileInfo{}, err
	}
	return MediaFileInfo{}, nil
}
