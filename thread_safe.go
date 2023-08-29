package set

import (
	"context"
	"sync"
)

type threadSafeSet struct {
	unsafe threadUnsafeSet
	rw     sync.RWMutex
}

func newSafeSet() Set {
	s := newThreadSafeSet()
	return &s
}

func newThreadSafeSet() threadSafeSet {
	return threadSafeSet{
		unsafe: newThreadUnsafeSet(),
		rw:     sync.RWMutex{},
	}
}

func (s *threadSafeSet) Add(items ...interface{}) {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.unsafe.Add(items...)
}

func (s *threadSafeSet) Remove(items ...interface{}) {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.unsafe.Remove(items...)
}

func (s *threadSafeSet) Clear() {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.unsafe = newThreadUnsafeSet()
}

func (s *threadSafeSet) Has(items ...interface{}) bool {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return s.unsafe.Has(items...)
}

func (s *threadSafeSet) Iter(ctx context.Context) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		s.rw.RLock()
		defer s.rw.RUnlock()
		defer close(ch)
		for key := range s.unsafe {
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

func (s *threadSafeSet) Count() int {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return s.unsafe.Count()
}

func (s *threadSafeSet) SubsetOf(other Set) bool {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return s.unsafe.SubsetOf(other)
}

func (s *threadSafeSet) HasSubset(other Set) bool {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return s.unsafe.HasSubset(other)
}

func (s *threadSafeSet) Equal(other Set) bool {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return s.unsafe.Equal(other)
}

func (s *threadSafeSet) Copy() Set {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return s.unsafe.Copy()
}

func (s *threadSafeSet) Union(other Set) Set {
	s.rw.RLock()
	defer s.rw.RUnlock()
	unionSet := newThreadSafeSet()
	o := other.(*threadSafeSet)
	for key := range s.unsafe {
		unionSet.unsafe.Add(key)
	}
	for key := range o.unsafe {
		unionSet.unsafe.Add(key)
	}
	return &unionSet
}

func (s *threadSafeSet) Intersect(other Set) Set {
	s.rw.RLock()
	defer s.rw.RUnlock()
	intersectSet := newThreadSafeSet()
	o := other.(*threadSafeSet)
	if s.unsafe.Count() > o.unsafe.Count() {
		for key := range s.unsafe {
			if o.unsafe.Has(key) {
				intersectSet.unsafe.Add(key)
			}
		}
	} else {
		for key := range o.unsafe {
			if s.unsafe.Has(key) {
				intersectSet.unsafe.Add(key)
			}
		}
	}
	return &intersectSet
}

func (s *threadSafeSet) Diff(other Set) Set {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return s.unsafe.Diff(other)
}

func (s *threadSafeSet) MirrorDiff(other Set) Set {
	s.rw.RLock()
	defer s.rw.RUnlock()
	mirrorDiffSet := newThreadSafeSet()
	o := other.(*threadSafeSet)
	for key := range s.unsafe {
		if !o.unsafe.Has(key) {
			mirrorDiffSet.unsafe.Add(key)
		}
	}
	for key := range o.unsafe {
		if !s.unsafe.Has(key) {
			mirrorDiffSet.unsafe.Add(key)
		}
	}
	return &mirrorDiffSet
}

func (s *threadSafeSet) All() []interface{} {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return s.unsafe.All()
}
