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

func (e *ExecStat) GetTimeIndex(timeSpan clock.TimeSpan) clock.TimeIndex {
	return e.Ts.MakeTimeIndex(timeSpan.Value, timeSpan.Span)
}

type AggregatedStat struct {
	TimeIndex clock.TimeIndex
	Open      decimal.Decimal
	Close     decimal.Decimal
	Min       decimal.Decimal
	Max       decimal.Decimal
	Size      decimal.Decimal
	BuySize   decimal.Decimal
	SellSize  decimal.Decimal
}

func NewAggregatedStat(timeIndex clock.TimeIndex) *AggregatedStat {
	return &AggregatedStat{
		TimeIndex: timeIndex,
		Open:      decimal.Zero,
		Close:     decimal.Zero,
		Min:       decimal.Zero,
		Max:       decimal.Zero,
		Size:      decimal.Zero,
		BuySize:   decimal.Zero,
		SellSize:  decimal.Zero,
	}
}
