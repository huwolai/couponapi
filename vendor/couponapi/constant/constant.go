package constant

import (
	"net/http"
	"gitlab.qiyunxin.com/tangtao/utils/util"
)

const (
	//等待使用
	COUPON_STACK_STATUS_WAIT_USE =0
	//已使用
	COUPON_STACK_STATUS_USED = 1
)

var Q func(code int) string
var strMap map[int]string
func init()  {
	strMap = map[int]string{
		10001:"查询错误,请联系管理员",
		10002: "认证失败",
		10003: "参数数据格式有误!",
		10004: "用户ID不能为空!",
		10005: "交易号不能为空",
		10006: "查询出错!",
		10007: "没有找到账户记录!",
		10008: "查询用户账户信息出错!",
		10009: "下发优惠凭证失败",
		10010: "不存在优惠信息!",
		10011: "优惠券不是待使用状态",
		10012: "优惠券信息更新失败",

	}

	Q = GetCodeMsg
}

func GetCodeMsg(code int) string  {

	return strMap[code]
}

func ResponseError400(w http.ResponseWriter,statusCode int)  {

	util.ResponseErrorSS(w,http.StatusBadRequest,statusCode,Q(statusCode))
}

func ResponseError401(w http.ResponseWriter,statusCode int)  {
	util.ResponseErrorSS(w,http.StatusUnauthorized,statusCode,Q(statusCode))
}