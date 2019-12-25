package clock

import "time"

// TimeSpan は時間の長さを value(大きさ) と span(期間の単位) で表す.
// value が 2 で span が SPAN_MINUTE であれば2分間を表す.
type TimeSpan struct {
	Value int64 // 大きさ
	Span  SPAN  // SECOND, MINUTEなど期間の単位
}

// Duration このTimeSpanの期間を返す
func (t *TimeSpan) Duration() time.Duration {
	return time.Duration(t.Value) * SpanToDuration(t.Span)
}
