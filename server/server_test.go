package server

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func BenchmarkAllocWithoutPool(b *testing.B) {
	s := New("test")
	s.config.Logging = false
	for n := 0; n < b.N; n++ {
		v := &visitor{deadline: (time.Duration(rand.Intn(s.config.Deadline) + 1)) * time.Second, lastSeen: time.Now()}
		v.id++
		s.visitors[n] = v
	}
}
func BenchmarkAllocWithPool(b *testing.B) {
	s := New("test")
	s.config.Logging = false
	for n := 0; n < b.N; n++ {
		v := s.pool.Get().(*visitor)
		v.id++
		s.pool.Put(v)
		s.visitors[n] = v
	}
}

func BenchmarkStart(b *testing.B) {
	s := New("test")
	s.config.Logging = false
	for n := 0; n < b.N; n++ {
		s.startTimer()
	}
}
func BenchmarkStartAndStop(b *testing.B) {
	s := New("test")
	s.config.Logging = false
	for n := 0; n < b.N; n++ {
		s.startTimer()
		s.stopTimer(strconv.Itoa(n))
	}
}

func BenchmarkStartAndHalfStopped(b *testing.B) {
	s := New("test")
	s.config.Logging = false
	for n := 0; n < b.N; n++ {
		s.startTimer()
		if n%2 == 0 {
			s.stopTimer(strconv.Itoa(n))
		}
	}
}

func TestTimerStop(t *testing.T) {
	s := New("test")
	s.config.Logging = false
	s.config.MaxVisitors = 10
	s.config.VisitorAlive = 5 * time.Second
	s.config.RemoveInterval = 5 * time.Second
	s.config.Deadline = 1
	server := httptest.NewServer(s)
	defer server.Close()

	// Get ID
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Wrong status: %d\n", resp.StatusCode)
	}
	defer resp.Body.Close()

	// Check with ID
	resp, err = http.Get(server.URL + "/" + strconv.Itoa(s.counter))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Wrong status: %d\n", resp.StatusCode)
	}

	// Check with same id again, shoud be deleted already
	resp, err = http.Get(server.URL + "/" + strconv.Itoa(s.counter))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Wrong status: %d\n", resp.StatusCode)
	}

	// Check with ID that aint created yet
	resp, err = http.Get(server.URL + "/" + strconv.Itoa(s.counter+1))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Wrong status: %d\n", resp.StatusCode)
	}
}

// Loop over counter max value and check if counter resets
// afterwards len(visitors) should be same as max value
// Then testing
func TestHomeHandle(t *testing.T) {
	s := New("test")
	s.config.Logging = false
	s.config.MaxVisitors = 10
	s.config.VisitorAlive = 2 * time.Second
	s.config.RemoveInterval = 2 * time.Second
	server := httptest.NewServer(s)
	defer server.Close()

	for i := 0; i < s.config.MaxVisitors*3; i++ {
		resp, err := http.Get(server.URL)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Wrong status: %d\n", resp.StatusCode)
		}

		wanted := (i + 1) % s.config.MaxVisitors
		if s.counter != wanted {
			t.Errorf("#%d request, counter was %d %d\n", i+1, wanted, s.counter)
		}
	}

	if s.config.MaxVisitors != len(s.visitors) {
		t.Errorf("#%d rounds, len(visitors): %d\n", s.config.MaxVisitors, len(s.visitors))
	}

	go s.CleanVisitors()

	// Just to be sure
	time.Sleep(s.config.VisitorAlive + (time.Second * 1))

	// Should be cleaned now
	if len(s.visitors) != 0 {
		t.Errorf("#%d rounds, len(visitors): %d\n", s.config.MaxVisitors, len(s.visitors))
	}
}
