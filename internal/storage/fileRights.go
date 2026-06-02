package storage

import "github.com/hurtki/fed/internal/domain"

type fileRightsEntry struct {
	Read bool
	Edit bool
}

// Not concurrent safe implementation of storage for rights of agent
// to read/write project files
type MemoryFileRightsStorage struct {
	m map[domain.ProjectFile]fileRightsEntry
}

func NewMemoryFileRightsStorage() *MemoryFileRightsStorage {
	return &MemoryFileRightsStorage{
		m: make(map[domain.ProjectFile]fileRightsEntry),
	}
}

func (m *MemoryFileRightsStorage) EligibleForEdit(pf domain.ProjectFile) (bool, error) {
	e := m.m[pf]
	return e.Edit, nil
}

func (m *MemoryFileRightsStorage) EligibleForRead(pf domain.ProjectFile) (bool, error) {
	e := m.m[pf]
	return e.Read, nil
}

func (m *MemoryFileRightsStorage) SetEligibleForEdit(pf domain.ProjectFile) error {
	e := m.m[pf]
	e.Edit = true
	m.m[pf] = e
	return nil
}

func (m *MemoryFileRightsStorage) SetEligibleForRead(pf domain.ProjectFile) error {
	e := m.m[pf]
	e.Read = true
	m.m[pf] = e
	return nil
}
