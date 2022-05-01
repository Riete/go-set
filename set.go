package go_set

import (
	"context"
)

type Set interface {
	Add(i ...interface{})
	Remove(i ...interface{})
	Clear()
	Count() int
	Has(i ...interface{}) bool
	Iter(context.Context) <-chan interface{}
	SubsetOf(Set) bool
	HasSubset(Set) bool
	Equal(Set) bool
	Copy() Set
	Union(Set) Set
	Intersect(Set) Set
	Diff(Set) Set
	MirrorDiff(Set) Set
	All() []interface{}
}

func NewSet(safe bool) Set {
	var s Set
	s = newUnsafeSet()
	if safe {
		s = newSafeSet()
	}
	return s
}

func NewSetWithValues(safe bool, items ...interface{}) Set {
	s := NewSet(safe)
	s.Add(items...)
	return s
}
