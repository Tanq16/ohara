package store

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/tanishqrupaal/ohara/internal/model"
)

func (s *Store) loadTouchpoints() ([]model.Touchpoint, error) {
	data, err := os.ReadFile(s.touchpointsPath())
	if err != nil {
		return nil, err
	}
	var tps []model.Touchpoint
	if err := json.Unmarshal(data, &tps); err != nil {
		return nil, err
	}
	return tps, nil
}

func (s *Store) saveTouchpoints(tps []model.Touchpoint) error {
	data, err := json.MarshalIndent(tps, "", "  ")
	if err != nil {
		return err
	}
	return atomicWrite(s.touchpointsPath(), data)
}

func (s *Store) sanitizeInput(input *model.TouchpointInput) {
	input.Description = s.sanitizer.Sanitize(input.Description)
	input.URL = s.sanitizer.Sanitize(input.URL)
	if input.Tags == nil {
		input.Tags = []string{}
	}
	if input.PeopleInvolved == nil {
		input.PeopleInvolved = []string{}
	}
	for i, p := range input.PeopleInvolved {
		input.PeopleInvolved[i] = s.sanitizer.Sanitize(p)
	}
}

func (s *Store) validateInput(input model.TouchpointInput) error {
	if input.Description == "" {
		return validationErr("description is required")
	}

	md, err := s.loadMetadata()
	if err != nil {
		return fmt.Errorf("failed to load metadata for validation: %w", err)
	}

	if !contains(md.Categories, input.Category) {
		return validationErr(fmt.Sprintf("unknown category: %s", input.Category))
	}

	for _, tag := range input.Tags {
		if !contains(md.Tags, tag) {
			return validationErr(fmt.Sprintf("unknown tag: %s", tag))
		}
	}

	return nil
}

func (s *Store) ListTouchpoints(category, tag, startDate string) ([]model.Touchpoint, error) {
	s.tpMu.RLock()
	defer s.tpMu.RUnlock()

	tps, err := s.loadTouchpoints()
	if err != nil {
		return nil, err
	}

	var start time.Time
	if startDate != "" {
		start, err = time.Parse(time.RFC3339, startDate)
		if err != nil {
			return nil, validationErr(fmt.Sprintf("invalid start_date format: %s", startDate))
		}
	}

	result := make([]model.Touchpoint, 0, len(tps))
	for _, tp := range tps {
		if category != "" && tp.Category != category {
			continue
		}
		if tag != "" && !contains(tp.Tags, tag) {
			continue
		}
		if !start.IsZero() {
			t, err := time.Parse(time.RFC3339, tp.Date)
			if err != nil {
				continue
			}
			if t.Before(start) {
				continue
			}
		}
		result = append(result, tp)
	}

	return result, nil
}

func (s *Store) CreateTouchpoint(input model.TouchpointInput) (model.Touchpoint, error) {
	s.sanitizeInput(&input)

	s.tpMu.Lock()
	defer s.tpMu.Unlock()

	s.mdMu.RLock()
	err := s.validateInput(input)
	s.mdMu.RUnlock()
	if err != nil {
		return model.Touchpoint{}, err
	}

	tps, err := s.loadTouchpoints()
	if err != nil {
		return model.Touchpoint{}, err
	}

	tp := model.Touchpoint{
		ID:             uuid.New().String(),
		Date:           time.Now().UTC().Format(time.RFC3339),
		Description:    input.Description,
		Category:       input.Category,
		Tags:           input.Tags,
		PeopleInvolved: input.PeopleInvolved,
		URL:            input.URL,
	}

	tps = append(tps, tp)
	if err := s.saveTouchpoints(tps); err != nil {
		return model.Touchpoint{}, err
	}

	return tp, nil
}

func (s *Store) UpdateTouchpoint(id string, input model.TouchpointInput) (model.Touchpoint, error) {
	s.sanitizeInput(&input)

	s.tpMu.Lock()
	defer s.tpMu.Unlock()

	s.mdMu.RLock()
	err := s.validateInput(input)
	s.mdMu.RUnlock()
	if err != nil {
		return model.Touchpoint{}, err
	}

	tps, err := s.loadTouchpoints()
	if err != nil {
		return model.Touchpoint{}, err
	}

	for i, tp := range tps {
		if tp.ID == id {
			tps[i].Description = input.Description
			tps[i].Category = input.Category
			tps[i].Tags = input.Tags
			tps[i].PeopleInvolved = input.PeopleInvolved
			tps[i].URL = input.URL

			if err := s.saveTouchpoints(tps); err != nil {
				return model.Touchpoint{}, err
			}
			return tps[i], nil
		}
	}

	return model.Touchpoint{}, notFoundErr("touchpoint", id)
}

func (s *Store) DeleteTouchpoint(id string) error {
	s.tpMu.Lock()
	defer s.tpMu.Unlock()

	tps, err := s.loadTouchpoints()
	if err != nil {
		return err
	}

	filtered := make([]model.Touchpoint, 0, len(tps))
	found := false
	for _, tp := range tps {
		if tp.ID == id {
			found = true
			continue
		}
		filtered = append(filtered, tp)
	}

	if !found {
		return notFoundErr("touchpoint", id)
	}

	return s.saveTouchpoints(filtered)
}
