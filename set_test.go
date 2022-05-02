package go_set

import (
	"context"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	s := NewSet()
	us := NewThreadUnsafeSet()
	s.Add(1, 2, 3)
	us.Add(1, 2, 3)
	t.Log(s.All())
	t.Log(us.All())
}

func TestRemove(t *testing.T) {
	s := NewSet(1, 2, 3)
	us := NewThreadUnsafeSet(1, 2, 4)
	s.Remove(1, 2)
	us.Remove(2, 4, 3)
	t.Log(s.All())
	t.Log(us.All())
}

func TestClear(t *testing.T) {
	s := NewSet(1, 2, 3)
	us := NewThreadUnsafeSet(1, 2, 4)
	s.Clear()
	us.Clear()
	t.Log(s.All())
	t.Log(us.All())
}

func TestHas(t *testing.T) {
	s := NewSet(1, 2, 3)
	us := NewThreadUnsafeSet(1, 2, 4)
	t.Log(s.Has(1, 2, 3), s.Has(2, 3, 4))
	t.Log(us.Has(1, 2, 4), us.Has(1, 2, 9))
}

func TestIter(t *testing.T) {
	ctx1, cancel1 := context.WithCancel(context.Background())
	ctx2, cancel2 := context.WithCancel(context.Background())

	s := NewSet(1, 2, 3)
	us := NewThreadUnsafeSet(1, 2, 3)
	go func() {
		time.Sleep(2 * time.Second)
		cancel1()
	}()
	for i := range s.Iter(ctx1) {
		time.Sleep(2 * time.Second)
		t.Log("safe set:", i)
	}
	go func() {
		time.Sleep(2 * time.Second)
		cancel2()
	}()
	for i := range us.Iter(ctx2) {
		time.Sleep(2 * time.Second)
		t.Log("unsafe set:", i)
	}
}

func TestCount(t *testing.T) {
	s := NewSet(1, 2, 3)
	us := NewThreadUnsafeSet(1, 2, 3)
	t.Log(s.Count(), us.Count())
}

func TestSubsetOf(t *testing.T) {
	s := NewSet(true, 1, 2, 3)
	o1 := NewSet(true, 1, 2, 3, 4)
	o2 := NewSet(true, 2, 3, 4)

	us := NewThreadUnsafeSet(false, 1, 2, 3)
	uo1 := NewThreadUnsafeSet(false, 1, 2, 3, 4)
	uo2 := NewThreadUnsafeSet(false, 2, 3, 4)
	t.Log(s.SubsetOf(o1), s.SubsetOf(o2))
	t.Log(us.SubsetOf(uo1), us.SubsetOf(uo2))
}

func TestHasSubset(t *testing.T) {
	s := NewSet(true, 1, 2, 3, 4, 5)
	o1 := NewSet(true, 1, 2, 3, 4)
	o2 := NewSet(true, 2, 3, 4, 6)

	us := NewThreadUnsafeSet(false, 1, 2, 3, 5, 4)
	uo1 := NewThreadUnsafeSet(false, 1, 2, 3, 4)
	uo2 := NewThreadUnsafeSet(false, 2, 3, 4, 6)
	t.Log(s.HasSubset(o1), s.HasSubset(o2))
	t.Log(us.HasSubset(uo1), us.HasSubset(uo2))
}

func TestEqual(t *testing.T) {
	s := NewSet(true, 1, 2, 3, 4, 5)
	s1 := NewSet(true, 1, 2, 3, 4, 5)
	s2 := NewSet(true, 1, 2, 3, 4)

	us := NewThreadUnsafeSet(false, 1, 2, 3, 5, 4)
	us1 := NewThreadUnsafeSet(false, 1, 2, 3, 5, 4)
	us2 := NewThreadUnsafeSet(false, 1, 2, 3, 5)
	t.Log(s.Equal(s1), s.Equal(s2))
	t.Log(us.Equal(us1), us.Equal(us2))
}

func TestUnion(t *testing.T) {
	s := NewSet(true, 1, 2)
	s1 := NewSet(true, "a", "b", 1, 2)
	us := NewThreadUnsafeSet(false, 1, 2)
	us1 := NewThreadUnsafeSet(false, "a", "b", 1, 2)
	t.Log(s.Union(s1).All())
	t.Log(us.Union(us1).All())
}

func TestIntersect(t *testing.T) {
	s := NewSet(true, 1, 2)
	s1 := NewSet(true, "a", "b", 1, 2)
	us := NewThreadUnsafeSet(false, 1, 2)
	us1 := NewThreadUnsafeSet(false, "a", "b", 1, 2)
	t.Log(s.Intersect(s1).All())
	t.Log(us.Intersect(us1).All())
}

func TestDiff(t *testing.T) {
	s := NewSet(true, 1, 2, 3)
	s1 := NewSet(true, "a", "b", 1, 2)
	us := NewThreadUnsafeSet(false, 1, 2, 3)
	us1 := NewThreadUnsafeSet(false, "a", "b", 1, 2)
	t.Log(s.Diff(s1).All(), s1.Diff(s).All())
	t.Log(us.Diff(us1).All(), us1.Diff(us).All())
}

func TestMirrorDiff(t *testing.T) {
	s := NewSet(true, 1, 2, 3)
	s1 := NewSet(true, "a", "b", 1, 2)
	us := NewThreadUnsafeSet(false, 1, 2, 3)
	us1 := NewThreadUnsafeSet(false, "a", "b", 1, 2)
	t.Log(s.MirrorDiff(s1).All())
	t.Log(us.MirrorDiff(us1).All())
}
