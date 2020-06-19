package logparser

import "time"

type TimeRange struct {
	Start, End time.Time
}

func (t TimeRange)Includes(other time.Time) bool {
	return !other.Before(t.Start) && !other.After(t.End)
}


func (t TimeRange)Intersects(other TimeRange) bool {
	// (StartA <= EndB) and (EndA >= StartB)
	return !t.Start.After(other.End) && !t.End.Before(other.Start)
}