package wxpay

import (
	"fmt"
	"net/http"
	"testing"
)

func TestClient_Notify(t *testing.T) {
	c := NewClient( /*your appId*/ "" /*your mchId*/, "" /*your key*/, "", Md5, http.DefaultClient)
	resp, err := c.Notify([]byte( /*you body*/ ""))
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
