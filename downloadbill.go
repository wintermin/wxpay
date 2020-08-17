package wxpay

import (
	"bytes"
	"context"
	"net/http"
)

type DownloadBillRequest struct {
	BillDate string   `map:"bill_date"`           //下载对账单的日期，格式：20140603
	BillType BillType `map:"bill_type,omitempty"` //ALL（默认值），返回当日所有订单信息（不含充值退款订单） SUCCESS，返回当日成功支付的订单（不含充值退款订单） REFUND，返回当日退款订单（不含充值退款订单） RECHARGE_REFUND，返回当日充值退款订单
	TarType  string   `map:"tar_type,omitempty"`  //压缩账单 非必传参数，固定值：GZIP，返回格式为.gzip的压缩包账单。不传则默认为数据流形式。
}

func (req DownloadBillRequest) Api() string {
	return ApiDownloadBill
}

func (req DownloadBillRequest) IgnoreKey() []string {
	return nil
}

func (req DownloadBillRequest) IsNeedCert() bool {
	return false
}

// GetDownloadBillHttpRequest 获取下载交易账单http参数 无需证书 https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_6&index=8
// 再使用httpClient.Do(request) 自行解析账单内容
func (c *Client) GetDownloadBillHttpRequest(ctx context.Context, request *DownloadBillRequest) (httpRequest *http.Request, err error) {
	buff := &bytes.Buffer{}
	buff.Write(c.Encode(request).Encode())
	httpRequest, err = http.NewRequestWithContext(ctx, http.MethodPost, MasterApiBaseUrl+request.Api(), buff)
	httpRequest.Header.Set(HttpContentType, HttpContentTypeXml)
	return httpRequest, err
}
