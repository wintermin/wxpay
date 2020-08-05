package wxpay

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestClient_UnifiedOrder(t *testing.T) {
	c := NewClient("", "", "", Md5, http.DefaultClient)
	c.apiInterceptor.Add(&testLog{})
	resp, err := c.UnifiedOrder(context.Background(), &UnifiedOrderRequest{
		Body:           "测试商品",
		Detail:         "测试商品详情",
		OutTradeNo:     "200805191943383785",
		TotalFee:       1,
		SpBillCreateIp: "127.0.0.1",
		TimeStart:      time.Now(),
		TimeExpire:     time.Now().Add(time.Hour),
		NotifyUrl:      "https://test.com/pay/notify/200805191943383785",
		TradeType:      Native,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
