package api

import (
	"github.com/gin-gonic/gin"
	"couponapi/service"
	"gitlab.qiyunxin.com/tangtao/utils/security"
	. "couponapi/constant"
	"github.com/Azure/azure-sdk-for-go/core/http"
)

func CouponAmount(c *gin.Context)  {

	openId :=c.Param("open_id")
	appId := security.GetAppId2(c.Request)

	amount,err :=service.CouponAmount(openId,appId)
	if err!=nil{
		ResponseError400(c.Writer,10001)
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"amount":int64(amount*1000),
		"open_id":openId,
	})
}
