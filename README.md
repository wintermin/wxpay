# WECHAT PAY

微信支付SDK

---

### 支持的功能
- [x] 统一下单(UnifiedOrder)

   -  [x] APP(SDK支付)
   -  [x] MWEB(H5支付)
   -  [x] JSAPI(公众号支付|小程序支付)
   -  [x] NATIVE(扫码支付)
- [x] 查询订单(OrderQuery) 
- [x] 支付结果通知(Notify)

---

### 快速开始

```go
    c:=NewClient( /*your appId*/, /*your mchId*/, /*your key*/, /*Md5|HmacSha256*/, http.DefaultClient)
    //SDK不会定制相关业务代码，而是开放业务扩展的能力
    //例如：debug请求响应相关信息、耗时统计，请自行实现Interceptor，参考testLog
    c.ApiInterceptor().Add(&testLog{})
    //查询订单
    c.OrderQuery(context.Background(), &OrderQueryRequest{
    		OutTradeNo: "200805191943383785",
    	})
```
 
