package bitflyer

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/util"
)

func TestParse(t *testing.T) {
	parser := NewParser()
	/*
		{"ts":"2019-12-22T22:07:35.12103609+09:00" ...
			{"channel":"lightning_executions_FX_BTC_JPY","message":
				[
					{"id":1474698784,"side":"BUY","price":795009.0,"size":0.02,"exec_date":"2019-12-22T13:07:34.8995155Z","buy_child_order_acceptance_id":"JRF20191222-130734-472319","sell_child_order_acceptance_id":"JRF20191222-130734-064154"},
					{"id":1474698785,"side":"BUY","price":795017.0,"size":0.06,"exec_date":"2019-12-22T13:07:34.8995155Z","buy_child_order_acceptance_id":"JRF20191222-130734-472319","sell_child_order_acceptance_id":"JRF20191222-130734-399460"},
					{"id":1474698786,"side":"BUY","price":795033.0,"size":0.10501928,"exec_date":"2019-12-22T13:07:34.8995155Z","buy_child_order_acceptance_id":"JRF20191222-130734-472319","sell_child_order_acceptance_id":"JRF20191222-130733-472314"},
					{"id":1474698787,"side":"BUY","price":795036.0,"size":0.01,"exec_date":"2019-12-22T13:07:34.8995155Z","buy_child_order_acceptance_id":"JRF20191222-130734-472319","sell_child_order_acceptance_id":"JRF20191222-130732-399445"},
					{"id":1474698788,"side":"BUY","price":795036.0,"size":0.18081405,"exec_date":"2019-12-22T13:07:34.8995155Z","buy_child_order_acceptance_id":"JRF20191222-130734-472319","sell_child_order_acceptance_id":"JRF20191222-130733-523010"},
					{"id":1474698789,"side":"SELL","price":794961.0,"size":0.01,"exec_date":"2019-12-22T13:07:34.9151456Z","buy_child_order_acceptance_id":"JRF20191222-130733-037464","sell_child_order_acceptance_id":"JRF20191222-130734-523014"}
				]
			}
		}
	*/
	jsonStr := `{"ts":"2019-12-22T22:07:35.12103609+09:00","jsonrpc":"2.0","method":"channelMessage","params":{"channel":"lightning_executions_FX_BTC_JPY","message":[{"id":1474698784,"side":"BUY","price":795009.0,"size":0.02,"exec_date":"2019-12-22T13:07:34.8995155Z","buy_child_order_acceptance_id":"JRF20191222-130734-472319","sell_child_order_acceptance_id":"JRF20191222-130734-064154"},{"id":1474698785,"side":"BUY","price":795017.0,"size":0.06,"exec_date":"2019-12-22T13:07:34.8995155Z","buy_child_order_acceptance_id":"JRF20191222-130734-472319","sell_child_order_acceptance_id":"JRF20191222-130734-399460"},{"id":1474698786,"side":"BUY","price":795033.0,"size":0.10501928,"exec_date":"2019-12-22T13:07:34.8995155Z","buy_child_order_acceptance_id":"JRF20191222-130734-472319","sell_child_order_acceptance_id":"JRF20191222-130733-472314"},{"id":1474698787,"side":"BUY","price":795036.0,"size":0.01,"exec_date":"2019-12-22T13:07:34.8995155Z","buy_child_order_acceptance_id":"JRF20191222-130734-472319","sell_child_order_acceptance_id":"JRF20191222-130732-399445"},{"id":1474698788,"side":"BUY","price":795036.0,"size":0.18081405,"exec_date":"2019-12-22T13:07:34.8995155Z","buy_child_order_acceptance_id":"JRF20191222-130734-472319","sell_child_order_acceptance_id":"JRF20191222-130733-523010"},{"id":1474698789,"side":"SELL","price":794961.0,"size":0.01,"exec_date":"2019-12-22T13:07:34.9151456Z","buy_child_order_acceptance_id":"JRF20191222-130733-037464","sell_child_order_acceptance_id":"JRF20191222-130734-523014"}]}}`
	execStats, err := parser.Parse(util.Bytes(jsonStr))
	assert.Nil(t, err)
	assert.Equal(t, 6, len(execStats))
	for _, stat := range execStats {
		fmt.Println(stat.Ts.Time())
	}
}
func TestTime(t *testing.T) {
	ti := "2019-12-22T13:07:34.8995155Z"
	ti2, err := time.Parse(time.RFC3339Nano, ti)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ti2)

	jst, _ := time.LoadLocation("Asia/Tokyo")
	fmt.Println(ti2.In(jst))

	ti2, err = time.ParseInLocation(time.RFC3339Nano, ti, jst)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ti2)
}
