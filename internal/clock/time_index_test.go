package clock

import (
	"testing"
	"time"
)

func TestCreateTimeIndex(t *testing.T) {
	test := []struct {
		ts       string
		value    int64
		span     SPAN
		expected int64
	}{
		{"1970-01-01T09:00:00.012777613+09:00", 1, SPAN_SECOND, 0},
		{"1970-01-01T09:00:05.012777613+09:00", 5, SPAN_SECOND, 1},
		{"2019-12-20T19:00:00.012777613+09:00", 1, SPAN_SECOND, 1576836000},
		{"2019-12-20T19:00:00.012777613+09:00", 5, SPAN_SECOND, 315367200},
		{"2019-12-20T19:00:00.012777613+09:00", 15, SPAN_SECOND, 105122400},
		{"2019-12-20T19:00:00.012777613+09:00", 1, SPAN_MINUTE, 26280600},
		{"2019-12-20T19:00:00.012777613+09:00", 5, SPAN_MINUTE, 5256120},
		{"2019-12-20T19:00:00.012777613+09:00", 1, SPAN_HOUR, 438010},
		{"2019-12-20T19:00:00.012777613+09:00", 5, SPAN_HOUR, 87602},
	}
	loc, _ := time.LoadLocation("Asia/Tokyo")
	for i, tt := range test {
		tsInUTC, _ := time.Parse(time.RFC3339Nano, tt.ts)
		tsInJST := tsInUTC.In(loc)
		timestamp := MakeSimpleTime(tsInJST)
		timeIndex := NewTimeIndex(timestamp.Sec, tt.value, tt.span)
		// fmt.Printf("sec=%v, value=%v, span=%v, %v\n", timestamp.Sec, tt.value, tt.span, timeIndex.Time())
		if tt.expected != timeIndex.Index {
			t.Errorf("Error at %v : %v != %v", i, tt.expected, timeIndex.Index)
		}
	}
}
