package wxpay

import (
	"context"
	"fmt"
	"strings"
	"time"
)

//打印请求前后相关报文
type testLog struct{}

func (l testLog) ApiList() []string {
	return ApiList
}

func (l testLog) Before(ctx context.Context, param *ContextParam) error {
	buff := &strings.Builder{}
	buff.WriteString("requestPath:")
	buff.WriteString(param.Api)
	buff.WriteString("|")
	buff.WriteString("requestBody:")
	buff.WriteString(string(param.RequestBody))
	fmt.Println(buff)
	return param.Error
}

func (l testLog) After(ctx context.Context, param *ContextParam) error {
	buff := &strings.Builder{}
	buff.WriteString("requestPath:")
	buff.WriteString(param.Api)
	buff.WriteString("|")
	buff.WriteString("responseBody:")
	buff.WriteString(string(param.ResponseBody))
	buff.WriteString("|")
	buff.WriteString("elapsedTime:")
	buff.WriteString(time.Since(param.StartTime).String())
	fmt.Println(buff)
	return param.Error
}
