// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Poll file for changes with this period.
	filePeriod = 10 * time.Second
)

var (
	homeTempl = template.Must(template.New("").Parse(homeHTML))
	filename  string
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func readFileIfModified(lastMod time.Time) (Stats, []byte, time.Time, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return Stats{}, nil, lastMod, err
	}
	if !fi.ModTime().After(lastMod) {
		return Stats{}, nil, lastMod, nil
	}
	fileScanner, err := os.Open(filename)

	if err != nil {
		log.Fatalf("failed to open")
	}
	scanner := bufio.NewScanner(fileScanner)

	stats := processData(scanner)

	fileScanner.Close()

	p, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return stats, p, fi.ModTime(), nil
}

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func writer(ws *websocket.Conn, lastMod time.Time) {
	lastError := ""
	pingTicker := time.NewTicker(pingPeriod)
	fileTicker := time.NewTicker(filePeriod)
	defer func() {
		pingTicker.Stop()
		fileTicker.Stop()
		ws.Close()
	}()
	for {
		select {
		case <-fileTicker.C:
			var p []byte
			var err error
			var stats Stats

			stats, p, lastMod, err = readFileIfModified(lastMod)

			if err != nil {
				if s := err.Error(); s != lastError {
					lastError = s
					p = []byte(lastError)
				}
			} else {
				lastError = ""
			}

			if p != nil {

				var content = struct {
					Logs       []byte `json:"logs"`
					Statistics Stats  `json:"statistics"`
				}{
					Logs:       p,
					Statistics: stats,
				}

				contentBytes, err := json.Marshal(&content)
				if err != nil {
					log.Fatal(err)
				}

				ws.SetWriteDeadline(time.Now().Add(writeWait))
				if err := ws.WriteMessage(websocket.TextMessage, contentBytes); err != nil {
					return
				}
			}

		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func pong(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	var lastMod time.Time
	if n, err := strconv.ParseInt(r.FormValue("lastMod"), 16, 64); err == nil {
		lastMod = time.Unix(0, n)
	}

	go writer(ws, lastMod)
	reader(ws)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	stats, p, lastMod, err := readFileIfModified(time.Time{})
	if err != nil {
		p = []byte(err.Error())
		lastMod = time.Unix(0, 0)
	}

	var v = struct {
		Host    string
		Data    string
		LastMod string
		Stats   Stats
	}{
		r.Host,
		string(p),
		strconv.FormatInt(lastMod.UnixNano(), 16),
		stats,
	}
	homeTempl.Execute(w, &v)
}

func Run(port int) {
	addr := flag.String("addr", fmt.Sprintf(":%v", port), "http service address")
	filename = "telmon.log"
	http.HandleFunc("/ping", pong)
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWs)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}

const homeHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <title>Telmon</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="w-screen relative overflow-x-hidden h-fit">
<nav class="nav sticky top-0 bg-white border-b flex p-4 items-center justify-between">
    <h1 class="font-bold px-4 text-2xl">telmon</h1>
    <h5 class="px-5 text-sm font-medium">STATISTICS</h5>
    <div class="p-2 space-x-3 text-sm flex">
        <div>status: <small class="text-gray-900" id="status">connected</small></div>
        <div>lines: <small class="text-gray-900" id="count">{{.Stats.Count}}</small></div>
    </div>
</nav>

<div style="padding-top: 30px; padding-bottom: 40px;"
     class="flex-col justify-center items-center h-full w-full overflow-hidden space-y-6 space-x-6 ">

    <canvas class="p-5 flex justify-center" id="statsChart" width="400" height="250"></canvas>
    <h5 class="border-b text-sm font-medium">LOGS</h5>
    <div class="w-full h-96 p-2 bg-gray-900 w-100 flex overflow-auto dark">
        <pre class="text-white w-100" id="fileData">{{.Data}}</pre>
    </div>
</div>

<!--chats.json-->
<script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
<script src="https://cdn.jsdelivr.net/npm/chartjs-adapter-date-fns/dist/chartjs-adapter-date-fns.bundle.min.js"></script>
<script type="text/javascript">

    const ctx = document.getElementById('statsChart');

    let chart = new Chart(ctx, {
        type: 'line',
        data: {
            datasets: []
        },
        options: {
            scales: {
                y: {
                    min: 0
                },
                x: {
                    type: 'time',
                    time: {
                        displayFormats: {
                            quarter: 'HHH DDD MMM'
                        }
                    }
                }
            }
        }
    });

    function genGraph(datasets) {
        chart.data.datasets = datasets
        chart.update()
    }

    (function () {
        const status = document.getElementById("status");
        const logs = document.getElementById("fileData");
        const count = document.getElementById("count");
        const conn = new WebSocket("ws://{{.Host}}/ws?lastMod={{.LastMod}}");
        conn.onclose = function (evt) {
            status.textContent = 'Connection closed';
        }
        conn.onmessage = function (evt) {
            console.log('file updated');
            const json = JSON.parse(evt.data);
            const stats = json["statistics"]
            logs.textContent = atob(json["logs"]);
            count.textContent = stats["count"];

            let datasets = []
            let keySet = new Set(stats["data"].map(line => line["kind"]))
            let keys = new Array(...keySet)

            keys.forEach(key => {
                let sortedData = new Array(...stats["data"].filter(x => x["kind"] === key.toString()))
                datasets.push(
                    {
                        label: key,
                        data: sortedData.map(line => {
                            return {
                                x: line["time"],
                                y: key.toString() === "restored" ? 2 : key.toString() === "failed" ? 1 : 3,
                            }
                        }),
                        fill: false,
                        tension: 0.1,
                        borderColor: key.toString() === "restored" ? 'rgb(255,255,0)' : key.toString() === "failed" ? 'rgb(192, 0, 0)' : 'rgb(0, 192, 0)',
                    }
                )
            })

            genGraph(datasets)
            console.log(datasets)

        }
    })();
</script>
</body>
</html>
`
