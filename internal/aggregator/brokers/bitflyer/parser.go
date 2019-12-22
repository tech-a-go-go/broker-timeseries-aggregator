package bitflyer

import (
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/aggregator/brokers"
)

const EXECUTIONS = "lightning_executions_FX_BTC_JPY"

type Parser struct {
}

/*
{"ts":"2019-12-22T22:07:35.12103609+09:00","jsonrpc":"2.0","method":"channelMessage",
	"params":{"channel":"lightning_executions_FX_BTC_JPY",
	"message":[{"id":1474698784,"side":"BUY","price":795009.0,"size":0.02,"exec_date":"2019-12-22T13:07:34.8995155Z","buy_child_order_acceptance_id":"JRF20191222-130734-472319","sell_child_order_acceptance_id":"JRF20191222-130734-064154"},{"id":1474698785,"side":"BUY","price":795017.0,"size":0.06,"exec_date":"2019-12-22T13:07:34.8995155Z","buy_child_order_acceptance_id":"JRF20191222-130734-472319","sell_child_order_acceptance_id":"JRF20191222-130734-399460"},{"id":1474698786,"side":"BUY","price":795033.0,"size":0.10501928,"exec_date":"2019-12-22T13:07:34.8995155Z","buy_child_order_acceptance_id":"JRF20191222-130734-472319","sell_child_order_acceptance_id":"JRF20191222-130733-472314"},{"id":1474698787,"side":"BUY","price":795036.0,"size":0.01,"exec_date":"2019-12-22T13:07:34.8995155Z","buy_child_order_acceptance_id":"JRF20191222-130734-472319","sell_child_order_acceptance_id":"JRF20191222-130732-399445"},{"id":1474698788,"side":"BUY","price":795036.0,"size":0.18081405,"exec_date":"2019-12-22T13:07:34.8995155Z","buy_child_order_acceptance_id":"JRF20191222-130734-472319","sell_child_order_acceptance_id":"JRF20191222-130733-523010"},{"id":1474698789,"side":"SELL","price":794961.0,"size":0.01,"exec_date":"2019-12-22T13:07:34.9151456Z","buy_child_order_acceptance_id":"JRF20191222-130733-037464","sell_child_order_acceptance_id":"JRF20191222-130734-523014"}]}}
*/
func (p *Parser) Parse(json []byte) ([]*brokers.ParsedValue, error) {
	// if bytes.Contains(json, util.Bytes(EXECUTIONS)) {
	// 	json.GetStringInJson()
	// }
	return nil, nil
}
