package dto

import (
	"time"
)

type TimeSegment struct {
	Now          time.Time
	Duration     time.Duration
	ActiveEvents []Event
}

type Timeline []TimeSegment
