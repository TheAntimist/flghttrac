package main

import (
	"time"
)

type Source interface {
	GetName() string
	GetLowest(start, end time.Time) float32
}
