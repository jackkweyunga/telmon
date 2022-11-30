// Api for prometheus logs.

package web

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func pong(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
}

func Run(port int) {
	addr := flag.String("addr", fmt.Sprintf(":%v", port), "http service address")
	http.HandleFunc("/ping", pong)
	http.Handle(
		"/prometheus",
		promhttp.Handler(),
	)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
