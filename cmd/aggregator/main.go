package main

import (
	"fmt"

	"github.com/tech-a-go-go/timeseries-aggregator/internal/log"
)

var logger = log.GetLogger()

func main() {
	fmt.Println("Hello Aggregator")
	logger.Info().Timestamp().Msg("Hello Aggregator")
}
