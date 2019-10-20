package main

import "time"

type Kiwi struct {
}

func (m Kiwi) GetName() string {
	return "momondo"
}

func (m Kiwi) GetLowest(start, end time.Time) float32 {
	return 1.0
}
