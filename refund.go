package wxpay

import "context"

type RefundRequest struct {
	TransactionId string `map:"transaction_id,omitempty"`  //微信订单号	二选一	微信的订单号，优先使用
	OutTradeNo    string `map:"out_trade_no,omitempty"`    //商户订单号	商户系统内部的订单号，当没提供transaction_id时需要传这个。
	OutRefundNo   string `map:"out_refund_no,omitempty"`   //商户退款单号	商户系统内部的退款单号，商户系统内部唯一，只能是数字、大小写字母_-|*@ ，同一退款单号多次请求只退一笔。
	TotalFee      int    `map:"total_fee"`                 //订单总金额，单位为分
	RefundFee     int    `map:"refund_fee"`                //退款总金额，订单金额，单位为分
	RefundFeeType string `map:"refund_fee_type,omitempty"` //退款货币种类
	RefundDesc    string `map:"refund_desc,omitempty"`     //退款原因	 若商户传入，会在下发给用户的退款消息中体现退款原因 注意：若订单退款金额≤1元，且属于部分退款，则不会在退款消息中体现退款原因
	RefundAccount string `map:"refund_account,omitempty"`  //退款资金来源	REFUND_SOURCE_UNSETTLED_FUNDS---未结算资金退款（默认使用未结算资金退款） REFUND_SOURCE_RECHARGE_FUNDS---可用余额退款
	NotifyUrl     string `map:"notify_url,omitempty"`      //退款结果通知url 异步接收微信支付退款结果通知的回调地址，通知URL必须为外网可访问的url，不允许带参数
}

func (req RefundRequest) Api() string {
	return ApiRefund
}

func (req RefundRequest) IgnoreKey() []string {
	return nil
}

func (req RefundRequest) IsNeedCert() bool {
	return true
}

type RefundResponse struct {
	Response
	TransactionId       string            `map:"transaction_id,omitempty"`        //微信订单号	二选一	微信的订单号，优先使用
	OutTradeNo          string            `map:"out_trade_no,omitempty"`          //商户订单号	商户系统内部的订单号，当没提供transaction_id时需要传这个。
	OutRefundNo         string            `map:"out_refund_no,omitempty"`         //商户退款单号	商户系统内部的退款单号，商户系统内部唯一，只能是数字、大小写字母_-|*@ ，同一退款单号多次请求只退一笔。
	RefundId            string            `map:"refund_id,omitempty"`             //微信退款单号
	TotalFee            int               `map:"total_fee"`                       //订单总金额，单位为分
	RefundFee           int               `map:"refund_fee"`                      //退款总金额，订单金额，单位为分
	SettlementRefundRee int               `map:"settlement_refund_fee,omitempty"` //应结退款金额	去掉非充值代金券退款金额后的退款金额，退款金额=申请退款金额-非充值代金券退款金额，退款金额<=申请退款金额
	FeeType             string            `map:"fee_type,omitempty"`              //标价币种
	CashFee             int               `map:"cash_fee"`                        //现金支付金额
	CashFeeType         string            `map:"cash_fee_type,omitempty"`         //现金支付币种
	CashRefundFee       string            `map:"cash_refund_fee,omitempty"`       //现金退款金额，单位为分
	CouponRefundFee     int               `map:"coupon_refund_fee,omitempty"`     //代金券退款总金额
	CouponType          map[string]string `map:"coupon_type,omitempty"`           //代金券类型
	CouponRefundFeeMap  map[string]string `map:"coupon_refund_fee_,omitempty"`    //单个代金券退款金额
	CouponRefundCount   int               `map:"coupon_refund_count,omitempty"`   //退款代金券使用数量
	CouponRefundId      map[string]string `map:"coupon_refund_id_,omitempty"`     //退款代金券ID, $n为下标，从0开始编号
}

// Refund 申请退款 请求需要双向证书 https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_4
func (c *Client) Refund(ctx context.Context, request *RefundRequest) (resp *RefundResponse, err error) {
	resp = &RefundResponse{}
	_, err = c.Do(ctx, MasterApiBaseUrl, request, resp)
	return resp, err
}
