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

// Server ...
type Server struct {
	// Works as ID
	counter int
	// All requests gets unique key defined by counter
	visitors map[int]*visitor
	// Compilation timestamp gets injected here
	build string

	pool   *sync.Pool
	mu     sync.Mutex
	templ  *template.Template
	config *config.Config
}

// New returns a new server
func New(buildtime string) *Server {
	s := &Server{build: buildtime}

	s.visitors = make(map[int]*visitor)

	// We use this, when creating new visitor
	s.pool = &sync.Pool{
		New: func() interface{} {
			seed := rand.Intn(s.config.Deadline) + 1
			deadline := time.Duration(seed) * time.Second
			v := visitor{lastSeen: time.Now(), deadline: deadline}
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
func (s *Server) Router(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		text(w, http.StatusMethodNotAllowed, "Error")
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
		timeSpent, visitor := s.stopTimer(path[1:])
		status, msg := computeResponse(timeSpent, visitor)
		if status == http.StatusBadRequest {
			text(w, status, msg)
			return
		}
		delete(s.visitors, visitor.id)
		text(w, status, msg)
	}
}

// Implement http.Handler interface, for httptest purposes
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router(w, r)
}

// Helper funcs

func computeResponse(timeSpent time.Duration, v *visitor) (int, string) {
	status := http.StatusOK
	var msg string

	// This must be checked first
	if v == nil {
		status = http.StatusBadRequest
		msg = "Error"
		return status, msg
	}

	info := fmt.Sprintf("Timelimit was: %.3s", v.deadline)
	if timeSpent > v.deadline {
		msg = fmt.Sprintf("Too slow: %.4s\n", timeSpent)
	} else {
		msg = fmt.Sprintf("Fast enough: %.4s\n", timeSpent)
	}
	msg += info
	return status, msg
}

func text(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintln(w, msg)
}

func (s *Server) render(w http.ResponseWriter) {
	data := struct {
		Counter   int
		Deadline  float64
		BuildTime string
	}{
		s.counter,
		s.pool.Get().(*visitor).deadline.Seconds(),
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
