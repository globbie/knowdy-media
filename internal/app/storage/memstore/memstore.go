package memstore

import (
	"github.com/globbie/knowdy-media/internal/app/domain/mediafile"

	"strconv"
	"sync"
	"time"
)

type MemStore struct {
	Id             string
	filesById      map[string]mediafile.MediaFile
	filesByName    map[string][]mediafile.MediaFile
	filesByHash    map[string][]mediafile.MediaFile
	filesByOwner   map[string][]mediafile.MediaFile
	nextId         uint64
	mu             *sync.Mutex
}

func New(storeId string) *MemStore {
	return &MemStore{
	        Id:          storeId,
		filesById:   make(map[string]mediafile.MediaFile),
		filesByName: make(map[string][]mediafile.MediaFile),
		filesByHash: make(map[string][]mediafile.MediaFile),
		mu:          &sync.Mutex{},
	}
}

func (m *MemStore) CreateFile(mimeType string, fileName string, fileSize uint64, ownerId string, path string, createdAt time.Time) (mediafile.MediaFile, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	f := mediafile.MediaFile{
		MimeType:  mimeType,
		Id:        strconv.FormatUint(m.nextId, 16),
		Name:      fileName,
		Size:      fileSize,
		Owner:     ownerId,
		Path:      path,
		CreatedAt: createdAt,
	}

        m.filesById[f.Id] = f

        files, ok := m.filesByName[fileName]
	if !ok {
		m.filesByName[fileName] = make([]mediafile.MediaFile, 0, 1)
	}
	m.filesByName[fileName] = append(files, f)

	files, ok = m.filesByHash[f.Hash]
	if !ok {
		m.filesByHash[f.Hash] = make([]mediafile.MediaFile, 0, 1)
   	        m.filesByHash[f.Hash] = append(files, f)
	} else {
	   // implement doublet logic
	}

        m.nextId++;
        return f, nil
}

func (m *MemStore) ListFiles(ownerId string) ([]mediafile.MediaFile, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	files, ok := m.filesByOwner[ownerId]
	if !ok {
		return []mediafile.MediaFile{}, nil
	}
	return files, nil
}

