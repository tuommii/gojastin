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
	s.config.Reset = 10
	server := httptest.NewServer(s)
	defer server.Close()

	var count int

	// Check counter reset also
	for i := 0; i < s.config.Reset*3; i++ {
		resp, err := http.Get(server.URL)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Wrong status: %d\n", resp.StatusCode)
		}
		count++
		wanted := count % s.config.Reset

		if s.counter != wanted {
			t.Errorf("#%d request, counter was %d %d\n", count, wanted, s.counter)
		} else {
			t.Log("#", count, "req, want:", wanted, "got:", s.counter)
		}
	}
}
