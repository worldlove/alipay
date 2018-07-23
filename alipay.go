package alipay

import (
	"encoding/json"
	"log"
	"time"
)

// 支付宝SDK使用之前需初始化ALipay客户端
// Alipay 一般一个业务只需要一个支付宝客户端, 不需要每次使用都初始化
var Alipay *AliClient

func SetAlipay(ali *AliClient) {
	Alipay = ali
}

// QueryBase 查询生成函数，基于method生成不同的查询函数
func QueryBase(method string) func(map[string]string) (Response, error) {
	return func(bizContent map[string]string) (Response, error) {
		var base = GetDefaultParams("default")
		base.SysBase.Set("timestamp", time.Now().Format("2006-01-02 15:04:05"))
		base.SysBase.Set("method", method)
		biz_content, _ := json.Marshal(bizContent)
		base.SysBase.Set("biz_content", string(biz_content))

		var res = make(Response)
		err := Alipay.SignAndRequest(base.SysBase, "POST", res)
		if err != nil {
			log.Println("Func:QueryPagePay Time:1 Error:", err)
			if err := Alipay.SignAndRequest(base.SysBase, "POST", res); err != nil {
				log.Println("Func:QueryPagePay Time:2 Error:", err)
				return nil, err
			}
		}
		return res, nil
	}
}

// 基本支付函数, 封装各种参数生成URL以备前端调用
func NewTradeToURL(param *AliParam) string {
	param.SysBase.Set("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	biz_content, _ := json.Marshal(param.BizContent)
	param.SysBase.Add("biz_content", string(biz_content))
	return Alipay.SignToUrl(param.SysBase)
}
