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

func NewSet(items ...interface{}) Set {
	s := newSafeSet()
	s.Add(items...)
	return s
}

func NewThreadUnsafeSet(items ...interface{}) Set {
	s := newUnsafeSet()
	s.Add(items...)
	return s
}
