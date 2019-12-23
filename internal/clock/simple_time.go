package clock

import (
	"time"
)

type SimpleTime struct {
	Sec  int64
	Nsec int32
}

func MakeSimpleTime(time time.Time) SimpleTime {
	return SimpleTime{
		Sec:  time.Unix(),
		Nsec: int32(time.Nanosecond()),
	}
}

func (t *SimpleTime) Time() time.Time {
	return time.Unix(t.Sec, int64(t.Nsec))
}

func (t *SimpleTime) MakeTimeIndex(value int64, span SPAN) TimeIndex {
	return CreateTimeIndex(t.Sec, value, span)
}
