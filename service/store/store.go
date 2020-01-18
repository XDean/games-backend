package store

import (
	"fmt"
)

type (
	Service struct {
		Root *Entity
	}

	Entity struct {
		Name     string
		Parent   *Entity
		Children map[string]*Entity
		Value    interface{}
	}
)

func (s *Service) Create() {
}

func (s *Service) Delete() {

}

func (s *Service) Get(namespace ...string) (*Entity, bool) {
	e := s.Root
	for _, name := range namespace {
		if child, ok := e.Get(name); ok {
			e = child
		} else {
			return nil, false
		}
	}
	return e, true
}

func (e *Entity) Create(name string) (*Entity, error) {
	if _, ok := e.Get(name); ok {
		return nil, fmt.Errorf("Entity already exist %s under %s.", name, e.GetNameSpace())
	}
	child := &Entity{
		Name:     name,
		Parent:   e,
		Children: map[string]*Entity{},
		Value:    nil,
	}
	e.Children[name] = child
	return child, nil
}

func (e *Entity) Delete(name string) (*Entity, error) {
	if child, ok := e.Get(name); ok {
		return nil, fmt.Errorf("No such Entity %s under %s.", name, e.GetNameSpace())
	} else {
		delete(e.Children, name)
		return child, nil
	}
}

func (e *Entity) Get(name string) (*Entity, bool) {
	child, ok := e.Children[name]
	return child, ok
}

func (e *Entity) GetNameSpace() []string {
	res := make([]string, 0)
	en := e
	for en.Parent != nil {
		res = append(res, en.Name)
		en = en.Parent
	}
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return res
}
