package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"gitlab.qiyunxin.com/tangtao/utils/startup"
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"gitlab.qiyunxin.com/tangtao/utils/queue"
	"github.com/streadway/amqp"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"couponapi/service"
	"couponapi/api"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, app_id, open_id")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT,DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {

	if os.Getenv("GO_ENV")=="" {
		os.Setenv("GO_ENV","tests")
		os.Setenv("APP_ID","couponapi")
	}
	err :=config.Init(false)
	util.CheckErr(err)
	err = startup.InitDBData()
	env := os.Getenv("GO_ENV")
	if env=="tests" {
		gin.SetMode(gin.DebugMode)
	}else if env== "production" {
		gin.SetMode(gin.ReleaseMode)
	}else if env == "preproduction" {
		gin.SetMode(gin.TestMode)
	}
	router := gin.Default()
	router.Use(CORSMiddleware())

	queue.SetupAMQP(config.GetValue("amqp_url").ToString())

	queue.ConsumeAccountEvent("couponapi_account_consumer", func(accountEvent *queue.AccountEvent, dv amqp.Delivery) {
		log.Error("获取到账户金额变化时间",accountEvent.EventKey)
		//账户充值
		if accountEvent.EventKey==queue.ACCOUNT_AMOUNT_EVENT_CHANGE {
			if accountEvent.Content!=nil&&accountEvent.Content.Action=="ACCOUNT_RECHARGE"{
				err :=service.RechargeCoupon(accountEvent.Content.OpenId,accountEvent.Content.SubTradeNo,float64(accountEvent.Content.ChangeAmount)/100,accountEvent.Content.AppId)
				if err!=nil{
					log.Error(err)
				}
				dv.Ack(false)
			}
		}
	})

	v1 := router.Group("/v1")
	{
		coupons :=v1.Group("/coupons")
		{
			coupons.POST("/callback",api.CouponUseCallback)
		}
		coupon :=v1.Group("/coupon")
		{	//获取用户优惠券总金额
			coupon.GET("/:open_id/amount",api.CouponAmount)
			//获取用户的优惠凭证
			coupon.GET("/:open_id/order/:order_no/tokens",api.CouponDistribute)
		}

	}
	router.Run(":8080")
}
