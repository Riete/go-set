package set

import (
	"context"
)

type threadUnsafeSet map[interface{}]struct{}

func newUnsafeSet() Set {
	s := newThreadUnsafeSet()
	return &s
}

func newThreadUnsafeSet() threadUnsafeSet {
	return make(threadUnsafeSet)
}

func (s *threadUnsafeSet) Add(items ...interface{}) {
	for _, item := range items {
		(*s)[item] = struct{}{}
	}
}

func (s *threadUnsafeSet) Remove(items ...interface{}) {
	for _, item := range items {
		delete(*s, item)
	}
}

func (s *threadUnsafeSet) Clear() {
	*s = newThreadUnsafeSet()
}

func (s *threadUnsafeSet) Has(items ...interface{}) bool {
	for _, item := range items {
		if _, found := (*s)[item]; !found {
			return false
		}
	}
	return true
}

func (s *threadUnsafeSet) Iter(ctx context.Context) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		for key := range *s {
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

func (s *threadUnsafeSet) Count() int {
	return len(*s)
}

func (s *threadUnsafeSet) SubsetOf(other Set) bool {
	if s.Count() > other.Count() {
		return false
	}
	for key := range *s {
		if !other.Has(key) {
			return false
		}
	}
	return true
}

func (s *threadUnsafeSet) HasSubset(other Set) bool {
	return other.SubsetOf(s)
}

func (s *threadUnsafeSet) Equal(other Set) bool {
	if s.Count() != other.Count() {
		return false
	}
	for key := range *s {
		if !other.Has(key) {
			return false
		}
	}
	return true
}

func (s *threadUnsafeSet) Copy() Set {
	set := newThreadUnsafeSet()
	for key := range *s {
		set.Add(key)
	}
	return &set
}

func (s *threadUnsafeSet) Union(other Set) Set {
	unionSet := newThreadUnsafeSet()
	o := other.(*threadUnsafeSet)

	for key := range *s {
		unionSet.Add(key)
	}
	for key := range *o {
		unionSet.Add(key)
	}
	return &unionSet
}

func (s *threadUnsafeSet) Intersect(other Set) Set {
	intersectSet := newThreadUnsafeSet()
	o := other.(*threadUnsafeSet)
	if s.Count() > o.Count() {
		for key := range *s {
			if o.Has(key) {
				intersectSet.Add(key)
			}
		}
	} else {
		for key := range *o {
			if s.Has(key) {
				intersectSet.Add(key)
			}
		}
	}
	return &intersectSet
}

func (s *threadUnsafeSet) Diff(other Set) Set {
	diffSet := newThreadUnsafeSet()
	for key := range *s {
		if !other.Has(key) {
			diffSet.Add(key)
		}
	}
	return &diffSet
}

func (s *threadUnsafeSet) MirrorDiff(other Set) Set {
	mirrorDiffSet := newThreadUnsafeSet()
	o := other.(*threadUnsafeSet)
	for key := range *s {
		if !o.Has(key) {
			mirrorDiffSet.Add(key)
		}
	}
	for key := range *o {
		if !s.Has(key) {
			mirrorDiffSet.Add(key)
		}
	}
	return &mirrorDiffSet
}

func (s *threadUnsafeSet) All() []interface{} {
	items := make([]interface{}, 0, s.Count())
	for key := range *s {
		items = append(items, key)
	}
	return items
}
