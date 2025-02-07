package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Caller struct {
	IP        string            `json:"ip"`
	UserAgent string            `json:"user_agent"`
	Headers   map[string]string `json:"headers"`
	Timestamp string            `json:"timestamp"`
	URI       string            `json:"uri"`
}

var callers []Caller

func main() {
	http.HandleFunc("/image", imageHandler)
	http.HandleFunc("/callers", callersHandler)

	fmt.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// Log caller information
	log.Printf("Caller: %s\n", r.RemoteAddr)
	for key, value := range r.Header {
		log.Printf("%s: %v\n", key, value)
	}

	// Save caller information
	headers := make(map[string]string)
	for key, value := range r.Header {
		headers[key] = value[0]
	}

	caller := Caller{
		IP:        r.RemoteAddr,
		UserAgent: r.UserAgent(),
		Headers:   headers,
		Timestamp: time.Now().UTC().Format("2006-01-02 - 15:04:05"),
		URI:       r.URL.RawQuery,
	}
	callers = append(callers, caller)
	saveCallersToFile()

	// Write a 1px transparent PNG image
	pngData := []byte{
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
		0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x01, 0x03, 0x00, 0x00, 0x00, 0x25, 0xdb, 0x56,
		0xca, 0x00, 0x00, 0x00, 0x03, 0x50, 0x4c, 0x54,
		0x45, 0x00, 0x00, 0x00, 0xa7, 0x7a, 0x3d, 0xda,
		0x00, 0x00, 0x00, 0x01, 0x74, 0x52, 0x4e, 0x53,
		0x00, 0x40, 0xe6, 0xd8, 0x66, 0x00, 0x00, 0x00,
		0x0a, 0x49, 0x44, 0x41, 0x54, 0x08, 0xd7, 0x63,
		0x60, 0x00, 0x00, 0x00, 0x02, 0x00, 0x01, 0xe2,
		0x21, 0xbc, 0x33, 0x00, 0x00, 0x00, 0x00, 0x49,
		0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
	}

	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pngData)))
	w.Write(pngData)
}

func callersHandler(w http.ResponseWriter, r *http.Request) {
	jsonData, err := readCallersFromFile()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading caller data: %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonData)
}

func saveCallersToFile() {
	data, err := json.MarshalIndent(callers, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}
	err = os.WriteFile("callers.json", data, 0644)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}
}

func readCallersFromFile() ([]Caller, error) {
	data, err := os.ReadFile("callers.json")
	if err != nil {
		if os.IsNotExist(err) {
			return []Caller{}, nil
		}
		return nil, err
	}
	var callers []Caller
	err = json.Unmarshal(data, &callers)
	if err != nil {
		return nil, err
	}
	return callers, nil
}
