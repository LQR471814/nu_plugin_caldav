package nutypes

import (
	"time"
)

type TimeSegment struct {
	Now          time.Time
	Duration     time.Duration
	ActiveEvents []EventReplica
}

type Timeline []TimeSegment
