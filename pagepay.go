package alipay

// 电脑网站支付，在前台打开组装好的url即可

// DefaultPagePay 使用前请先设置pagepay 公用参数
// SetDefaultParams("pagepay", *AliParam)
func DefaultPagePay(bizContent map[string]string) string {
	var base = GetDefaultParams(AliPagePayMethod)
	base.SysBase.Set("method", AliPagePayMethod)
	for key, value := range bizContent {
		base.BizContent[key] = value
	}
	return NewTradeToURL(base)
}
