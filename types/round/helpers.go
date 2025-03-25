package roundtypes

import (
	"time"
)

func DescriptionPtr(desc string) *Description {
	d := Description(desc)
	return &d
}

func LocationPtr(loc string) *Location {
	l := Location(loc)
	return &l
}

func StartTimePtr(t time.Time) *StartTime {
	st := StartTime(t)
	return &st
}
