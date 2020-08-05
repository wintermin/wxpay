package wxpay

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestClient_OrderQuery(t *testing.T) {
	c := NewClient("", "", "", Md5, http.DefaultClient)
	c.apiInterceptor.Add(&testLog{})
	resp, err := c.OrderQuery(context.Background(), &OrderQueryRequest{
		OutTradeNo:    "",
		TransactionId: "",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
