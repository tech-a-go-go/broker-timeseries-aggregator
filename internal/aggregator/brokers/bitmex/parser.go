package bitmex

import (
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/aggregator/brokers"
)

type Parser struct {
}

func (p *Parser) Parse(jsonBytes []byte) ([]*brokers.ExecStat, error) {
	return nil, nil
}
