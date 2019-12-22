package aggregator

import (
	"time"

	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/aggregator/brokers"
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/aggregator/brokers/bitflyer"
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/aggregator/brokers/bitmex"
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/aggregator/loader"
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/log"
)

var logger = log.GetLogger()

type BrokerType int64

const (
	BrokerType_Bitflyer BrokerType = iota
	BrokerType_Bitmex
)

var timeSeriesResolutions = []time.Duration{
	1 * time.Minute,
	3 * time.Minute,
	5 * time.Minute,
	15 * time.Minute,
	30 * time.Minute,
	1 * time.Hour,
}

type Aggregator struct {
	broker     BrokerType
	dataLoader *loader.DataLoader
}

func NewAggregator(broker BrokerType, dataPath string) *Aggregator {
	return &Aggregator{
		broker:     broker,
		dataLoader: loader.NewDataLoader(dataPath),
	}
}

func (a *Aggregator) getBrokerParser() brokers.BrokerParserInterface {
	if a.broker == BrokerType_Bitflyer {
		return &bitflyer.Parser{}
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
	a.startLoop()
}

func (a *Aggregator) startLoop() {
	brokerParser := a.getBrokerParser()
	_ = brokerParser
	dataCh := a.dataLoader.GetDataCh()
	endCh := a.dataLoader.GetEndCh()
L:
	for {
		select {
		case data := <-dataCh:
			if len(data) == 0 {
				break
			}

			// 1m, 5m, 10m,
			// price = open, close, max, min
			// volume = sum
			//parsedValues, err := brokerParser.Parse(data)

		case <-endCh:
			break L
		}
	}
}
