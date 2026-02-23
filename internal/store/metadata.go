package store

import (
	"encoding/json"
	"os"

	"github.com/tanq16/ohara/internal/model"
)

func (s *Store) loadMetadata() (model.Metadata, error) {
	data, err := os.ReadFile(s.metadataPath())
	if err != nil {
		return model.Metadata{}, err
	}
	var md model.Metadata
	if err := json.Unmarshal(data, &md); err != nil {
		return model.Metadata{}, err
	}
	return md, nil
}

func (s *Store) saveMetadata(md model.Metadata) error {
	data, err := json.MarshalIndent(md, "", "  ")
	if err != nil {
		return err
	}
	return atomicWrite(s.metadataPath(), data)
}

func (s *Store) GetMetadata() (model.Metadata, error) {
	s.mdMu.RLock()
	defer s.mdMu.RUnlock()
	return s.loadMetadata()
}

func (s *Store) AddCategory(name string) error {
	if name == "" {
		return validationErr("category name is required")
	}

	s.mdMu.Lock()
	defer s.mdMu.Unlock()

	md, err := s.loadMetadata()
	if err != nil {
		return err
	}

	if contains(md.Categories, name) {
		return alreadyExistsErr("category", name)
	}

	md.Categories = append(md.Categories, name)
	return s.saveMetadata(md)
}

func (s *Store) RemoveCategory(name string) error {
	if name == "" {
		return validationErr("category name is required")
	}

	s.mdMu.Lock()
	defer s.mdMu.Unlock()

	md, err := s.loadMetadata()
	if err != nil {
		return err
	}

	idx := -1
	for i, c := range md.Categories {
		if c == name {
			idx = i
			break
		}
	}
	if idx == -1 {
		return notFoundErr("category", name)
	}

	md.Categories = append(md.Categories[:idx], md.Categories[idx+1:]...)
	return s.saveMetadata(md)
}

func (s *Store) AddTag(name string) error {
	if name == "" {
		return validationErr("tag name is required")
	}

	s.mdMu.Lock()
	defer s.mdMu.Unlock()

	md, err := s.loadMetadata()
	if err != nil {
		return err
	}

	if contains(md.Tags, name) {
		return alreadyExistsErr("tag", name)
	}

	md.Tags = append(md.Tags, name)
	return s.saveMetadata(md)
}

func (s *Store) RemoveTag(name string) error {
	if name == "" {
		return validationErr("tag name is required")
	}

	s.mdMu.Lock()
	defer s.mdMu.Unlock()

	md, err := s.loadMetadata()
	if err != nil {
		return err
	}

	idx := -1
	for i, t := range md.Tags {
		if t == name {
			idx = i
			break
		}
	}
	if idx == -1 {
		return notFoundErr("tag", name)
	}

	md.Tags = append(md.Tags[:idx], md.Tags[idx+1:]...)
	return s.saveMetadata(md)
}
