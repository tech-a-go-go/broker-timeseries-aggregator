package aggregator

import (
	"fmt"

	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/aggregator/brokers"
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/aggregator/brokers/bitflyer"
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/aggregator/loader"
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/aggregator/types"
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/clock"
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/log"
)

var logger = log.GetLogger()

type BrokerType int64

const (
	BrokerType_Bitflyer BrokerType = iota
	BrokerType_Bitmex
)

type Aggregator struct {
	broker     BrokerType
	spans      []clock.TimeSpan
	dataLoader *loader.DataLoader
}

func NewAggregator(broker BrokerType, spans []clock.TimeSpan, dataPath string) *Aggregator {
	return &Aggregator{
		broker:     broker,
		spans:      spans,
		dataLoader: loader.NewDataLoader(dataPath),
	}
}

func (a *Aggregator) getBrokerParser() brokers.BrokerParserInterface {
	if a.broker == BrokerType_Bitflyer {
		return bitflyer.NewParser()
	} else if a.broker == BrokerType_Bitmex {
		return nil
	}
	panic("Unknown broker type")
}

func (a *Aggregator) Start() {
	filenams, err := a.dataLoader.GetDataFilenames()
	if err != nil {
		logger.Error().Err(err).Msg("")
	}
	logger.Info().Int("size", len(filenams)).Msg("broker data archives")
	if len(filenams) == 0 {
		logger.Info().Msg("archives not found. exit.")
		return
	}
	go a.dataLoader.Load()
	a.startLoop()
}

func (a *Aggregator) startLoop() {
	brokerParser := a.getBrokerParser()
	dataCh := a.dataLoader.GetDataCh()
	endCh := a.dataLoader.GetEndCh()
	var execStats []*types.ExecStat
	var err error
	accumulatorMap := make(map[clock.TimeSpan][]*Accumulator)

	i := 0
L:
	for {
		select {
		case data := <-dataCh: // 1行ずつデータを受け取る
			if len(data) == 0 {
				break
			}
			i++
			if i%10000 == 0 {
				fmt.Println(i)
			}
			execStats, err = brokerParser.Parse(data)

			if err != nil {
				//logger.Error().Err(err).Msg("Error on brokerParser.Parse()")
				break
			}

			for _, execStat := range execStats {
				for _, span := range a.spans { // 1m, 5m, ...
					timeIndex := execStat.GetTimeIndex(span)
					accums, found := accumulatorMap[span]
					if !found {
						accums = make([]*Accumulator, 0, 32)
						accum := NewAccumulator(timeIndex)
						accums = append(accums, accum)
						accumulatorMap[span] = accums
					}
					accum := accums[len(accums)-1]
					if !accum.IsSameTimeIndex(timeIndex) {
						// 今回のTimeIndexが前回までと異なるということは前回までの
						// Accumulatorのデータが全て溜まったということなので集計する
						accum.Calculate()
						//fmt.Printf("Span=%v, Count=%v, Time=%v\n", span, accum.Count(), accum.timeIndex.Time())
						fmt.Println(accum.ToString())
						accum = NewAccumulator(timeIndex)
						accums = append(accums, accum)
						accumulatorMap[span] = accums
					}
					accum.Add(execStat)
				}
			}
		case <-endCh:
			break L
		}
	}

	// 全てのデータの読み込みが完了したので accumulatorMap に残っている accumulator の集約を行って終了
	for _, span := range a.spans { // 1m, 5m, ...
		accums, found := accumulatorMap[span]
		if !found {
			continue
		}
		accum := accums[len(accums)-1]
		accum.Calculate()
		fmt.Println(accum.ToString())
	}
}
