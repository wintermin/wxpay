package wxpay

import (
	"fmt"
	"net/http"
	"testing"
)

func TestClient_NonceStr(t *testing.T) {
	c := NewClient("", "", "", "", http.DefaultClient)
	fmt.Println(c.NonceStr())
}
