package memstore

import (
	mf "github.com/globbie/knowdy-media/internal/app/domain/mediafile"

	"errors"
	"sync"
	"time"
)

type MemStore struct {
	Id             string
	filesById      map[string]mf.MediaFile
	filesByName    map[string][]mf.MediaFile
	filesByCheckSum    map[string][]mf.MediaFile
	filesByOwner   map[string][]mf.MediaFile
	mu             *sync.Mutex
}

func New(storeId string) *MemStore {
	return &MemStore{
	        Id:          storeId,
		filesById:   make(map[string]mf.MediaFile),
		filesByName: make(map[string][]mf.MediaFile),
		filesByCheckSum: make(map[string][]mf.MediaFile),
		mu:          &sync.Mutex{},
	}
}

func (m *MemStore) SaveFileMeta(mimeType, fileId, fileName string, fileSize uint64,
	checkSum string, ownerId string) (mf.MediaFile, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	f := mf.MediaFile{
		MimeType:  mimeType,
		Id:        fileId,
		Name:      fileName,
		Size:      fileSize,
		CheckSum:  checkSum,
		Owner:     ownerId,
		CreatedAt: time.Now(),
	}

        m.filesById[f.Id] = f

        files, ok := m.filesByName[fileName]
	if !ok {
		m.filesByName[fileName] = make([]mf.MediaFile, 0, 1)
	}
	m.filesByName[fileName] = append(files, f)

	files, ok = m.filesByCheckSum[f.CheckSum]
	if !ok {
		m.filesByCheckSum[f.CheckSum] = make([]mf.MediaFile, 0, 1)
   	        m.filesByCheckSum[f.CheckSum] = append(files, f)
	} else {
		// thread safe doublet warning
		return mf.MediaFile{}, errors.New("File already exists")
	}
        return f, nil
}

func (m *MemStore) CheckDoublets(checkSum string) ([]mf.MediaFile, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	files, ok := m.filesByCheckSum[checkSum]
	if ok {
		// return existing file copies
		return files, errors.New("File already exists")
	}
	return []mf.MediaFile{}, nil
}

func (m *MemStore) ListFiles(ownerId string) ([]mf.MediaFile, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	files, ok := m.filesByOwner[ownerId]
	if !ok {
		return []mf.MediaFile{}, nil
	}
	return files, nil
}


