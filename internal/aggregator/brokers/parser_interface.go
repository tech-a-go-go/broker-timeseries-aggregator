package brokers

import (
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

type SIDE int8

const (
	SIDE_BUY  SIDE = 0
	SIDE_SELL SIDE = 1
)

type BrokerParserInterface interface {
	Parse(json []byte) ([]*ExecStat, error)
}

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

type ExecStat struct {
	Ts      SimpleTime
	Side    SIDE
	Price   decimal.Decimal
	PriceOk bool
	Size    decimal.Decimal
	SizeOk  bool
}

var execStatPool = &sync.Pool{
	New: func() interface{} {
		return new(ExecStat)
	},
}

func GetExecStat() *ExecStat {
	value := execStatPool.Get().(*ExecStat)
	return value
}

func PutExecStat(value *ExecStat) {
	value.PriceOk = false
	value.SizeOk = false
	execStatPool.Put(value)
}
