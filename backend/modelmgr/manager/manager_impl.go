package manager

import (
	"context"
	"strings"
	"sync"
)

// staticManager 是静态模型管理器的实现。
type staticManager struct {
	models  []*InternalModel
	mapping map[int64]*InternalModel
	mu      sync.RWMutex
}

// NewStaticManager 创建一个静态模型管理器。
func NewStaticManager(models []*InternalModel) Manager {
	mapping := make(map[int64]*InternalModel, len(models))
	for _, m := range models {
		mapping[m.ID] = m
	}
	return &staticManager{
		models:  models,
		mapping: mapping,
		mu:      sync.RWMutex{},
	}
}

// List 列出符合条件的模型列表。
func (s *staticManager) List(ctx context.Context, opts *ListOptions) ([]*InternalModel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*InternalModel
	limit := opts.Limit
	if limit <= 0 {
		limit = 100
	}

	for _, m := range s.models {
		if len(result) >= limit {
			break
		}
		if len(opts.Status) > 0 {
			matched := false
			for _, status := range opts.Status {
				if m.Meta.Status == status {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}
		if opts.NameFilter != "" && !strings.Contains(m.Name, opts.NameFilter) {
			continue
		}
		result = append(result, m)
	}
	return result, nil
}

// Get 根据 ID 获取单个模型。
func (s *staticManager) Get(ctx context.Context, id int64) (*InternalModel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if m, ok := s.mapping[id]; ok {
		return m, nil
	}
	return nil, nil
}

// MGet 批量获取多个模型。
func (s *staticManager) MGet(ctx context.Context, ids []int64) ([]*InternalModel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*InternalModel, 0, len(ids))
	for _, id := range ids {
		if m, ok := s.mapping[id]; ok {
			result = append(result, m)
		}
	}
	return result, nil
}
