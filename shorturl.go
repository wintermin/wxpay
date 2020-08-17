package wxpay

import (
	"context"
)

type ShortUrlRequest struct {
	LongUrl string `map:"long_url"` //需要转换的URL，签名用原串，传输需URLencode
}

func (req ShortUrlRequest) Api() string {
	return ApiShortUrl
}

func (req ShortUrlRequest) IgnoreKey() []string {
	return nil
}

func (req ShortUrlRequest) IsNeedCert() bool {
	return false
}

type ShortUrlResponse struct {
	Response
	ShortUrl string `map:"short_url,omitempty"` //转换后的URL
}

// ShortUrl 无需证书 该接口主要用于Native支付模式一中的二维码链接转成短链接(weixin://wxpay/s/XXXXXX)，减小二维码数据量，提升扫描速度和精确度。
// https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_9&index=10
// 只针对此模式生效 https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=6_4
func (c *Client) ShortUrl(ctx context.Context, request *ShortUrlRequest) (resp *ShortUrlResponse, err error) {
	resp = &ShortUrlResponse{}
	_, err = c.Do(ctx, MasterApiBaseUrl, request, resp)
	return resp, err
}
