package store

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

var validFilename = regexp.MustCompile(`^[\w][\w\-]*\.md$`)

func (s *Store) validateReportFilename(filename string) error {
	if !validFilename.MatchString(filename) {
		return fmt.Errorf("report filename %s (must match alphanumeric/hyphens ending in .md): %w", filename, ErrInvalidFilename)
	}
	return nil
}

func (s *Store) ListReports() ([]string, error) {
	entries, err := os.ReadDir(s.reportsDir())
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".md" {
			names = append(names, e.Name())
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(names)))
	return names, nil
}

func (s *Store) GetReport(filename string) (string, error) {
	if err := s.validateReportFilename(filename); err != nil {
		return "", err
	}

	data, err := os.ReadFile(filepath.Join(s.reportsDir(), filename))
	if err != nil {
		if os.IsNotExist(err) {
			return "", notFoundErr("report", filename)
		}
		return "", err
	}
	return string(data), nil
}

func (s *Store) CreateReport(filename, content string) error {
	if err := s.validateReportFilename(filename); err != nil {
		return err
	}

	path := filepath.Join(s.reportsDir(), filename)
	if _, err := os.Stat(path); err == nil {
		return alreadyExistsErr("report", filename)
	}

	return os.WriteFile(path, []byte(content), 0644)
}
