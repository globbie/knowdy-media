package upload

import (
	"github.com/globbie/knowdy-media/internal/app/domain/mediafile"

	"github.com/gofrs/uuid"
	"github.com/codingsince1985/checksum"

	"io"
	"log"
	"os"
	"time"
)

const (
        uploadPath = "/tmp"
)

type MediaFileInfo struct {
	MimeType  string
	Id        string
	Name      string
	Size      uint64
	CheckSum  string
	Owner     string
	CreatedAt time.Time
}

type FileSaver interface {
	SaveFile(src io.Reader, mimeType string, fileName string, fileSize uint64, ownerId string) (MediaFileInfo, error)
}

type FileStorage struct {
	FileMetaSaver mediafile.FileMetaSaver
	FileMetaQuery mediafile.FileMetaQuery
}

func (fs *FileStorage) SaveFileBody(src io.Reader, mimeType string) (string, string, error) {
	fileId, err := uuid.NewV4()
	if err != nil {
		log.Printf("failed to generate file UUID: %v", err)
		return "", "", err
	}
	log.Printf("generated Version 4 file UUID %v", fileId.String())

	filePath := uploadPath + "/" + fileId.String()
	out, err := os.Create(filePath)
	if err != nil {
		log.Printf("failed to open the file %s for writing", filePath)
		return "", "", err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		log.Printf("file copy failed: %s\n", err)
		return "", "", err
	}

	return fileId.String(), filePath, err
}

func (fs *FileStorage) SaveFile(src io.Reader, mimeType string,
	fileName string, fileSize uint64, ownerId string) (MediaFileInfo, error) {

	fileId, filePath, err := fs.SaveFileBody(src, mimeType)
	if err != nil {
		log.Printf("-- content copying failed: %v", fileId)

		return MediaFileInfo{}, err
	}
	
	checkSum, err := checksum.SHA256sum(filePath)
	if err != nil {
		log.Printf("-- file checksum failed: %v", err)
		return MediaFileInfo{}, err
	}

	checkSum = "000"
	// preliminary doublet checking,
	// concurrently running transactions
	// might still introduce new ones (to be resolved later)
	files, err := fs.FileMetaQuery.CheckDoublets(checkSum)
	if err != nil {
		log.Printf("-- file doublet detected: %v", checkSum)
		for _, f := range files {
			log.Printf("-- file doublet detected: %v", f.Id)
		}
		return MediaFileInfo{}, ErrAlreadyExists
	}

	mf, err := fs.FileMetaSaver.SaveFileMeta(mimeType, fileId, fileName, fileSize,
		checkSum, ownerId)
	if err != nil {
		log.Printf("-- metadata saving failed: %v", fileId)
		return MediaFileInfo{}, err
	}
	return MediaFileInfo{mimeType, fileId, fileName, fileSize,
		checkSum, ownerId, mf.CreatedAt}, nil
}
