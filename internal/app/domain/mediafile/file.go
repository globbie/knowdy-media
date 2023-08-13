package mediafile

import "time"

type FileChunk struct {
	StorageId string
	Id        string
	Size      uint64
}

type MediaFile struct {
	Id        string
	Name      string
	Owner     string
	Path      string
	Size      uint64
	Hash      string
	MimeType  string
	CreatedAt time.Time
	IsChunk   bool
	Chunks    []FileChunk
}

type Interface interface {
	CreateFile(mimetype string, fileName string, fileSize uint64, ownerId string, path string, createdAt time.Time) (MediaFile, error)
	ListFiles(owner string) ([]MediaFile, error)
}
