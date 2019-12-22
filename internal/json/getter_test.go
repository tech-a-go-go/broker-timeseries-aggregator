package json

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/tech-a-go-go/broker-timeseries-aggregator/internal/util"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fastjson"
)

func TestGetStringInJson(t *testing.T) {
	test := []struct {
		jsonStr   string
		param     string
		value     string
		endOffset int
		ok        bool
	}{
		{`{"jsonrpc":"2.0","method":"channelMessage","params":{"channel":"lightning_ticker_FX_BTC_JPY","message":{"product_code":"FX_BTC_JPY"`, "channel", "lightning_ticker_FX_BTC_JPY", 91, true},
		{`{"jsonrpc":"2.0","action":"update","params":{"channel":"lightning_ticker_FX_BTC_JPY","message":{"product_code":"FX_BTC_JPY"`, "channel", "lightning_ticker_FX_BTC_JPY", 83, true},
		{`{"action":"update","table":"orderBookL2_25","data":[{"symbol":"XBTUSD","id":8799294600,"side":"Sell","size":837797}]}`, "table", "orderBookL2_25", 42, true},
		{`{"action":"update","data":[{"symbol":"XBTUSD","id":0,"side":"Sell","size":837797}]}`, "table", "", 0, false},
		{`{}`, "table", "", 0, false},
	}

	var bytes [64]byte
	buf := bytes[:0]
	for _, tt := range test {
		//buf = buf[:0]
		param := util.Bytes(tt.param)
		fmt.Printf("1 %s\n", buf)
		value, endOffset, ok := GetStringInJson(util.Bytes(tt.jsonStr), param, 0, buf)
		fmt.Printf("2 %s\n", buf)
		fmt.Printf("3 %s\n", value)
		if ok == tt.ok {
			if ok {
				if util.String(value) != tt.value {
					t.Errorf("Error id : %v, %v", tt.param, util.String(value))
				}
				if endOffset != tt.endOffset {
					t.Errorf("Error endOffset : %v, %v", tt.endOffset, endOffset)
				}
			}
		} else {
			t.Errorf("Error ok : %v, %v", tt.ok, ok)
		}
	}
}

func TestGetNumberInJson(t *testing.T) {
	test := []struct {
		jsonStr   string
		idValue   string
		endOffset int
		ok        bool
	}{
		{`{"table":"orderBookL2_25","action":"update","data":[{"symbol":"XBTUSD","id":8799294600,"side":"Sell","size":837797}]}`, "8799294600", 86, true},
		{`{"table":"orderBookL2_25","action":"update","data":[{"symbol":"XBTUSD","id":0,"side":"Sell","size":837797}]}`, "0", 77, true},
		{`{"id":8799294600}`, "8799294600", 16, true},
		{`{"id":16.0}`, "16.0", 10, true},
		{`{"id":1.6E-35}`, "1.6E-35", 13, true},
		{`{"id":-200.12}`, "-200.12", 13, true},
		{`{"iid":-200.12}`, "", 0, false}, // wrong param name iid
		{`{"table":"orderBookL2_25","action":"update","data":[{"symbol":"XBTUSD","side":"Sell","size":837797}]}`, "0", 0, false},
		{`{}`, "0", 0, false},
	}
	idBytes := util.Bytes("id")
	buf := bytebufferpool.Get()
	for _, tt := range test {
		buf.Reset()
		idValue, endOffset, ok := GetNumberInJson(util.Bytes(tt.jsonStr), idBytes, 0, buf)
		if ok == tt.ok {
			if ok {
				if util.String(idValue) != tt.idValue {
					t.Errorf("Error id : %v, %v", tt.idValue, util.String(idValue))
				}
				if endOffset != tt.endOffset {
					t.Errorf("Error endOffset : %v, %v", tt.endOffset, endOffset)
				}
			}
		} else {
			t.Errorf("Error ok : %v, %v", tt.ok, ok)
		}
	}
}

func TestS(t *testing.T) {
	boardJson := `[{"product_code":"FX_BTC_JPY","side":"SELL","price":789613.000000000000,"size":0.010000000000,"commission":0.000000000000,"swap_point_accumulate":0.0,"require_collateral":1974.032500000000,"open_date":"2019-12-19T07:40:23.66","leverage":4.000000000000,"pnl":1.430000000000000000000000,"sfd":0.000000000000},{"product_code":"FX_BTC_JPY","side":"SELL","price":789495.000000000000,"size":0.010000000000,"commission":0.000000000000,"swap_point_accumulate":0.0,"require_collateral":1973.737500000000,"open_date":"2019-12-19T07:40:36.347","leverage":4.000000000000,"pnl":0.250000000000000000000000,"sfd":0.000000000000}]`
	values, _ := fastjson.Parse(boardJson)
	elems, _ := values.Array()
	sizeDecimal := decimal.Zero
	for _, elem := range elems {
		size := elem.GetFloat64("size")
		sizeDecimal = sizeDecimal.Add(decimal.NewFromFloat(size))
	}
	sum, _ := sizeDecimal.Float64()
	if float64(0.02) != sum {
		t.Errorf("Error : %v, %v", float64(0.02), sum)
	}
}
