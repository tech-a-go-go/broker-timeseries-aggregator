package brokers

import (
	"github.com/shopspring/decimal"
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/clock"
)

type SIDE int8

const (
	SIDE_BUY  SIDE = 0
	SIDE_SELL SIDE = 1
)

type BrokerParserInterface interface {
	Parse(json []byte) ([]*ExecStat, error)
}

type ExecStat struct {
	Ts      clock.SimpleTime
	Side    SIDE
	Price   decimal.Decimal
	PriceOk bool
	Size    decimal.Decimal
	SizeOk  bool
}

func (e *ExecStat) GetTimeIndex(value int64, span clock.SPAN) clock.TimeIndex {
	return e.Ts.MakeTimeIndex(value, span)
}
