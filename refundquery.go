package wxpay

import (
	"context"
)

type RefundQueryRequest struct {
	//四选一
	TransactionId string `map:"transaction_id,omitempty"` //微信订单号	二选一	微信的订单号，优先使用
	OutTradeNo    string `map:"out_trade_no,omitempty"`   //商户订单号	商户系统内部的订单号，当没提供transaction_id时需要传这个。
	OutRefundNo   string `map:"out_refund_no,omitempty"`  //商户退款单号	商户系统内部的退款单号，商户系统内部唯一，只能是数字、大小写字母_-|*@ ，同一退款单号多次请求只退一笔。
	RefundId      string `map:"refund_id,omitempty"`      //微信退款单号

	Offset string `map:"offset,omitempty"` //偏移量，当部分退款次数超过10次时可使用，表示返回的查询结果从这个偏移量开始取记录
}

func (req RefundQueryRequest) Api() string {
	return ApiRefundQuery
}

func (req RefundQueryRequest) IgnoreKey() []string {
	return nil
}

func (req RefundQueryRequest) IsNeedCert() bool {
	return false
}

type RefundQueryResponse struct {
	Response
	TotalRefundCount   int    `map:"total_refund_count,omitempty"`   //订单总共已发生的部分退款次数，当请求参数传入offset后有返回
	TransactionId      string `map:"transaction_id,omitempty"`       //微信支付订单号
	OutTradeNo         string `map:"out_trade_no,omitempty"`         //商户订单号
	TotalFee           int    `map:"total_fee,omitempty"`            //订单总金额，单位为分
	SettlementTotalFee int    `map:"settlement_total_fee,omitempty"` //应结订单金额 当订单使用了免充值型优惠券后返回该参数，应结订单金额=订单金额-免充值优惠券金额。
	FeeType            string `map:"fee_type,omitempty"`             //货币类型，符合ISO 4217标准的三位字母代码，默认人民币：CNY
	CashFee            int    `map:"cash_fee,omitempty"`             //现金支付金额订单现金支付金额，详见支付金额
	CashFeeType        string `map:"cash_fee_type,omitempty"`        //现金支付货币类型
	CouponFee          int    `map:"coupon_fee,omitempty"`           //代金券金额
	RefundCount        int    `map:"refund_count,omitempty"`         //退款笔数

	OutRefundNo   map[string]string `map:"out_refund_no_,omitempty"`  //商户退款单号	商户系统内部的退款单号，商户系统内部唯一，只能是数字、大小写字母_-|*@ ，同一退款单号多次请求只退一笔。
	RefundId      map[string]string `map:"refund_id_,omitempty"`      //微信退款单号
	RefundChannel map[string]string `map:"refund_channel_,omitempty"` //ORIGINAL—原路退款 BALANCE—退回到余额 OTHER_BALANCE—原账户异常退到其他余额账户 OTHER_BANKCARD—原银行卡异常退到其他银行卡

	RefundFee           map[string]string `map:"refund_fee_,omitempty"`            //申请退款金额，订单金额，单位为分
	SettlementRefundFee map[string]string `map:"settlement_refund_fee_,omitempty"` //退款金额=申请退款金额-非充值代金券退款金额，退款金额<=申请退款金额

	CouponType         map[string]string `map:"coupon_type_,omitempty"`         //代金券类型
	CouponRefundFeeMap map[string]string `map:"coupon_refund_fee_,omitempty"`   //总代金券退款金额 & 单个代金券退款金额
	CouponRefundCount  map[string]string `map:"coupon_refund_count_,omitempty"` //退款代金券使用数量
	CouponRefundId     map[string]string `map:"coupon_refund_id_,omitempty"`    //退款代金券ID
	RefundStatus       map[string]string `map:"refund_status_,omitempty"`       //退款状态 SUCCESS—退款成功 REFUNDCLOSE—退款关闭。PROCESSING—退款处理中 CHANGE—退款异常，退款到银行发现用户的卡作废或者冻结了，导致原路退款银行卡失败，可前往商户平台（pay.weixin.qq.com）-交易中心，手动处理此笔退款。$n为下标，从0开始编号。
	RefundAccount      map[string]string `map:"refund_account_,omitempty"`      //退款资金来源
	RefundRecvAccout   map[string]string `map:"refund_recv_accout_,omitempty"`  //退款入账账户
	RefundSuccessTime  map[string]string `map:"refund_success_time_,omitempty"` //退款成功时间
}

// RefundQuery 查询退款 无需证书 https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_5
func (c *Client) RefundQuery(ctx context.Context, request *RefundQueryRequest) (resp *RefundQueryResponse, err error) {
	resp = &RefundQueryResponse{}
	_, err = c.Do(ctx, MasterApiBaseUrl, request, resp)
	return resp, err
}
