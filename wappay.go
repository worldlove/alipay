package alipay

// wap网站支付，在前台打开组装好的url即可

// DefaultWapPay 使用前请先设置wappay 公用参数
// SetDefaultParams("wappay", *AliParam)
func DefaultWapPay(bizContent map[string]string) string {
	var base = GetDefaultParams(AliWapPayMethod)
	base.SysBase.Set("method", AliWapPayMethod)
	for key, value := range bizContent {
		base.BizContent[key] = value
	}
	return NewTradeToURL(base)
}
