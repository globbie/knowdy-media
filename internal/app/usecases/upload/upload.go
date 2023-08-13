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
	CreateFile(creatorId, roomId string, text string) error
}

type UseCases struct {
	FileStorage mediafile.Interface
}

func (uc *UseCases) CreateFile(ownerId string, fileName string, fileSize uint64) ([]MediaFileInfo, error) {
	t := time.Now()
	_, err := uc.MediaStorage.CreateFile(ownerId, roomId, text, t)
	if err != nil {
		return err
	}
	return nil
}
