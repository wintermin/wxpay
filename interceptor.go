package wxpay

import (
	"context"
	"net/http"
	"time"
)

type ContextParam struct {
	Api          string
	StartTime    time.Time
	Request      *http.Request
	Response     *http.Response
	Error        error
	RequestBody  []byte
	ResponseBody []byte
	IsNeedRetry  bool
}

type Interceptor interface {
	ApiList() []string
	Before(ctx context.Context, param *ContextParam) error
	After(ctx context.Context, param *ContextParam) error
}

type ApiInterceptor map[string][]Interceptor

func (ai ApiInterceptor) Add(i Interceptor) {
	for _, api := range i.ApiList() {
		ai[api] = append(ai[api], i)
	}
}

func (ai ApiInterceptor) Before(ctx context.Context, param *ContextParam) error {
	for _, it := range ai[param.Api] {
		if err := it.Before(ctx, param); err != nil {
			return err
		}
	}
	return param.Error
}

func (ai ApiInterceptor) After(ctx context.Context, param *ContextParam) error {
	for _, it := range ai[param.Api] {
		if err := it.After(ctx, param); err != nil {
			return err
		}
	}
	return param.Error
}
