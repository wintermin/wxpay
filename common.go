package wxpay

import (
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	//签名类型
	Md5        SignType = "MD5"
	HmacSha256 SignType = "HMAC-SHA256"

	//通信标识
	ReturnCodeSuccess ReturnCode = "SUCCESS"
	ReturnCodeFail    ReturnCode = "FAIL"

	//交易类型
	JsApi  TradeType = "JSAPI"  //JSAPI支付(公众号支付)或小程序支付
	Native TradeType = "NATIVE" //Native支付
	App    TradeType = "APP"    //APP支付
	MWeb   TradeType = "MWEB"   //H5支付

	//交易状态
	TradeSuccess    TradeState = "SUCCESS"    //支付成功
	TradeRefund     TradeState = "REFUND"     //转入退款
	TradeNotPay     TradeState = "NOTPAY"     //未支付
	TradeClosed     TradeState = "CLOSED"     //已关闭
	TradeRevoked    TradeState = "REVOKED"    //已撤销（刷卡支付
	TradeUserPaying TradeState = "USERPAYING" //用户支付中
	TradePayError   TradeState = "PAYERROR"   //支付失败(其他原因，如银行返回失败)

	KeySign       = "sign"
	KeyPackage    = "package"
	KeySignType   = "sign_type"
	KeyAppId      = "appid"
	KeyMchId      = "mch_id"
	KeyNonceStr   = "nonce_str"
	KeyReturnCode = "return_code"

	//微信域名
	MasterApiBaseUrl = "https://api.mch.weixin.qq.com"
	SlaveApiBaseUrl  = "https://api2.mch.weixin.qq.com"

	//api列表
	ApiUnifiedOrder = "/pay/unifiedorder" //统一下单
	ApiOrderQuery   = "/pay/orderquery"   //查询订单

	// 响应渠道处理成功
	responseSuccessBody = "<xml><return_code><![CDATA[SUCCESS]]></return_code><return_msg><![CDATA[OK]]></return_msg></xml>"
)

var (
	ApiList = []string{ApiUnifiedOrder, ApiOrderQuery}
	//标准北京时间，时区为东八区；如果商户的系统时间为非标准北京时间。参数值必须根据商户系统所在时区先换算成标准北京时间
	BeijingLocation = time.FixedZone("Asia/Shanghai", 8*60*60)
)

type (
	SignType   string //签名类型
	ReturnCode string //通信标识
	TradeType  string //交易类型
	TradeState string //交易状态
	Request    interface {
		Api() string         //渠道接口
		IgnoreKey() []string //过滤不需要上送的Key
	}
	Response struct {
		ReturnCode ReturnCode `map:"return_code"`          //返回状态码 SUCCESS/FAIL 此字段是通信标识
		ReturnMsg  string     `map:"return_msg,omitempty"` //返回信息 返回信息，如非空，为错误原因 签名失败 参数格式校验错误
		//以下字段在return_code为SUCCESS的时候有返回
		AppId      string `map:"appid,omitempty"`        //微信开放平台审核通过的应用APPID
		MchId      string `map:"mch_id,omitempty"`       //微信支付分配的商户号
		ResultCode string `map:"result_code,omitempty"`  //业务结果
		ErrCode    string `map:"err_code,omitempty"`     //错误代码
		ErrCodeDes string `map:"err_code_des,omitempty"` //错误代码描述
	}
	Values map[string]string
)

// 是否通信成功
func (r Response) IsOk() bool {
	if r.ReturnCode == ReturnCodeSuccess {
		return true
	}
	return false
}

func (val Values) Add(k, v string) {
	val[k] = v
}

func (val Values) Get(k string) string {
	return val[k]
}

type xmlValuesEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func (val Values) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(val) == 0 {
		return nil
	}
	start.Name = xml.Name{Local: "xml"}
	err := e.EncodeToken(start)
	if err != nil {
		return err
	}
	for k, v := range val {
		e.Encode(xmlValuesEntry{XMLName: xml.Name{Local: k}, Value: v})
	}
	return e.EncodeToken(start.End())
}

func (val Values) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		e := &xmlValuesEntry{}
		err := d.Decode(e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		(val)[e.XMLName.Local] = e.Value
	}
	return nil
}

func (val Values) Encode() []byte {
	body, _ := xml.Marshal(val)
	return body
}

func (val Values) Decode(body []byte) {
	xml.Unmarshal(body, &val)
}

func ToString(object interface{}) string {
	if object == nil {
		return ""
	}
	switch v := object.(type) {
	case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%g", v)
	case string:
		return v
	case time.Time:
		return FormatTime(object.(time.Time))
	default:
		return fmt.Sprintf("%s", v)
	}
}

func (val Values) ForStruct(object interface{}) {
	elem := reflect.ValueOf(object).Elem()
	relType := elem.Type()
	for i := 0; i < relType.NumField(); i++ {
		f := relType.Field(i)
		name := f.Tag.Get("map")
		if name == "-" {
			continue
		}
		value := ""
		if elem.Field(i).CanInterface() {
			value = ToString(elem.Field(i).Interface())
		}
		if name == "" {
			name = f.Name
		} else {
			arr := strings.Split(name, ",")
			if len(arr) > 1 && arr[1] == "omitempty" && value == "" {
				continue
			}
			name = arr[0]
		}
		val.Add(name, value)
	}
}

func (val Values) ToStruct(object interface{}) {
	elem := reflect.ValueOf(object).Elem()
	relType := elem.Type()
	for i := 0; i < relType.NumField(); i++ {
		f := relType.Field(i)
		key := getKey(f)
		if key == "" {
			continue
		}
		if f.Type.Kind() == reflect.Struct {
			// 暂时只支持一层 嵌套struct赋值
			for j := 0; j < f.Type.NumField(); j++ {
				key = getKey(f.Type.Field(j))
				value := elem.FieldByIndex([]int{i, j})
				if vv, ok := val[key]; ok {
					valueSet(value, vv)
				}
			}
		} else if f.Type.Kind() == reflect.Map {
			mv := elem.Field(i)
			if mv.IsNil() { //为nil,反射初始化map
				mv.Set(reflect.MakeMap(f.Type))
			}
			for k, v := range val {
				//key模糊匹配
				if strings.Contains(k, key) {
					mv.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
				}
			}
		} else if vv, ok := val[key]; ok {
			valueSet(elem.Field(i), vv)
		}
	}
}

func getKey(rsf reflect.StructField) string {
	name := rsf.Tag.Get("map")
	if name == "-" {
		return ""
	}
	if name == "" {
		name = rsf.Name
	} else {
		arr := strings.Split(name, ",")
		name = arr[0]
	}
	return name
}

// 暂时只支持基本数据类型
func valueSet(value reflect.Value, v string) {
	if value.CanSet() {
		switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			vv, _ := strconv.Atoi(v)
			value.SetInt(int64(vv))
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			vv, _ := strconv.Atoi(v)
			value.SetUint(uint64(vv))
		case reflect.Float32:
			vv, _ := strconv.ParseFloat(v, 32)
			value.SetFloat(vv)
		case reflect.Float64:
			vv, _ := strconv.ParseFloat(v, 64)
			value.SetFloat(vv)
		case reflect.Bool:
			vv, _ := strconv.ParseBool(v)
			value.SetBool(vv)
		case reflect.String:
			value.SetString(v)
		}
	}
}

// 响应渠道处理成功
func ResponseSuccess() string {
	return responseSuccessBody
}

// ResponseFail 收到回调后响应渠道的报文 处理失败
// msg 失败原因
func ResponseFail(msg string) string {
	return "<xml><return_code><![CDATA[FAIL]]></return_code><return_msg><![CDATA[" + msg + "]]></return_msg></xml>"
}

// FormatTime 北京时间yyyyMMddHHmmss
func FormatTime(t time.Time) string {
	return t.In(BeijingLocation).Format("20060102150405")
}

// ParseTime 北京时间yyyyMMddHHmmss
func ParseTime(value string) (time.Time, error) {
	return time.ParseInLocation("20060102150405", value, BeijingLocation)
}
