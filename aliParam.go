package alipay

import (
	"net/url"
)

// Aliparam 支付宝参数结构体 包含基础参数参数与业务参数
// 基础参数key=value&  业务参数json字符串
type AliParam struct {
	SysBase    url.Values
	BizContent map[string]string
}

var defaultAliParam = make(map[string]*AliParam)

func GetDefaultParams(key string) *AliParam {
	// copy公共参数 必须深copy，不然多次请求共用同一数据实体，公共参数修改，数据错乱
	if params, ok := defaultAliParam[key]; ok {
		return copyAliParam(params)
	}
	return copyAliParam(defaultAliParam["default"])
}

func SetDefaultParams(key string, param *AliParam) {
	defaultAliParam[key] = param
}

func copyAliParam(src *AliParam) *AliParam {
	var dst = AliParam{
		SysBase:    copyUrlValues(src.SysBase),
		BizContent: copyMap(src.BizContent),
	}
	return &dst
}
