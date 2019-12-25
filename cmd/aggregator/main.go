package main

import (
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/aggregator"
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/clock"
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/log"
)

var logger = log.GetLogger()

func main() {
	// ld := loader.NewDataLoader("./data/bitflyer-2019*.log")
	// filenames, err := ld.GetDataFilenames()
	// if err != nil {
	// 	fmt.Errorf("%v\n", err.Error())
	// 	return
	// }
	// for _, f := range filenames {
	// 	fmt.Printf("%s\n", f)
	// }

	aggr := aggregator.NewAggregator(aggregator.BrokerType_Bitflyer, []clock.TimeSpan{ /*clock.TimeSpan{Value: 1, Span: clock.SPAN_MINUTE},*/ clock.TimeSpan{Value: 5, Span: clock.SPAN_MINUTE}}, "./data/broker.bitflyer-20191220*.gz")
	aggr.Start()
}
