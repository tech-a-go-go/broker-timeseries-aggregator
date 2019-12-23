package clock

import (
	"testing"
)

func TestCreateTimeIndex(t *testing.T) {
	test := []struct {
		sec      int64
		value    int64
		span     SPAN
		expected int64
	}{
		{1, 1, SPAN_SECOND, 1},
		{0, 1, SPAN_SECOND, 0},
		{60, 1, SPAN_SECOND, 60},
		{60, 5, SPAN_SECOND, 12},

		{59, 1, SPAN_MINUTE, 0},
		{60, 1, SPAN_MINUTE, 1},
		{61, 1, SPAN_MINUTE, 1},
		{610, 1, SPAN_MINUTE, 10},
		{610, 5, SPAN_MINUTE, 2},
	}

	for i, tt := range test {
		timeIndex := CreateTimeIndex(tt.sec, tt.value, tt.span)
		if tt.expected != timeIndex.Index {
			t.Errorf("Error at %v : %v != %v", i, tt.expected, timeIndex.Index)
		}
	}
}
