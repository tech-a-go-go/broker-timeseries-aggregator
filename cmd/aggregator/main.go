package main

import (
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/aggregator"
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

	aggr := aggregator.NewAggregator(aggregator.BrokerType_Bitflyer, "./datda/bitflyer-20d19*.log")
	aggr.Start()
}
