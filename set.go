package set

import (
	"context"
	"sync"
)

type Set[T comparable] struct {
	s  map[T]struct{}
	rw sync.RWMutex
}

func New[T comparable](items ...T) *Set[T] {
	s := &Set[T]{s: make(map[T]struct{})}
	s.Add(items...)
	return s
}

func (s *Set[T]) Add(items ...T) {
	s.rw.Lock()
	defer s.rw.Unlock()
	for _, item := range items {
		s.s[item] = struct{}{}
	}
}

func (s *Set[T]) Remove(items ...T) {
	s.rw.Lock()
	defer s.rw.Unlock()
	for _, item := range items {
		delete(s.s, item)
	}
}

func (s *Set[T]) Clear() {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.s = make(map[T]struct{})
}

func (s *Set[T]) Has(items ...T) bool {
	s.rw.RLock()
	defer s.rw.RUnlock()
	for _, item := range items {
		if _, found := s.s[item]; !found {
			return false
		}
	}
	return true
}

func (s *Set[T]) Iter(ctx context.Context) <-chan T {
	ch := make(chan T)
	go func() {
		s.rw.RLock()
		defer s.rw.RUnlock()
		defer close(ch)
		for key := range s.s {
			select {
			case <-ctx.Done():
				return
			default:
				ch <- key
			}
		}
	}()
	return ch
}

func (s *Set[T]) Count() int {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return len(s.s)
}

func (s *Set[T]) SubsetOf(other *Set[T]) bool {
	s.rw.RLock()
	defer s.rw.RUnlock()
	if s.Count() > other.Count() {
		return false
	}
	for key := range s.s {
		if !other.Has(key) {
			return false
		}
	}
	return true
}

func (s *Set[T]) HasSubset(other *Set[T]) bool {
	return other.SubsetOf(s)
}

func (s *Set[T]) Equal(other *Set[T]) bool {
	s.rw.RLock()
	defer s.rw.RUnlock()
	if s.Count() != other.Count() {
		return false
	}
	for key := range s.s {
		if !other.Has(key) {
			return false
		}
	}
	return true
}

func (s *Set[T]) Copy() *Set[T] {
	s.rw.RLock()
	defer s.rw.RUnlock()
	set := New[T]()
	for key := range s.s {
		set.Add(key)
	}
	return set
}

func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	s.rw.RLock()
	defer s.rw.RUnlock()
	unionSet := New[T]()

	for key := range s.s {
		unionSet.Add(key)
	}
	for key := range other.s {
		unionSet.Add(key)
	}
	return unionSet
}

func (s *Set[T]) Intersect(other *Set[T]) *Set[T] {
	s.rw.RLock()
	defer s.rw.RUnlock()
	intersectSet := New[T]()
	if s.Count() > other.Count() {
		for key := range s.s {
			if other.Has(key) {
				intersectSet.Add(key)
			}
		}
	} else {
		for key := range other.s {
			if s.Has(key) {
				intersectSet.Add(key)
			}
		}
	}
	return intersectSet
}

func (s *Set[T]) Diff(other *Set[T]) *Set[T] {
	s.rw.RLock()
	defer s.rw.RUnlock()
	diffSet := New[T]()
	for key := range s.s {
		if !other.Has(key) {
			diffSet.Add(key)
		}
	}
	return diffSet
}

func (s *Set[T]) MirrorDiff(other *Set[T]) *Set[T] {
	s.rw.RLock()
	defer s.rw.RUnlock()
	mirrorDiffSet := New[T]()
	mirrorDiffSet.Add(s.Diff(other).All()...)
	mirrorDiffSet.Add(other.Diff(s).All()...)
	return mirrorDiffSet
}

func (s *Set[T]) All() []T {
	s.rw.RLock()
	defer s.rw.RUnlock()
	var items []T
	for key := range s.s {
		items = append(items, key)
	}
	return items
}
