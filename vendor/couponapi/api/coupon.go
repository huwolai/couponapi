package api

import (
	"github.com/gin-gonic/gin"
	"couponapi/service"
	"gitlab.qiyunxin.com/tangtao/utils/security"
	. "couponapi/constant"
	"net/http"
	"strings"
	"couponapi/dao"
	"gitlab.qiyunxin.com/tangtao/utils/log"
)

type CouponUserDto struct  {
	Id int64 `json:"id"`
	OpenId string `json:"open_id"`
	CouponCode string `json:"coupon_code"`
	Title string `json:"title"`
	Remark string `json:"remark"`
	Amount float64 `json:"amount"`
	Balance float64 `json:"balance"`
	CouponToken string `json:"coupon_token"`
	IsOne int `json:"is_one"`
	UseStatus int `json:"use_status"`
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
		"amount":int64(amount*1000),
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
	couponusers,err := service.CouponDistribute(openId,orderNo,flag,couponcodeArray,appId)
	if err!=nil{
		ResponseError400(c.Writer,10001)
		return
	}


	if couponusers!=nil{

	}
}

func CouponUserToDto(model *dao.CouponUser) *CouponUserDto  {
	dto := &CouponUserDto{}
	dto.Id = model.Id
	dto.Amount = model.Amount
	dto.Balance = model.Balance
	dto.CouponCode = model.CouponCode
	dto.IsOne = model.IsOne
	dto.OpenId = model.OpenId
	dto.Title = model.Title
	dto.Remark = model.Remark
	dto.UseStatus = model.UseStatus

	return dto
}