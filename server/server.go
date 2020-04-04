package server

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"text/template"

	"golang.org/x/time/rate"
	"miikka.xyz/gojastin/config"
)

// Responses for different cases
const (
	onEarly = "Fast enough"
	onLate  = "Too slow"
	onError = "Error"
)

type server struct {
	// Works as ID/IP
	counter int
	// All requests gets unique key defined by counter
	visitors map[int]*visitor
	// Compilation timestamp gets injected here
	build string

	templ   *template.Template
	mu      sync.Mutex
	limiter *rate.Limiter
	config  *config.Config
}

// New return's new server
func New(buildtime string) *server {
	s := &server{build: buildtime}
	c := config.New()
	s.config = c
	s.visitors = make(map[int]*visitor)
	s.limiter = rate.NewLimiter(1, 100)
	templ, err := template.New("home").Parse(html)
	if err != nil {
		log.Fatal(err)
	}
	s.templ = templ
	return s
}

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
	case "_status":
		text(w, http.StatusOK, "OK")
		return
	default:
		measure, timeLimit, err := s.stopTimer(path[1:])
		if err != nil {
			if s.config.Logging {
				log.Println(err)
			}
			text(w, http.StatusBadRequest, onError)
			return
		}
		if measure > timeLimit {
			text(w, http.StatusOK, fmt.Sprintf("%s: %.4s\nTimelimit was: %.3s", onLate, measure, timeLimit))
			return
		}
		text(w, http.StatusOK, fmt.Sprintf("%s: %.4s\nTimelimit was: %.3s", onEarly, measure, timeLimit))
	}
}

// Implement http.Handler interface, for httptest purposes
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router(w, r)
}

// Global limiter
func (s *server) Limit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.limiter.Allow() == false {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func text(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintln(w, msg)
}

func (s *server) render(w http.ResponseWriter) {
	data := struct {
		Counter   int
		Deadline  float64
		BuildTime string
	}{
		s.counter,
		s.visitors[s.counter].deadline.Seconds(),
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
