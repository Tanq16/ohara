package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/microcosm-cc/bluemonday"
	"github.com/tanq16/ohara/internal/model"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrAlreadyExists   = errors.New("already exists")
	ErrValidation      = errors.New("validation error")
	ErrInvalidFilename = errors.New("invalid filename")
)

func notFoundErr(kind, id string) error {
	return fmt.Errorf("%s %s: %w", kind, id, ErrNotFound)
}

func alreadyExistsErr(kind, name string) error {
	return fmt.Errorf("%s %s: %w", kind, name, ErrAlreadyExists)
}

func validationErr(msg string) error {
	return fmt.Errorf("%s: %w", msg, ErrValidation)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

type Config struct {
	DataDir string
}

type Store struct {
	dataDir   string
	tpMu      sync.RWMutex
	mdMu      sync.RWMutex
	sanitizer *bluemonday.Policy
}

func New(cfg Config) (*Store, error) {
	s := &Store{
		dataDir:   cfg.DataDir,
		sanitizer: bluemonday.StrictPolicy(),
	}

	if err := os.MkdirAll(filepath.Join(cfg.DataDir, "reports"), 0755); err != nil {
		return nil, err
	}

	tpPath := filepath.Join(cfg.DataDir, "touchpoints.json")
	if _, err := os.Stat(tpPath); os.IsNotExist(err) {
		if err := os.WriteFile(tpPath, []byte("[]"), 0644); err != nil {
			return nil, err
		}
	}

	mdPath := filepath.Join(cfg.DataDir, "metadata.json")
	if _, err := os.Stat(mdPath); os.IsNotExist(err) {
		defaults := model.Metadata{
			Categories: []string{
				"Feature Development",
				"Bug Fix",
				"Code Review",
				"Technical Design",
				"Tooling & Infrastructure",
				"Documentation",
				"Incident Response",
				"Mentorship",
				"Cross-team Collaboration",
				"Knowledge Sharing",
			},
			Tags: []string{"backend", "frontend", "devops", "testing", "security", "performance"},
		}
		data, err := json.MarshalIndent(defaults, "", "  ")
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(mdPath, data, 0644); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Store) touchpointsPath() string {
	return filepath.Join(s.dataDir, "touchpoints.json")
}

func (s *Store) metadataPath() string {
	return filepath.Join(s.dataDir, "metadata.json")
}

func (s *Store) reportsDir() string {
	return filepath.Join(s.dataDir, "reports")
}

func atomicWrite(path string, data []byte) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
