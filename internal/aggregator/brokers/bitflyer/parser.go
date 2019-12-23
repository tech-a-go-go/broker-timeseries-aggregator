package bitflyer

import (
	"bytes"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/aggregator/brokers"
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/clock"
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/json"
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/util"
	"github.com/valyala/fastjson"
)

var CHANNEL_EXECUTIONS_BYTES = util.Bytes("lightning_executions_FX_BTC_JPY")

var NO_DATA_ERROR = errors.New("NO DATA")

type Parser struct {
	parser   *fastjson.Parser
	location *time.Location
}

func NewParser() *Parser {
	var p fastjson.Parser
	JST, _ := time.LoadLocation("Asia/Tokyo")
	return &Parser{
		parser:   &p,
		location: JST,
	}
}

/*
{"ts":"2019-12-22T22:07:35.12103609+09:00","jsonrpc":"2.0","method":"channelMessage",
	"params":{"channel":"lightning_executions_FX_BTC_JPY",
	"message":[{"id":1474698784,"side":"BUY","price":795009.0,"size":0.02,"exec_date":"2019-12-22T13:07:34.8995155Z","buy_child_order_acceptance_id":"JRF20191222-130734-472319","sell_child_order_acceptance_id":"JRF20191222-130734-064154"},{"id":1474698785,"side":"BUY","price":795017.0,"size":0.06,"exec_date":"2019-12-22T13:07:34.8995155Z","buy_child_order_acceptance_id":"JRF20191222-130734-472319","sell_child_order_acceptance_id":"JRF20191222-130734-399460"},{"id":1474698786,"side":"BUY","price":795033.0,"size":0.10501928,"exec_date":"2019-12-22T13:07:34.8995155Z","buy_child_order_acceptance_id":"JRF20191222-130734-472319","sell_child_order_acceptance_id":"JRF20191222-130733-472314"},{"id":1474698787,"side":"BUY","price":795036.0,"size":0.01,"exec_date":"2019-12-22T13:07:34.8995155Z","buy_child_order_acceptance_id":"JRF20191222-130734-472319","sell_child_order_acceptance_id":"JRF20191222-130732-399445"},{"id":1474698788,"side":"BUY","price":795036.0,"size":0.18081405,"exec_date":"2019-12-22T13:07:34.8995155Z","buy_child_order_acceptance_id":"JRF20191222-130734-472319","sell_child_order_acceptance_id":"JRF20191222-130733-523010"},{"id":1474698789,"side":"SELL","price":794961.0,"size":0.01,"exec_date":"2019-12-22T13:07:34.9151456Z","buy_child_order_acceptance_id":"JRF20191222-130733-037464","sell_child_order_acceptance_id":"JRF20191222-130734-523014"}]}}
*/
func (p *Parser) Parse(jsonBytes []byte) ([]*brokers.ExecStat, error) {
	var bufBytes [64]byte
	buf := bufBytes[:0]
	channel, _, ok := json.GetStringInJson(jsonBytes, util.Bytes("channel"), 0, buf)
	if !ok {
		return nil, NO_DATA_ERROR
	}
	if !bytes.Equal(CHANNEL_EXECUTIONS_BYTES, channel) {
		return nil, NO_DATA_ERROR
	}
	buf = buf[:0]
	value, err := p.parser.ParseBytes(jsonBytes)
	if err != nil {
		return nil, err
	}
	// tsBytes := value.GetStringBytes("ts")
	// ts, err := time.Parse(time.RFC3339Nano, util.String(tsBytes))
	// if err != nil {
	// 	return nil, err
	// }
	// timestamp := brokers.MakeSimpleTime(ts)
	execValues := value.GetArray("params", "message")
	execStats := make([]*brokers.ExecStat, 0, len(execValues))
	for _, execValue := range execValues {
		// execStat := brokers.GetExecStat()
		execStat := &brokers.ExecStat{}
		execDataBytes, err := execValue.Get("exec_date").StringBytes()
		if err != nil {
			return nil, err
		}
		tsInUTC, err := time.Parse(time.RFC3339Nano, util.String(execDataBytes))
		if err != nil {
			return nil, err
		}
		tsInJST := tsInUTC.In(p.location)
		timestamp := clock.MakeSimpleTime(tsInJST)
		execStat.Ts = timestamp
		sideBytes, err := execValue.Get("side").StringBytes()
		if err != nil {
			return nil, errors.Wrap(err, tsInJST.String())
		}
		if util.String(sideBytes) == "BUY" {
			execStat.Side = brokers.SIDE_BUY
		} else {
			execStat.Side = brokers.SIDE_SELL
		}
		price, err := execValue.Get("price").Float64()
		if err != nil {
			return nil, errors.Wrap(err, tsInJST.String())
		}
		execStat.Price = decimal.NewFromFloat(price)
		execStat.PriceOk = true

		size, err := execValue.Get("size").Float64()
		if err != nil {
			return nil, errors.Wrap(err, tsInJST.String())
		}
		execStat.Size = decimal.NewFromFloat(size)
		execStat.SizeOk = true

		execStats = append(execStats, execStat)
	}
	return execStats, nil
}
