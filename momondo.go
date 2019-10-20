package main

import "time"

type Momondo struct {
	Source
}

func (m Momondo) GetName() string {
	return "momondo"
}

func (m Momondo) GetLowest(start, end time.Time) float32 {
	return 1.0
}
