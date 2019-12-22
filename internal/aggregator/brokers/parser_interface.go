package brokers

import (
	"github.com/shopspring/decimal"
	"time"
)

type BrokerParserInterface interface {
	Parse(json []byte) ([]*ParsedValue, error)
}

type ParsedValue struct {
	ts      time.Time
	price   decimal.Decimal
	priceOk bool
	size    decimal.Decimal
	sizeOk  bool
}
