package mediafile

import (
	"time"
)

type FileChunk struct {
	StorageId string
	Id        string
	Size      uint64
}

type MediaFile struct {
	Id        string
	Name      string
	Owner     string
	Size      uint64
	CheckSum  string
	MimeType  string
	CreatedAt time.Time
	IsChunk   bool
	Chunks    []FileChunk
}

type FileMetaSaver interface {
	SaveFileMeta(mimetype string, fileId string, fileName string, fileSize uint64, CheckSum string, ownerId string) (MediaFile, error)
}

type FileMetaQuery interface {
	CheckDoublets(checkSum string) ([]MediaFile, error)
	ListFiles(owner string) ([]MediaFile, error)
}

