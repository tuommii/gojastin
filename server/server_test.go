package server

import (
	"net/http"
	"net/http/httptest"
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

func TestHomeHandle(t *testing.T) {
	s := New("test")
	s.Log(false)
	server := httptest.NewServer(s)
	defer server.Close()

	// for i := range []int{100, -200} {
	for i := 0; i < 100; i++ {
		resp, err := http.Get(server.URL)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Wrong status: %d\n", resp.StatusCode)
		}

		// i starts from 0, my counter not
		if s.counter != i+1 {
			t.Errorf("Sended nth(%d) request, and counter was %d\n", i, s.counter)
		} else {
			t.Log("#", i, " request, and counter was ", s.counter)
		}
	}
}
