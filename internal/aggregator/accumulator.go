package aggregator

import (
	"github.com/shopspring/decimal"
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/aggregator/types"
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/clock"
)

type accumulrator struct {
	TimeIndex      clock.TimeIndex
	ExecStats      []*types.ExecStat
	AggregatedStat types.AggregatedStat
}

func (a *accumulrator) Start() {
	open := decimal.Zero
	close := decimal.Zero
	min := decimal.NewFromFloat(9999999999)
	max := decimal.Zero

	size := decimal.Zero
	buySize := decimal.Zero
	sellSize := decimal.Zero

	len := len(a.ExecStats)
	for i := 0; i < len; i++ {
		stat := a.ExecStats[i]
		if stat.PriceOk && open.IsZero() {
			open = stat.Price
		}
		if stat.PriceOk {
			close = stat.Price
		}
		if stat.PriceOk && stat.Price.LessThan(min) {
			min = stat.Price
		}
		if stat.PriceOk && stat.Price.GreaterThan(max) {
			max = stat.Price
		}
		if stat.SizeOk {
			size = size.Add(stat.Size)
			if stat.Side == types.SIDE_BUY {
				buySize = buySize.Add(stat.Size)
			} else {
				sellSize = sellSize.Add(stat.Size)
			}
		}
	}
	a.AggregatedStat = types.AggregatedStat{
		TimeIndex: a.TimeIndex,
		Open:      open,
		Close:     close,
		Min:       min,
		Max:       max,
		Size:      size,
		BuySize:   buySize,
		SellSize:  sellSize,
	}
}
