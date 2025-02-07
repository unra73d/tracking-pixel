package main

import (
	"fmt"
	"log"
	"os"
	"tracking-pixel-api/util"

	"encoding/json"

	"github.com/gin-gonic/gin"
)

type Caller struct {
	IP        string              `json:"ip"`
	UserAgent string              `json:"user_agent"`
	Headers   map[string][]string `json:"headers"`
}

var callers []Caller
var l = util.Logger{}

func main() {
	r := gin.Default()

	// API to return a transparent 1px PNG image and log caller information.
	r.GET("/image", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "image/png")
		c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Writer.Header().Set("Pragma", "no-cache")
		c.Writer.Header().Set("Expires", "0")

		// Log caller information
		l.Debug(fmt.Sprintf("Caller: %s\n", c.ClientIP()))
		for key, value := range c.Request.Header {
			l.Debug(fmt.Sprintf("%s: %v\n", key, value))
		}

		// Save caller information
		headers := make(map[string][]string)
		for key, value := range c.Request.Header {
			headers[key] = value
		}
		caller := Caller{
			IP:        c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			Headers:   headers,
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

		// Set the correct Content-Length header
		c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(pngData)))

		// Write the image data to the response
		c.Writer.Write(pngData)
	})

	// API to list saved callers' information
	r.GET("/callers", func(c *gin.Context) {
		l.Debug("Listing saved callers...")
		jsonData, err := readCallersFromFile()
		if err != nil {
			c.String(500, "Error reading caller data: %s", err)
			return
		}
		c.JSON(200, jsonData)
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err) // Replace l.Fatal with log.Fatal
	}
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
			return []Caller{}, nil //Return empty array if file doesn't exist
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
