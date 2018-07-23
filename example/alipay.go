package example

// 使用示范文件
import (
	"github.com/worldlove/alipay"
	"io/ioutil"
	"log"
	"net/url"
)

func init() {
	initAliClient()
	initBaseParams()
}

// 初始化Alipay客户端
func initAliClient() {
	const priPath = "secret/private.pem"
	const pubPath = "secret/public.pem"
	priKey, err1 := ioutil.ReadFile(priPath)
	pubKey, err2 := ioutil.ReadFile(pubPath)
	if err1 != nil || err2 != nil {
		log.Println(err1, err2)
		panic("RSAKey Read Error")
	}
	payClient := alipay.NewAliClient(
		// 沙箱
		// AliTestAppID
		// AliTestGateway
		AliAppID,
		AliGateway,
		AliSignType,
		priKey,
		pubKey,
		// AliTestPubkey
		AliPubKey,
	)
	alipay.SetAlipay(payClient)
}

// 设置公用默认基础参数
func initBaseParams() {
	setDefaultParam()
	setDefaultPagePayParam()
	setDefaultWapPayParam()
}

// 设置基础公用参数
func setDefaultParam() {
	var defaultSysBase = make(url.Values)
	defaultSysBase.Add("sign_type", AliSignType)
	defaultSysBase.Add("format", AliFormat)
	defaultSysBase.Add("charset", AliCharset)
	defaultSysBase.Add("version", AliVersion)
	var defaultParam = alipay.AliParam{
		SysBase: defaultSysBase,
	}
	alipay.SetDefaultParams("default", &defaultParam)
}

// 设置电脑网站支付基础参数
func setDefaultPagePayParam() {
	var base = alipay.GetDefaultParams("default")
	//设置电脑网页支付特有公用参数
	base.SysBase.Set("return_url", AliPagePayReturn)
	base.SysBase.Set("notify_url", AliPagePayNotify)

	//设置电脑网页支付默认业务字段
	base.BizContent["timeout_express"] = "2d"
	base.BizContent["product_code"] = "FAST_INSTANT_TRADE_PAY"

	alipay.SetDefaultParams(alipay.AliPagePayMethod, base)
}

// 设置Wap网站支付基础参数
func setDefaultWapPayParam() {
	var base = alipay.GetDefaultParams("default")
	//设置手机网页支付特有公用参数
	base.SysBase.Set("return_url", AliPagePayReturn)
	base.SysBase.Set("notify_url", AliPagePayNotify)

	//设置手机网页支付默认业务字段
	base.BizContent["timeout_express"] = "2d"
	base.BizContent["product_code"] = "QUICK_WAP_WAY"
	alipay.SetDefaultParams(alipay.AliWapPayMethod, base)
}

func setDefaultTransToAccountParam() {
	var base = alipay.GetDefaultParams("default")
	//设置转账支付公用基础参数

	//设置转账支付公用业务参数
	base.BizContent["payee_type"] = "ALIPAY_LOGINID"
	base.BizContent["remark"] = "业务转账"
	alipay.SetDefaultParams(alipay.AliTransToAccountMethod, base)
}
