package clock

import "time"

type SPAN int8

const (
	SPAN_NANO        SPAN = 0
	SPAN_MICROSECOND SPAN = 1
	SPAN_MILLISECOND SPAN = 2
	SPAN_SECOND      SPAN = 3
	SPAN_MINUTE      SPAN = 4
	SPAN_HOUR        SPAN = 5
)

// TimeIndex は January 1, 1970 UTC からのTimeSpan間隔での番号(順番)を表す
type TimeIndex struct {
	TimeSpan TimeSpan
	Index    int64 // 番号(順番)
}

func (t *TimeIndex) Equal(t2 TimeIndex) bool {
	if t.TimeSpan.Value == t2.TimeSpan.Value && t.TimeSpan.Span == t2.TimeSpan.Span && t.Index == t2.Index {
		return true
	}
	return false
}

func SpanToDuration(span SPAN) time.Duration {
	if span == SPAN_NANO {
		return time.Nanosecond
	} else if span == SPAN_MICROSECOND {
		return time.Microsecond
	} else if span == SPAN_MILLISECOND {
		return time.Millisecond
	} else if span == SPAN_SECOND {
		return time.Second
	} else if span == SPAN_MINUTE {
		return time.Minute
	} else if span == SPAN_HOUR {
		return time.Hour
	}
	panic("Span2Duration : Unknown span")
}

// CreateTimeIndex は TimeIndex を返す.
func CreateTimeIndex(sec int64, value int64, span SPAN) TimeIndex {
	// 引数secが秒のためそれより小さいspanは対応しない
	if span == SPAN_NANO || span == SPAN_MICROSECOND || span == SPAN_MILLISECOND {
		panic("Unsupported span")
	}
	return TimeIndex{
		TimeSpan: TimeSpan{Value: value, Span: span},
		Index:    int64(sec / (value * int64(SpanToDuration(span)/time.Second))),
	}
}
