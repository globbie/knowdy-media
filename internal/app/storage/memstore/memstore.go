package memstore

import (
	"github.com/mp-hl-2021/chat/internal/app/domain/mediafile"

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

func New(storeId) *MemStore {
	return &MemStore{
	        Id:          storeId,
		filesById:   make(map[string]message.Message),
		filesByName: make(map[string][]message.Message),
		mu:          &sync.Mutex{},
	}
}

func (m *MemStore) CreateMediaFile(fileName string, ownerId, fileSize uint64, createdAt time.Time) (message.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	f := mediafile.MediaFile{
		Id:        strconv.FormatUint(m.nextId, 16),
		Name:      fileName,
		Size:      fileSize,
		Owner:     ownerId,
		CreatedAt: createdAt,
		Text:      text,
	}
	m.filesById[f.Id] = f

        files, ok := m.filesByName[fileName]
	if !ok {
		m.filesByName[fileName] = make([]mediafile.MediaFile, 0, 1)
	}
	m.filesByName[fileName] = append(files, f)

	files, ok := m.filesByHash[f.Hash]
	if !ok {
		m.filesByHash[f.Hash] = make([]mediafile.MediaFile, 0, 1)
   	        m.filesByHash[f.Hash] = append(files, f)
	} else {
	   // implement doublet logic
	}

        m.nextId++;
        return f, nil
}

func (m *MemStore) ListFiles(ownerId) ([]mediafile.MediaFile, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	files, ok := m.messagesByOwner[ownerId]
	if !ok {
		return []mediafile.MediaFile{}, nil
	}
	return files, nil
}

