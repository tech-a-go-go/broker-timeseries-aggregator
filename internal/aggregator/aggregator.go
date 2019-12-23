package aggregator

import (
	"fmt"

	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/aggregator/brokers"
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/aggregator/brokers/bitflyer"
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/aggregator/brokers/bitmex"
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/aggregator/loader"
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
	var stats []*brokers.ExecStat
	var err error
	allStats := make([]*brokers.ExecStat, 100)
	// var aggregatedStats map[clock.TimeSpan][]*brokers.AggregatedStat
	// var processingSpans map[clock.TimeSpan]*brokers.AggregatedStat
	aggregatedStats := make(map[clock.TimeSpan][]*brokers.AggregatedStat)
	processingSpans := make(map[clock.TimeSpan]*brokers.AggregatedStat)
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
			stats, err = brokerParser.Parse(data)

			if err != nil {
				//logger.Error().Err(err).Msg("Error on brokerParser.Parse()")
				break
			}

			for _, stat := range stats {
				for _, span := range a.spans {
					aggreStat, found := processingSpans[span]
					if !found {
						timeIndex := stat.GetTimeIndex(span)
						aggreStat = brokers.NewAggregatedStat(timeIndex)
						processingSpans[span] = aggreStat
					} else {
						timeIndex := stat.GetTimeIndex(span)
						aggreTimeIndex := aggreStat.TimeIndex
						if aggreTimeIndex.Equal(timeIndex) {
							// TODO: aggreStateが現在のtimeIndexと同じなのでaggreTimeIndexを必要があれば更新する
						} else {
							// TODO: for-loopでaggreTimeIndexのIndexがtimeIndexになるまでIndexをインクリメントしたaggreStatを作成して保存
						}
					}

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
