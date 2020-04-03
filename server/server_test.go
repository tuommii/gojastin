package server

import (
	"testing"
)

func TestNew(t *testing.T) {
	s := New("test")
	if s.build != "test" {
		t.Errorf("Error while creating server")
	}
}

func BenchmarkStartTimer(b *testing.B) {
	s := New("test")
	s.Log(false)
	for n := 0; n < b.N; n++ {
		s.startTimer()
	}
}
