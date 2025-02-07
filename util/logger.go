package util

import (
	"log"
)

type Logger struct{}

func (l *Logger) Debug(message string) {
	log.Println("DEBUG:", message)
}

func (l *Logger) Info(message string) {
	log.Println("INFO:", message)
}

func (l *Logger) Warn(message string) {
	log.Println("WARN:", message)
}

func (l *Logger) Error(message string) {
	log.Println("ERROR:", message)
}
