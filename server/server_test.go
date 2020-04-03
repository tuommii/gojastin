package server

import (
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	res := New("test")
	if res.build != "test" {
		t.Errorf("Error while creating server")
	}
}

func BenchmarkPool(b *testing.B) {
	var p sync.Pool
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.Put(1)
			p.Get()
		}
	})
}

func BenchmarkAllocation(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i := 0
			i = i
		}
	})
}
