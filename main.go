package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const maxLines = 100

// CircularBuffer to hold the last n log lines
type CircularBuffer struct {
	lines []string
	index int
	mutex sync.Mutex
}

func NewCircularBuffer(size int) *CircularBuffer {
	return &CircularBuffer{
		lines: make([]string, size),
		index: 0,
	}
}

func (cb *CircularBuffer) Add(line string) {
	cb.mutex.Lock()
	cb.lines[cb.index] = line
	cb.index = (cb.index + 1) % maxLines
	cb.mutex.Unlock()
}

func (cb *CircularBuffer) GetAll() []string {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	result := make([]string, 0, maxLines)
	for i := 0; i < maxLines; i++ {
		line := cb.lines[(cb.index+i)%maxLines]
		if line != "" {
			result = append(result, line)
		}
	}
	return result
}

func (cb *CircularBuffer) Search(term string) []string {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	results := []string{}
	lowerTerm := strings.ToLower(term)
	for i := 0; i < maxLines; i++ {
		line := cb.lines[(cb.index+i)%maxLines]
		if line != "" && strings.Contains(strings.ToLower(line), lowerTerm) {
			results = append(results, line)
		}
	}
	return results
}

// Handler to serve the latest logs via HTTP
func logHandler(logBuffer *CircularBuffer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		term := r.URL.Query().Get("search")
		var logs []string

		if term != "" {
			logs = logBuffer.Search(term)
		} else {
			logs = logBuffer.GetAll()
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, "<h1>Logs</h1>")
		fmt.Fprintln(w, "<form method='GET'>")
		fmt.Fprintln(w, "<input type='text' name='search' placeholder='Search logs' value='"+term+"'/>")
		fmt.Fprintln(w, "<input type='submit' value='Search'/>")
		fmt.Fprintln(w, "<input type='button' value='Refresh' onclick='window.location.reload();'/>")
		fmt.Fprintln(w, "</form>")
		fmt.Fprintln(w, "<hr>")
		for _, logLine := range logs {
			// Append <br> for line breaks in HTML
			fmt.Fprintf(w, "%s<br>", logLine)
		}
	}
}

func main() {
	// Create or open the log file
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer logFile.Close()

	// Set log output to the log file
	log.SetOutput(logFile)

	// Create a CircularBuffer for the last 100 log entries
	logBuffer := NewCircularBuffer(maxLines)

	// Log the start time
	startLog := fmt.Sprintf("UDP Logger started at: %s", time.Now().Format(time.RFC3339))
	log.Println(startLog)
	logBuffer.Add(startLog)

	// Start the HTTP server for the web interface
	go func() {
		http.HandleFunc("/logs", logHandler(logBuffer))
		log.Println("Starting HTTP server on :8080...")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("error starting HTTP server: %v", err)
		}
	}()

	// Listen for UDP messages
	addr, err := net.ResolveUDPAddr("udp", ":514")
	if err != nil {
		log.Fatalf("error resolving UDP address: %v", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("error listening on UDP port 514: %v", err)
	}
	defer conn.Close()

	log.Println("Listening on UDP port 514...")

	// Buffer for incoming messages
	buffer := make([]byte, 1024)

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("error reading from UDP connection: %v", err)
			continue
		}

		// Log the received message along with the sender's address and timestamp
		logLine := fmt.Sprintf("[%s] Received message from %s: %s", time.Now().Format(time.RFC3339), remoteAddr.String(), string(buffer[:n]))
		log.Println(logLine)
		logBuffer.Add(logLine)
	}
}
