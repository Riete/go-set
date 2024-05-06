package set

import (
	"context"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	s := New(1, 2)
	s.Add(1, 2, 3)
	t.Log(s.All())
}

func TestRemove(t *testing.T) {
	s := New(1, 2, 3)
	s.Remove(1, 2)
	t.Log(s.All())
}

func TestClear(t *testing.T) {
	s := New(1, 2, 3)
	s.Clear()
	t.Log(s.All())
}

func TestHas(t *testing.T) {
	s := New(1, 2, 3)
	t.Log(s.Has(1, 2, 3), s.Has(2, 3, 4))
}

func TestIter(t *testing.T) {
	ctx1, cancel1 := context.WithCancel(context.Background())

	s := New(1, 2, 3)
	go func() {
		time.Sleep(2 * time.Second)
		cancel1()
	}()
	for i := range s.Iter(ctx1) {
		time.Sleep(2 * time.Second)
		t.Log("safe set:", i)
	}
}

func TestCount(t *testing.T) {
	s := New(1, 2, 3)
	t.Log(s.Count())
}

func TestSubsetOf(t *testing.T) {
	s := New(1, 2, 3)
	o1 := New(1, 2, 3, 4)
	o2 := New(2, 3, 4)

	t.Log(s.SubsetOf(o1), s.SubsetOf(o2))
}

func TestHasSubset(t *testing.T) {
	s := New(1, 2, 3, 4, 5)
	o1 := New(1, 2, 3, 4)
	o2 := New(2, 3, 4, 6)

	t.Log(s.HasSubset(o1), s.HasSubset(o2))
}

func TestEqual(t *testing.T) {
	s := New(1, 2, 3, 4, 5)
	s1 := New(1, 2, 3, 4, 5)
	s2 := New(1, 2, 3, 4)

	t.Log(s.Equal(s1), s.Equal(s2))
}

func TestCopy(t *testing.T) {
	s := New("a", "b", "c")
	c := s.Copy()
	s.Remove("a", "b")
	t.Log(s.All(), c.All())
}

func TestUnion(t *testing.T) {
	s := New("a", "b")
	s1 := New("c", "d")
	t.Log(s.Union(s1).All())
}

func TestIntersect(t *testing.T) {
	s := New("aa", "b")
	s1 := New("a", "b")
	t.Log(s.Intersect(s1).All())
}

func TestDiff(t *testing.T) {
	s := New(1, 2, 3)
	s1 := New(1, 2)
	t.Log(s.Diff(s1).All(), s1.Diff(s).All())
}

func TestMirrorDiff(t *testing.T) {
	s := New(1, 2, 3)
	s1 := New(1, 2, 5)
	t.Log(s.MirrorDiff(s1).All())
}
