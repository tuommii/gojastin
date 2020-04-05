package server

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"text/template"
	"time"

	"miikka.xyz/gojastin/config"
)

// Responses for different cases
const (
	onEarly = "Fast enough"
	onLate  = "Too slow"
	onError = "Error"
)

type server struct {
	// Works as ID
	counter int
	// All requests gets unique key defined by counter
	visitors map[int]*visitor
	// Compilation timestamp gets injected here
	build string

	templ  *template.Template
	mu     sync.Mutex
	config *config.Config
	pool   *sync.Pool
}

// New returns a new server
func New(buildtime string) *server {
	s := &server{build: buildtime}

	s.visitors = make(map[int]*visitor)

	// We use this, when creating new visitor
	s.pool = &sync.Pool{
		New: func() interface{} {
			v := visitor{lastSeen: time.Now(), deadline: (time.Duration(rand.Intn(s.config.Deadline) + 1)) * time.Second}
			return &v
		},
	}

	s.config = config.New()
	templ, err := template.New("home").Parse(html)
	if err != nil {
		log.Fatal(err)
	}
	s.templ = templ
	return s
}

// Router handles all routes
func (s *server) Router(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		text(w, http.StatusMethodNotAllowed, onError)
		return
	}

	path := r.URL.String()
	switch path {
	case "/":
		s.startTimer()
		s.render(w)
		return
	case "/favicon.ico":
	case "/_status":
		text(w, http.StatusOK, "OK")
		return
	default:
		total, visitor := s.stopTimer(path[1:])
		status, msg := computeResponse(total, visitor)
		if status == http.StatusBadRequest {
			text(w, status, msg)
			return
		}
		delete(s.visitors, visitor.id)
		text(w, status, msg)
	}
}

// Implement http.Handler interface, for httptest purposes
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router(w, r)
}

// Helper funcs

func computeResponse(total time.Duration, v *visitor) (int, string) {
	status := http.StatusOK
	var msg string

	// This must be checked first
	if v == nil {
		status = http.StatusBadRequest
		msg = onError
		return status, msg
	}
	if total > v.deadline {
		msg = fmt.Sprintf("%s: %.4s\nTimelimit was: %.3s", onLate, total, v.deadline)
	} else {
		msg = fmt.Sprintf("%s: %.4s\nTimelimit was: %.3s", onEarly, total, v.deadline)
	}
	return status, msg
}

func text(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintln(w, msg)
}

func (s *server) render(w http.ResponseWriter) {
	// v := s.pool.Get().(*visitor)
	data := struct {
		Counter   int
		Deadline  float64
		BuildTime string
	}{
		s.counter,
		s.pool.Get().(*visitor).deadline.Seconds(),
		// s.visitors[s.counter].deadline.Seconds(),
		s.build,
	}
	s.templ.Execute(w, data)
}

// For sake of simplicity. Force reload on back-button
const html = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Timer</title>
	<style>
		body { margin-top: 1rem; }
		h1, h3, p a { font-family: -apple-system, BlinkMacSystemFont, Ubuntu, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;}
		h1, h3, { color: #444; }
		h1 { font-weight: 900; }
		h1 span { font-size: 1rem; font-weight: normal; }
		a, a:visited { color: #0366d6; }
		.compiled { font-size: 12px; margin-top: 0.5rem}
		.link { font-size: 2rem; }
	</style>
	<script>
		if(performance.navigation.type == 2) {
			location.reload(true);
		}
	</script>
</head>
<body>
	<center>
		<h1><span>TIME TO CLICK:&nbsp;&nbsp;&nbsp;</span> {{.Deadline}}<span> sec</span></h1>
		<br />
		<a class="link" href="/{{.Counter}}">LINK</a>
		<br />
		<br />
		<br />
		<br />
		<p>Created by: <a href="https://github.com/tuommii">Miikka Tuominen</a></p>
		<p class="compiled"><a href="https://github.com/tuommii/gojastin">Repository</a> on Github. Server compiled: {{.BuildTime}}</p>
	</center>
</body>
</html>
`
