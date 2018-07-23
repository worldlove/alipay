package alipay

const (
	ServerURL = "https://example.com/"

	AliAppID   = ""
	AliPubKey  = ""
	AliGateway = "https://openapi.alipay.com/gateway.do"

	// 沙箱
	AliTestAppID   = ""
	AliTestPubKey  = ""
	AliTestGateway = "https://openapi.alipaydev.com/gateway.do"

	AliSignType = "RSA2"
	AliFormat   = "JSON"
	AliCharset  = "utf-8"
	AliVersion  = "1.0"

	AliPagePayMethod       = "alipay.trade.page.pay"
	AliPagePayReturnMethod = "alipay.trade.page.pay.return"
	AliQueryMethod         = "alipay.trade.query"

	AliWapPayMethod       = "alipay.trade.wap.pay"
	AliWapPayReturnMethod = "alipay.trade.wap.pay.return"

	AliTransToAccountMethod  = "alipay.fund.trans.toaccount.transfer"
	AliTransOrderQueryMethod = "alipay.fund.trans.order.query"
)

var (
	AliPagePayReturn = ServerURL + "/alipay/return"
	AliPagePayNotify = ServerURL + "/alipay/notify"
	AliTransRemark   = ServerURL + "转账"
)
