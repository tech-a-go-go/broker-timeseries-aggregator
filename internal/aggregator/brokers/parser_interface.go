package brokers

import "github.com/tech-a-go-go/broker-timeseries-aggregator/internal/aggregator/types"

type BrokerParserInterface interface {
	Parse(json []byte) ([]*types.ExecStat, error)
}
