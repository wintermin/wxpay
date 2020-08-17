package wxpay

import (
	"context"
	"strconv"
	"time"
)

type UnifiedOrderRequest struct {
	DeviceInfo     string    `map:"device_info,omitempty"` //终端设备号(门店号或收银设备ID)，默认请传"WEB"
	Body           string    `map:"body"`                  //商品描述
	Detail         string    `map:"detail,omitempty"`      //商品详情
	Attach         string    `map:"attach,omitempty"`      //附加数据，在查询API和支付通知中原样返回，可作为自定义参数使用
	OutTradeNo     string    `map:"out_trade_no"`          //商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|* 且在同一个商户号下唯一
	FeeType        string    `map:"fee_type,omitempty"`    //标价币种
	TotalFee       int       `map:"total_fee"`             //订单总金额，单位为分
	SpBillCreateIp string    `map:"spbill_create_ip"`      //支持IPV4和IPV6两种格式的IP地址。用户的客户端IP
	TimeStart      time.Time `map:"time_start,omitempty"`  //订单生成时间，格式为yyyyMMddHHmmss
	TimeExpire     time.Time `map:"time_expire,omitempty"` //订单失效时间，格式为yyyyMMddHHmmss 建议：最短失效时间间隔大于1分钟
	GoodsTag       string    `map:"goods_tag,omitempty"`   //订单优惠标记，使用代金券或立减优惠功能时需要的参数
	NotifyUrl      string    `map:"notify_url"`            //异步接收微信支付结果通知的回调地址，通知url必须为外网可访问的url，不能携带参数。
	TradeType      TradeType `map:"trade_type"`            //交易类型
	ProductId      string    `map:"product_id,omitempty"`  //商品ID trade_type=NATIVE时，此参数必传。此参数为二维码中包含的商品ID，商户自行定义
	LimitPay       string    `map:"limit_pay,omitempty"`   //指定支付方式 上传此参数no_credit--可限制用户不能使用信用卡支付
	OpenId         string    `map:"openid,omitempty"`      //trade_type=JSAPI时（即JSAPI支付），此参数必传
	Receipt        string    `map:"receipt,omitempty"`     //电子发票入口开放标识 Y，传入Y时，支付成功消息和支付详情页将出现开票入口。需要在微信支付商户平台或微信公众平台开通电子发票功能，传此字段才可生效
	SceneInfo      string    `map:"scene_info,omitempty"`  //场景信息 该字段常用于线下活动时的场景信息上报，支持上报实际门店信息，商户也可以按需求自己上报相关信息。该字段为JSON对象数据，对象格式为{"store_info":{"id": "门店ID","name": "名称","area_code": "编码","address": "地址" }}
}

func (req UnifiedOrderRequest) Api() string {
	return ApiUnifiedOrder
}

func (req UnifiedOrderRequest) IgnoreKey() []string {
	return nil
}

func (req UnifiedOrderRequest) IsNeedCert() bool {
	return false
}

type UnifiedOrderResponse struct {
	Response
	TradeType   TradeType         `map:"trade_type,omitempty"` //调用接口提交的交易类型
	PrepayId    string            `map:"prepay_id,omitempty"`  //微信生成的预支付回话标识，用于后续接口调用中使用，该值有效期为2小时
	MWebUrl     string            `map:"mweb_url,omitempty"`   //H5支付：支付跳转链接 拉起微信支付收银台的中间页面，可通过访问该url来拉起微信客户端，完成支付,mweb_url的有效期为5分钟。
	CodeUrl     string            `map:"code_url,omitempty"`   //扫码支付：trade_type=NATIVE时有返回，此url用于生成支付二维码，然后提供给用户进行扫码支付
	ClientParam map[string]string `map:"-"`                    //JSAPI支付(公众号支付)或小程序支付 sdk支付时有效: 客户端发起请求相关参数，原样返回给客户端即可
}

// UnifiedOrder 统一下单 无需证书 https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_1
func (c *Client) UnifiedOrder(ctx context.Context, request *UnifiedOrderRequest) (resp *UnifiedOrderResponse, err error) {
	resp = &UnifiedOrderResponse{}
	_, err = c.Do(ctx, MasterApiBaseUrl, request, resp)
	if err != nil {
		return resp, err
	}
	if resp.TradeType == JsApi || resp.TradeType == App {
		val := make(Values)
		switch resp.TradeType {
		case JsApi: //公众号或小程序支付 https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=7_7&index=6
			val.Add("appId", c.appId)
			val.Add("nonceStr", c.NonceStr())
			val.Add("signType", string(c.signType))
			val.Add(KeyPackage, "prepay_id="+resp.PrepayId)
			val.Add("timeStamp", FormatTime(time.Now()))
			c.Sign(val)
			val.Add("paySign", val.Get(KeySign))
			delete(val, KeySign)
		case App: //sdk支付 https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_12&index=2
			val.Add(KeyAppId, c.appId)
			val.Add("partnerid", c.mchId)
			val.Add("noncestr", c.NonceStr())
			val.Add(KeyPackage, "Sign=WXPay")
			val.Add("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
			val.Add("prepayid", resp.PrepayId)
			c.Sign(val)
		}
		resp.ClientParam = val
	}
	return resp, err
}
