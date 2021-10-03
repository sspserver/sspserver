package targetaccessor

import (
	"sort"

	"geniusrabbit.dev/sspserver/internal/models"
)

// SimpleTargetAccessorLoader function
type SimpleTargetAccessorLoader func() ([]models.Target, error)

// SimpleTargetAccessor wrapper
type SimpleTargetAccessor struct {
	loader SimpleTargetAccessorLoader
	list   []models.Target
}

// NewSimpleTargetAccessor object
func NewSimpleTargetAccessor(loader SimpleTargetAccessorLoader) *SimpleTargetAccessor {
	return &SimpleTargetAccessor{
		loader: loader,
	}
}

// Reload data list
func (ta *SimpleTargetAccessor) Reload() error {
	list, err := ta.loader()
	if err != nil {
		return err
	}
	sort.Slice(list, func(i, j int) bool { return list[i].ID() < list[j].ID() })
	ta.list = list
	return nil
}

// TargetByID returns target object
func (ta *SimpleTargetAccessor) TargetByID(id uint64) models.Target {
	var (
		list = ta.list
		idx  = sort.Search(len(list), func(i int) bool { return list[i].ID() >= id })
	)
	if idx >= 0 && idx < len(list) && list[idx].ID() == id {
		return list[idx]
	}
	return nil
}
