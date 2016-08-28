package service

import (
	"couponapi/dao"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"fmt"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"gitlab.qiyunxin.com/tangtao/utils/log"
)

func RechargeCoupon(openId string,subTradeNo string,amount float64,appId string) error  {

	couponUser :=dao.NewCouponUser()
	couponUser.Amount = amount
	couponUser.OpenId = openId
	couponUser.CouponCode = util.GenerUUId()
	couponUser.Balance = amount
	couponUser.IsOne = 0
	couponUser.Title="冲多少送多少活动"
	couponUser.Remark=fmt.Sprintf("冲%f送%f",amount,amount)
	couponUser.UseStatus = 1
	couponUser.AppId = appId
	tx,_ :=db.NewSession().Begin()
	err :=couponUser.InsertTx(tx)
	if err!=nil{
		log.Error(err)
		tx.Rollback()
		return err
	}

	couponTrack := dao.NewCouponTrack()
	couponTrack.OpenId = openId
	couponTrack.Title = couponUser.Title
	couponTrack.Remark = couponUser.Remark
	couponTrack.Amount = amount
	couponTrack.TradeNo =subTradeNo
	couponTrack.TradeType = 1 //充值
	couponTrack.TrackCode = util.GenerUUId()
	couponTrack.TrackType = 2
	err =couponTrack.InsertTx(tx)
	if err!=nil{
		log.Error(err)
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}