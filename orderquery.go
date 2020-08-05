package wxpay

import (
	"context"
)

type OrderQueryRequest struct {
	TransactionId string `map:"transaction_id,omitempty"` //微信订单号	二选一	微信的订单号，优先使用
	OutTradeNo    string `map:"out_trade_no,omitempty"`   //商户订单号	商户系统内部的订单号，当没提供transaction_id时需要传这个。
}

func (req OrderQueryRequest) Api() string {
	return ApiOrderQuery
}

func (req OrderQueryRequest) IgnoreKey() []string {
	return []string{KeySignType}
}

type OrderQueryResponse struct {
	Response
	DeviceInfo         string            `map:"device_info,omitempty"`          //微信支付分配的终端设备号
	OpenId             string            `map:"openid,omitempty"`               //用户在商户appid下的唯一标识
	IsSubscribe        string            `map:"is_subscribe,omitempty"`         //用户是否关注公众账号，Y-关注，N-未关注
	TradeType          TradeType         `map:"trade_type,omitempty"`           //调用接口提交的交易类型
	TradeState         TradeState        `map:"trade_state,omitempty"`          //调用接口提交的交易类型
	BankType           string            `map:"bank_type,omitempty"`            //付款银行类型，采用字符串类型的银行标识
	TotalFee           int               `map:"total_fee,omitempty"`            //订单总金额，单位为分
	FeeType            string            `map:"fee_type,omitempty"`             //货币类型，符合ISO 4217标准的三位字母代码，默认人民币：CNY
	CashFee            int               `map:"cash_fee,omitempty"`             //现金支付金额订单现金支付金额，详见支付金额
	CashFeeType        string            `map:"cash_fee_type,omitempty"`        //现金支付货币类型
	SettlementTotalFee int               `map:"settlement_total_fee,omitempty"` //应结订单金额 当订单使用了免充值型优惠券后返回该参数，应结订单金额=订单金额-免充值优惠券金额。
	CouponFee          int               `map:"coupon_fee,omitempty"`           //代金券金额
	CouponCount        int               `map:"coupon_count,omitempty"`         //代金券或立减优惠使用数量
	CouponId           map[string]string `map:"coupon_id,omitempty"`            //代金券ID
	CouponType         map[string]string `map:"coupon_type,omitempty"`          //代金券类型
	TransactionId      string            `map:"transaction_id,omitempty"`       //微信支付订单号
	OutTradeNo         string            `map:"out_trade_no,omitempty"`         //商户订单号
	Attach             string            `map:"attach,omitempty"`               //附加数据
	TimeEnd            string            `map:"time_end,omitempty"`             //支付完成时间
	TradeStateDesc     string            `map:"trade_state_desc,omitempty"`     //交易状态描述
}

// OrderQuery 查询订单 无需证书 https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_2&index=4
func (c *Client) OrderQuery(ctx context.Context, request *OrderQueryRequest) (resp *OrderQueryResponse, err error) {
	resp = &OrderQueryResponse{}
	_, err = c.Do(ctx, MasterApiBaseUrl, request, resp)
	return resp, err
}
