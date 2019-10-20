package main

import "time"

type Skyscanner struct {
}

func (s Skyscanner) GetName() string {
	return "skyscanner"
}

func (s Skyscanner) GetLowest(start, end time.Time) float32 {
	return 1.0
}
