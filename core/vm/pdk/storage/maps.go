package storage

import "fmt"

type Map[K any, V any] struct {
	path *Path
}

func NewMap[K any, V any](p *Path) *Map[K, V] {
	return &Map[K, V]{path: p}
}

func (m *Map[K, V]) BindPath(p *Path) {
	m.path = p
}

func (m *Map[K, V]) Get(key K) (V, error) {
	var zero V
	if m.path == nil {
		return zero, fmt.Errorf("[K, V]: cannot get value from nil path")
	}
	var val V
	err := Load(m.path.Map(key), &val)
	return val, err
}

func (m *Map[K, V]) Set(key K, value V) error {
	if m.path == nil {
		return fmt.Errorf("[K, V]: cannot set value to nil path")
	}
	return Save(m.path.Map(key), value)
}
