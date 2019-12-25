package aggregator

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/aggregator/types"
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/clock"
)

// Accumulator データの集約を行う
type Accumulator struct {
	timeIndex      clock.TimeIndex
	execStats      []*types.ExecStat
	aggregatedStat types.AggregatedStat
}

func NewAccumulator(timeIndex clock.TimeIndex) *Accumulator {
	return &Accumulator{
		timeIndex: timeIndex,
		execStats: make([]*types.ExecStat, 0, 128),
	}
}

func (a *Accumulator) IsSameTimeIndex(timeIndex clock.TimeIndex) bool {
	return a.timeIndex.Equal(timeIndex)
}

func (a *Accumulator) Add(stat *types.ExecStat) {
	a.execStats = append(a.execStats, stat)
}

func (a *Accumulator) Count() int {
	return len(a.execStats)
}

func (a *Accumulator) ToString() string {
	return fmt.Sprintf("%v, Open=%v, Close=%v, Min=%v, Max=%v, Size=%v, BuySize=%v, SellSize=%v", a.timeIndex.Time(), a.aggregatedStat.Open, a.aggregatedStat.Close, a.aggregatedStat.Min, a.aggregatedStat.Max, a.aggregatedStat.Size, a.aggregatedStat.BuySize, a.aggregatedStat.SellSize)
}

func (a *Accumulator) GetAggregatedStat() types.AggregatedStat {
	return a.aggregatedStat
}

func (a *Accumulator) Calculate() {
	open := decimal.Zero
	close := decimal.Zero
	min := decimal.NewFromFloat(9999999999)
	max := decimal.Zero
	size := decimal.Zero
	buySize := decimal.Zero
	sellSize := decimal.Zero

	len := len(a.execStats)
	for i := 0; i < len; i++ {
		stat := a.execStats[i]
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
	a.aggregatedStat = types.AggregatedStat{
		TimeIndex: a.timeIndex,
		Open:      open,
		Close:     close,
		Min:       min,
		Max:       max,
		Size:      size,
		BuySize:   buySize,
		SellSize:  sellSize,
	}
}
