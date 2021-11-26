// Penggunaan Sederhana SSE (Server Sent Events)
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// struktur data yang dikirimkan ke klien
type Client struct {
	name   string
	events chan *DashBoard
}

type DashBoard struct {
	User uint
}

func main() {
	// Membuat layanan web server pada port 3000

	// Routing & Callback function
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Service Up !!!")
	})
	http.HandleFunc("/sse", dashboardHandle)

	// Listen address
	address := fmt.Sprintf("%s:%d", "localhost", 3000)

	log.Printf("Listen %s", address)
	http.ListenAndServe(address, nil) // start web server
}

func dashboardHandle(w http.ResponseWriter, r *http.Request) {
	// Instan variabel client
	client := &Client{name: r.RemoteAddr, events: make(chan *DashBoard, 10)}

	// Proses membuat nomor acak dan mengirimkan data melalui chanel ke goroutine
	go updateDashBoard(client)

	// Seting header response
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// set timeout
	timeout := time.After(1 * time.Second)

	select {
	case ev := <-client.events:
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.Encode(ev)

		fmt.Fprintf(w, "data: %v\n\n", buf.String()) // respon hasil pengacakan nomor
		fmt.Printf("data: %v\n", buf.String())

	case <-timeout:
		fmt.Fprintf(w, "nothing to sent\n\n")

	}

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

func updateDashBoard(client *Client) {
	// perulangan tak terhingga untuk membuat nomor acak
	// dan mengirimkan melalui chanel events

	for {
		db := &DashBoard{User: uint(rand.Uint32())}

		client.events <- db
	}
}
