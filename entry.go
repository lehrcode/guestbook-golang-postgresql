package main

import (
	"time"
)

type Entry struct {
	ID      int
	Name    string
	Email   string
	Message string
	Posted  time.Time
}
