package alipay

// 通用查询接口，查询用户支付是否成功（电脑支付和wap支付都可以使用）
// func(bizContent map[string]string) (Response, error)
var QueryTrade = QueryBase(AliQueryMethod)

// 转账查询接口， 查询转账是否成功
// func(bizContent map[string]string) (Response, error)
var QueryTransToAccount = QueryBase(AliTransOrderQueryMethod)
