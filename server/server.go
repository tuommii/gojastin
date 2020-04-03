package server

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"text/template"

	"golang.org/x/time/rate"
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
}

// New return's new server
func New(buildtime string) *server {
	s := &server{build: buildtime}
	s.visitors = make(map[int]*visitor)
	s.limiter = rate.NewLimiter(1, 100)
	s.templ, _ = template.New("home").Parse(html)
	return s
}

func (s *server) SetRateLimit(r rate.Limit, b int) {
	s.limiter = rate.NewLimiter(r, b)
}

func (s *server) Router(w http.ResponseWriter, r *http.Request) {
	path := r.URL.EscapedPath()
	if path == "/" {
		s.startTimer()
		s.serveTemplate(w)
	} else if path == "/favicon.ico" {
		// Dont count favicon
	} else {
		took, timeLimit, err := s.stopTimer(path[1:])
		if err != nil {
			log.Println(err)
			w.Write([]byte("Error"))
			return
		}
		if took > timeLimit {
			w.Write([]byte(fmt.Sprintf("Too slow: %s", took)))
			return
		}
		w.Write([]byte(fmt.Sprintf("Fast enough: %s", took)))
	}
}

// Middleware limiter
func (s *server) Limit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Global limiter
		if s.limiter.Allow() == false {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (s *server) serveTemplate(w http.ResponseWriter) {
	data := struct {
		Counter   int
		MaxTime   float64
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
		h1, h3, p a { color: #444; }
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
		<h1><span>TIME TO CLICK:</span> {{.MaxTime}}<span> sec</span></h1>
		<br />
		<a class="link" href="/{{.Counter}}">LINK</a>
		<br />
		<br />
		<br />
		<br />
		<p>Created by: <a href="https://github.com/tuommii">Miikka Tuominen</a></p>
		<p class="compiled">Server compiled: {{.BuildTime}}</p>
	</center>
</body>
</html>
`
