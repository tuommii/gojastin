package main

import (
	"log"
	"net/http"

	"miikka.xyz/gojastin/server"
)

// sudo vim /lib/systemd/system/timergo.service
// [Unit]
// Description=goweb

// [Service]
// Type=simple
// Restart=always
// RestartSec=5s
// ExecStart=/home/user/go/go-web/main

// [Install]
// WantedBy=multi-user.target

var buildTime = "NOT DYNAMIC"

func main() {
	s := server.New(buildTime)

	// Delete inactive visitors in background
	go s.CleanVisitors()

	http.HandleFunc("/", s.Limit(s.Router))
	log.Println("Started server on :3030...")
	if err := http.ListenAndServe("0.0.0.0:3030", nil); err != nil {
		log.Fatal(err)
	}
}
