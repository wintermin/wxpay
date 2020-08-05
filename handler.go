package wxpay

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"time"
)

func (c *Client) Do(ctx context.Context, baseUrl string, req Request, resp interface{}) (bool, error) {
	ctxParam := &ContextParam{
		Api:         req.Api(),
		StartTime:   time.Now(),
		IsNeedRetry: false,
	}
	val := c.Encode(req)
	ctxParam.RequestBody = val.Encode()
	buff := &bytes.Buffer{}
	buff.Write(ctxParam.RequestBody)
	request, reqErr := http.NewRequestWithContext(ctx, http.MethodPost, baseUrl+req.Api(), buff)
	ctxParam.Error = reqErr
	err := c.apiInterceptor.Before(ctx, ctxParam)
	if err != nil {
		return ctxParam.IsNeedRetry, err
	}
	request.Header.Set("Content-Type", "text/xml")
	response, respErr := c.httpClient.Do(request)
	if respErr != nil {
		//请求渠道失败才重试，用于实现跨城容灾
		ctxParam.IsNeedRetry = true
	}
	ctxParam.Error = respErr
	ctxParam.Response = response
	if response != nil && response.Body != nil {
		ctxParam.ResponseBody, err = ioutil.ReadAll(response.Body)
		ctxParam.Error = err
		response.Body.Close()
		if err == nil {
			ctxParam.Error = c.Decode(resp, ctxParam.ResponseBody)
		}
	}
	err = c.apiInterceptor.After(ctx, ctxParam)
	return ctxParam.IsNeedRetry, err
}
