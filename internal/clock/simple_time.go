package clock

import (
	"time"
)

// SimpleTime time.Time構造体からLocationを省いたシンプルなもの.
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

func (t *SimpleTime) GetTimeIndex(value int64, span SPAN) TimeIndex {
	return NewTimeIndex(t.Sec, value, span)
}
