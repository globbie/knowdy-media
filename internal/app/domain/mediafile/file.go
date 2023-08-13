package mediafile

import "time"

type FileChunk {
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
	CreateFile(fileName string, fileSize int, ownerId string, path string, mimetype string, createdAt time.Time) (MediaFile, error)
	ListFiles(owner string) ([]MediaFile, error)
}
