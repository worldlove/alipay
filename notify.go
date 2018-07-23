package alipay

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// 此处使用通用http.Request作为请求参数，适配各种框架
// 比如使用gin框架：
// func AlipayReturn(c *gin.Context) {
//     query, err := alipay.AliPayReturn(c.Request)
//     ...
// }

// 支付宝异步通知（对应notify_url）
func AliPayNotify(req *http.Request) (url.Values, error) {
	if err := req.ParseForm(); err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}
	var query = req.PostForm
	if err := Alipay.VerifyResponseURL(query); err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}
	var status = query.Get("trade_status")
	if status != "TRADE_FINISHED" && status != "TRADE_SUCCESS" {
		return query, fmt.Errorf("AlipayNotify Error:%v", status)
	}
	return query, nil
}

// 支付宝支付完成后页面跳转返回（对应return_url）
// 注意返回后会丢失cookie信息，需要做特殊处理重新设置cookie（比如把cookie放入passback_params字段），不然会掉线
func AliPayReturn(req *http.Request) (url.Values, error) {
	var query = req.URL.Query()
	if err := Alipay.VerifyResponseURL(query); err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}
	// 注意返回后会丢失cookie信息，需要做特殊处理重新设置cookie,
	// 比如把cookieID(避免把真实cookie写入，产生cookie劫持风险)放入passback_params字段，
	// 然后收到返回信息时重新设置cookie, 不然会掉线.
	// cookieID := query.Get("passback_params") //然后通过cookieID重新设置cookie
	return query, nil
}
