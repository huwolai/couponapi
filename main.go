package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"gitlab.qiyunxin.com/tangtao/utils/startup"
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"gitlab.qiyunxin.com/tangtao/utils/util"
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
		os.Setenv("APP_ID","shopapi")
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

	v1 := router.Group("/v1")
	{

	}


}
