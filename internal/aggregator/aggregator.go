package aggregator

import (
	"fmt"

	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/aggregator/brokers"
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/aggregator/brokers/bitflyer"
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/aggregator/brokers/bitmex"
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
		return &bitmex.Parser{}
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
	allStats := make([]*types.ExecStat, 100)
	// var aggregatedStats map[clock.TimeSpan][]*brokers.AggregatedStat
	// var processingSpans map[clock.TimeSpan]*brokers.AggregatedStat
	execStatSpanMap := make(map[clock.TimeSpan][][]*types.ExecStat)
	// TimeSpan毎に集約するデータを保存する配列を用意
	for _, span := range a.spans {
		execStatSpanMap[span] = make([][]*types.ExecStat, 128)
	}

	i := 0
L:
	for {
		select {
		case data := <-dataCh:
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

			// 1m = 1 2 3
			// 5m = 1 2 3
			for _, execStat := range execStats {
				for _, span := range a.spans { // 1m, 5m, ...
					aggreStat, found := execStatSpanMap[span]
					_ = aggreStat
					if !found {
						timeIndex := execStat.GetTimeIndex(span)
						aggreStat := types.NewAggregatedStat(timeIndex)
						_ = aggreStat
						// processingSpans[span] = aggreStat
					}
					execTimeIndex := execStat.GetTimeIndex(span)
					_ = execTimeIndex
					// aggreTimeIndex := aggreStat.TimeIndex
					// if execTimeIndex.Equal(aggreTimeIndex) {

					// 	// TODO: aggreStateが現在のtimeIndexと同じなのでaggreTimeIndexを必要があれば更新する
					// } else {
					// 	// TODO: for-loopでaggreTimeIndexのIndexがtimeIndexになるまでIndexをインクリメントしたaggreStatを作成して保存
					// }

				}

			}

			// allStats = append(allStats, stat)
			// 1m, 5m, 10m,
			// price = open, close, max, min
			// volume = sum
			//stat.

		case <-endCh:
			break L
		}
	}
	fmt.Printf("all stats = %v\n", len(allStats))
}
