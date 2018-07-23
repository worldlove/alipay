# alipay
支付宝SDK（部分功能） go语言实现
支付宝api文档:[文档地址](https://docs.open.alipay.com/api)

## 安装
```bash
go get github.com/worldlove/alipay
```

## 开发理念
1. 简单
   - 使用url.Values 和 map[string]string 表达数据结构
   - 付宝重新设计之后，多数参数都是通用的，没必要为每个接口创建一个struct
   - 直接可以使用ulr.Values.Encode()自动排序和编码

2. 高效
   - 复用全局Alipay 客户端对象(必须在调用接口之前先初始化全局Alipay对象)
   - 增加异步调用逻辑，可以将转账、查询等操作发送到异步调用链

## 使用示例
```go

// 初始化alipay客户端
import (
	"github.com/worldlove/alipay"
	"io/ioutil"
	"log"
)

func init() {
	initAliClient()
    setDefaultParam()
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

```
## 实现API

- [x] [电脑网页支付 API](#电脑网页支付-api)
- [x] [手机WAP支付 API](#手机wap支付-api)
- [x] [统一支付查询 API](#统一支付查询)
- [x] [账号转账 API](#账号转账-api)
- [x] [转账查询 API](#转账查询-api)

### 电脑网页支付宝 API
```go
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

var content = map[string]string{
    "out_trade_no":    "123456789",
    "total_amount":    "1.11",
    "subject":         "名称",
    "body":            "描述",
    "passback_params": "原样返回",
}

payURL := alipay.DefaultPagePay(content)

// 网页支付返回(参数为http.Request, 支持各种框架)
// 例：使用gin框架
func AlipayPageReturn(c *gin.Context) {
    values, err := alipay.AlipayReturn(c.Request)
}

// 各种支付异步通知(参数为http.Request, 支持各种框架)
// 例：使用gin框架
func AlipayNotify(c *gin.Context) {
    values, err := alipay.AlipayNotify(c.Request)
}

```
### 手机WAP支付 API
```go
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

var content = map[string]string{
    "out_trade_no":    "123456789",
    "total_amount":    "1.11",
    "subject":         "名称",
    "body":            "描述",
    "passback_params": "原样返回",
}

payURL := alipay.DefaultWapPay(content)
```
### 统一支付查询 API
```go
var payComment = map[string]string{
    "out_trade_no": "123456789",
}
res, err := alipay.QueryTrade(payComment)

// 将查询异步话
// // 处理返回结果
func queryTradeReturn(alipay.Response, error) {
....
}

var action = alipay.PayAction{
    Payload: payComment,
    Action:  alipay.QueryTrade,
    Return:  queryTradeReturn,
}

// 将相关事件插入异步链
alipay.PushPayActionToChan(action)
```

### 账号转账 API
```go
func setDefaultTransToAccountParam() {
	var base = alipay.GetDefaultParams("default")
	//设置转账支付公用基础参数

	//设置转账支付公用业务参数
	base.BizContent["payee_type"] = "ALIPAY_LOGINID"
	base.BizContent["remark"] = "业务转账"
	alipay.SetDefaultParams(alipay.AliTransToAccountMethod, base)
}
setDefaultTransToAccountParam()

var payComment = map[string]string{
    "out_biz_no":    "123456789",
    "payee_account": "example@account.com",
    "amount":        "1.88",
}
res, err := alipay.DefaultTransToAccount(payComment)

// 插入异步链
var action = alipay.PayAction{
    Payload: payComment,
    Action:  alipay.DefaultTransToAccount,
    Return:  proxyTransReturn,
}
alipay.PushPayActionToChan(&action)

```

### 转账查询 API
```go

var payComment = map[string]string{
    "out_trade_no": "123456789",
}
res, err := alipay.QueryTransToAccount(payComment)

```
