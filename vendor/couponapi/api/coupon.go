package api

import (
	"github.com/gin-gonic/gin"
	"couponapi/service"
	"gitlab.qiyunxin.com/tangtao/utils/security"
	. "couponapi/constant"
	"net/http"
	"strings"
	"gitlab.qiyunxin.com/tangtao/utils/log"
)

type CouponUser struct  {
	//优惠代码
	CouponCode string `json:"coupon_code"`
	//优惠凭证
	CouponToken string `json:"coupon_token"`
	//优惠金额
	CouponAmount float64 `json:"coupon_amount"`
	TrackCode string `json:"track_code"`
	//订单号
	OrderNo string `json:"order_no"`
	OpenId string `json:"open_id"`
	//appid
	AppId string `json:"app_id"`
}

//用户优惠券总金额
func CouponAmount(c *gin.Context)  {

	openId :=c.Param("open_id")
	appId := security.GetAppId2(c.Request)

	amount,err :=service.CouponAmount(openId,appId)
	if err!=nil{
		log.Error(err)
		ResponseError400(c.Writer,10001)
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"amount":amount,
		"open_id":openId,
	})
}

//获取用户的优惠信息
func CouponDistribute(c *gin.Context)  {
	//用户ID
	openId :=c.Param("open_id")
	orderNo := c.Param("order_no")
	//优惠券标识
	flag :=c.Query("flag")
	//优惠券代号多个以逗号分开
	couponcodes :=c.Query("codes")
	var couponcodeArray []string
	if couponcodes!=""&&len(couponcodes)>0 {
		couponcodeArray = strings.Split(couponcodes,",")
	}
	appId :=security.GetAppId2(c.Request)
	log.Info("appId:",appId)
	log.Info("flag:",flag)
	couponuser,err := service.CouponDistribute(openId,orderNo,flag,couponcodeArray,appId)
	if err!=nil{
		log.Error(err)
		ResponseError400(c.Writer,10009)
		return
	}
	couponusers :=make([]*CouponUser,0)
	if couponuser!=nil{
		couponusers = append(couponusers,CouponUserToDto(couponuser))
	}
	c.JSON(http.StatusOK,couponusers)
}

func CouponUserToDto(model *service.CouponUser) *CouponUser  {
	dto := &CouponUser{}
	dto.CouponCode = model.CouponCode
	dto.OpenId = model.OpenId
	dto.TrackCode = model.TrackCode
	dto.AppId = model.AppId
	dto.CouponAmount = model.CouponAmount
	dto.CouponToken = model.CouponToken
	dto.OrderNo = model.OrderNo

	return dto
}