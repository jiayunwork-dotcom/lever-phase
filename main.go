package main

import (
	"embed"
	"flag"
	"log"
	"net/http"

	"lever-phase/internal/web"
)

//go:embed web
var webFS embed.FS

//go:embed example/alloy-60.json
var alloy60JSON []byte

func main() {
	httpAddr := flag.String("http", ":8080", "serve the lever-phase web console on this address (e.g. :8080)")
	flag.Parse()
	handler, err := web.NewServer(web.Assets{
		WebFS: webFS,
		Examples: map[string][]byte{
			"alloy-60": alloy60JSON,
		},
	})
	if err != nil {
		log.Fatalf("build lever-phase server: %v", err)
	}
	log.Printf("lever-phase web console on http://localhost%s", *httpAddr)
	if err := http.ListenAndServe(*httpAddr, handler); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}
