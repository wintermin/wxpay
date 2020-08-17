package wxpay

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"strings"
)

type RefundNotifyResponse struct {
	Response
	ReqInfo string `map:"req_info,omitempty"` //加密信息
	ReqInfoResponse
}

type ReqInfoResponse struct {
	TransactionId string `xml:"transaction_id"` //微信订单号
	OutTradeNo    string `xml:"out_trade_no"`   //商户订单号
	OutRefundNo   string `xml:"out_refund_no"`  //商户退款单号
	RefundId      string `xml:"refund_id"`      //微信退款单号

	TotalFee           int `xml:"total_fee"`                      //订单总金额，单位为分
	SettlementTotalFee int `xml:"settlement_total_fee,omitempty"` //应结订单金额 当订单使用了免充值型优惠券后返回该参数，应结订单金额=订单金额-免充值优惠券金额。

	RefundFee           int    `xml:"refund_fee"`            //申请退款金额，订单金额，单位为分
	SettlementRefundFee int    `xml:"settlement_refund_fee"` //退款金额
	RefundStatus        string `xml:"refund_status"`         //退款状态

	RefundAccount       string `xml:"refund_account"`                //退款资金来源
	RefundRecvAccout    string `xml:"refund_recv_accout"`            //退款入账账户
	RefundSuccessTime   string `xml:"refund_success_time,omitempty"` //退款成功时间
	RefundRequestSource string `xml:"refund_request_source"`         //退款成功时间
}

// 退款结果通知 返回结果封装 https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_16&index=11
// 响应渠道使用 ResponseSuccess() 或 ResponseFail(msg string)
func (c *Client) RefundNotify(body []byte) (*RefundNotifyResponse, error) {
	resp := &RefundNotifyResponse{}
	val := make(Values)
	val.Decode(body)
	returnCode := val.Get(KeyReturnCode)
	if returnCode == string(ReturnCodeSuccess) {
		val.ToStruct(resp)
		reqInfo := val.Get(KeyReqInfo) //加密报文
		/*
			（1）对加密串A做base64解码，得到加密串B
			（2）对商户key做md5，得到32位小写key* ( key设置路径：微信商户平台(pay.weixin.qq.com)-->账户设置-->API安全-->密钥设置 )
			（3）用key*对加密串B做AES-256-ECB解密（PKCS7Padding）
		*/
		bytes, err := base64.StdEncoding.DecodeString(reqInfo)
		if err != nil {
			return resp, err
		}
		h := md5.New()
		h.Write([]byte(c.Key()))
		md5Key := hex.EncodeToString(h.Sum(nil))
		md5Key = strings.ToLower(md5Key)
		plaintext, err := AesECBDecrypt([]byte(md5Key), bytes)
		if err != nil {
			return resp, err
		}
		rir := &ReqInfoResponse{}
		xml.Unmarshal(plaintext, rir)
		resp.ReqInfoResponse = *rir
	}
	return resp, nil
}
