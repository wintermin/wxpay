package wxpay

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"hash"
	"net/http"
	"sort"
)

// 微信支付sdk核心对象：提供编码、解码等核心能力
type Client struct {
	appId          string
	mchId          string
	key            string
	signType       SignType
	apiInterceptor ApiInterceptor
	httpClient     *http.Client
	certHttpClient *http.Client
	buildNonceStr  func() string
}

func NewClient(appId, mchId, key string, signType SignType, http *http.Client) *Client {
	return &Client{
		appId:          appId,
		mchId:          mchId,
		key:            key,
		signType:       signType,
		apiInterceptor: make(ApiInterceptor),
		httpClient:     http,
	}
}

// http 无需证书的api client certHttp 需要证书的api client
func NewCertClient(appId, mchId, key string, signType SignType, http, certHttp *http.Client) *Client {
	return &Client{
		appId:          appId,
		mchId:          mchId,
		key:            key,
		signType:       signType,
		apiInterceptor: make(ApiInterceptor),
		httpClient:     http,
		certHttpClient: certHttp,
	}
}

func (c *Client) AppId() string {
	return c.appId
}

func (c *Client) MchId() string {
	return c.mchId
}

func (c *Client) Key() string {
	return c.key
}

func (c *Client) SignType() SignType {
	return c.signType
}

func (c *Client) HttpClient() *http.Client {
	return c.httpClient
}

func (c *Client) CertHttpClient() *http.Client {
	return c.certHttpClient
}

func (c *Client) ApiInterceptor() ApiInterceptor {
	return c.apiInterceptor
}

func (c *Client) SetBuildNonceStr(f func() string) {
	c.buildNonceStr = f
}

// NonceStr 生成随机字符串、通过SetBuildNonceStr可以自定义生成方法
func (c *Client) NonceStr() string {
	if c.buildNonceStr == nil {
		h := md5.New()
		h.Write([]byte(RandStringRunes(16)))
		return hex.EncodeToString(h.Sum(nil))
	}
	return c.buildNonceStr()
}

// Encode 上送渠道报文编码、签名
func (c *Client) Encode(req Request) Values {
	val := make(Values)
	val.Add(KeyAppId, c.appId)
	val.Add(KeyMchId, c.MchId())
	val.Add(KeySignType, string(c.signType))
	val.Add(KeyNonceStr, c.NonceStr())
	val.ForStruct(req)

	ignoreKey := req.IgnoreKey()
	if ignoreKey != nil {
		for _, k := range ignoreKey {
			delete(val, k)
		}
	}

	c.Sign(val)
	return val
}

// Decode 渠道响应报文解码、验签、对象转换赋值
func (c *Client) Decode(resp interface{}, body []byte) error {
	val := make(Values)
	val.Decode(body)
	returnCode := val.Get(KeyReturnCode)
	if returnCode == "SUCCESS" {
		sign := val.Get(KeySign)
		c.Sign(val)
		newSign := val.Get(KeySign)
		if sign != newSign {
			return errors.New("signature verification failed")
		}
	}
	val.ToStruct(resp)
	return nil
}

func (c *Client) Sign(val Values) {
	var h hash.Hash
	switch c.signType {
	case HmacSha256:
		h = sha256.New()
	default:
		h = md5.New()
	}
	keys := make([]string, 0, len(val))
	for k := range val {
		if k == KeySign {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	buff := bufio.NewWriterSize(h, 128)
	for _, k := range keys {
		v := val[k]
		if len(v) == 0 {
			continue
		}
		buff.WriteString(k)
		buff.WriteByte('=')
		buff.WriteString(v)
		buff.WriteByte('&')
	}
	buff.WriteString("key=")
	buff.WriteString(c.key)
	buff.Flush()

	signature := make([]byte, hex.EncodedLen(h.Size()))
	hex.Encode(signature, h.Sum(nil))
	val.Add(KeySign, string(bytes.ToUpper(signature)))
}
