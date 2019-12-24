package bitmex

import "github.com/tech-a-go-go/broker-timeseries-aggregator/internal/aggregator/types"

type Parser struct {
}

func (p *Parser) Parse(jsonBytes []byte) ([]*types.ExecStat, error) {
	return nil, nil
}
