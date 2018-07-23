package alipay

import (
	"encoding/json"
	"errors"
	"log"
	"time"
)

func DefaultTransToAccount(bizContent map[string]string) (Response, error) {
	var base = GetDefaultParams(AliTransToAccountMethod)
	base.SysBase.Set("method", AliTransToAccountMethod)
	for key, value := range bizContent {
		base.BizContent[key] = value
	}
	return NewTransToAccount(base)
}

func NewTransToAccount(param *AliParam) (Response, error) {
	param.SysBase.Set("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	biz_content, _ := json.Marshal(param.BizContent)
	param.SysBase.Set("biz_content", string(biz_content))

	var res = make(Response)

	err := Alipay.SignAndRequest(param.SysBase, "POST", res)
	if err != nil {
		log.Println("Func:NewTransToAccount Time:1 Error:", err)
		if err := Alipay.SignAndRequest(param.SysBase, "GET", res); err != nil {
			log.Println("Func:NewTransToAccount Time:2 Error:", err)
			return nil, err
		}
	}
	var resCode = res["code"]
	if resCode != "10000" {
		var resSubCode = res["sub_code"]
		if resCode == "20000" || resCode == "40004" || resSubCode == "SYSTEM_ERROR" {
			var queryContent = map[string]string{
				"out_biz_no": param.BizContent["out_biz_no"],
			}
			if res["order_id"] != "" {
				queryContent["order_id"] = res["order_id"]
			}
			query, err := QueryTransToAccount(queryContent)
			if err != nil {
				return nil, err
			}
			if query["code"] == "10000" && query["error_code"] == "ORDER_NOT_EXIST" {
				if err := Alipay.SignAndRequest(param.SysBase, "GET", res); err != nil {
					log.Println("Func:NewTransToAccount Time:3 Error:", err)
					return nil, err
				}
				if resCode != "10000" {
					return res, errors.New(res["msg"])
				}
				return res, nil
			}
		}
		return res, errors.New(res["msg"])
	}
	return res, nil
}
