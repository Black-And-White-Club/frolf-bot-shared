package roundtypes

import (
	"time"

	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

func DescriptionPtr(desc string) *Description {
	d := Description(desc)
	return &d
}

func LocationPtr(loc string) *Location {
	l := Location(loc)
	return &l
}

func StartTimePtr(t time.Time) *sharedtypes.StartTime {
	st := sharedtypes.StartTime(t)
	return &st
}
