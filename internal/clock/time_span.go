package clock

// TimeSpan は時間の間隔を表す.
// value が 2 で span が SPAN_MINUTE であれば2分の間隔を表す.
type TimeSpan struct {
	Value int64
	Span  SPAN
}
