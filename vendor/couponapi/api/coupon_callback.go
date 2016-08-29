package api

import (
	"github.com/gin-gonic/gin"
	. "couponapi/constant"
	"couponapi/dao"
	"errors"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"gitlab.qiyunxin.com/tangtao/utils/log"
)

type CouponCallbackDto struct  {
	//优惠代号
	CouponCode string `json:"coupon_code"`
	//追踪代号
	TrackCode string `json:"track_code"`
	OpenId string `json:"open_id"`
	//子交易号
	SubTradeNo string `json:"sub_trade_no"`

}

func CouponUseCallback(c *gin.Context)  {

	var resultDto *CouponCallbackDto
	err :=c.BindJSON(&resultDto)
	if err!=nil{
		log.Error(err)
		ResponseError400(c.Writer,10003)
		return
	}
	couponTrack := dao.NewCouponTrack()
	couponTrack,err = couponTrack.WithTrackCode(resultDto.TrackCode)
	if err!=nil{
		log.Error(err)
		ResponseError400(c.Writer,10001)
		return
	}
	if couponTrack==nil{
		log.Error("没有找到对应的coupon track")
		ResponseError400(c.Writer,10010)
		return
	}
	if couponTrack.Status != COUPON_STACK_STATUS_WAIT_USE {
		ResponseError400(c.Writer,10011)
		return
	}

	err =CouponInfoUpdate(couponTrack,resultDto)
	if err!=nil{
		log.Error(err)
		ResponseError400(c.Writer,10012)
		return
	}
	util.ResponseSuccess(c.Writer)
}

//更新优惠券信息
func CouponInfoUpdate(couponTrack *dao.CouponTrack,callbackDto *CouponCallbackDto)  error {

	couponuser,err :=dao.NewCouponUser().WithCouponCode(couponTrack.OpenId,couponTrack.CouponCode,couponTrack.AppId)
	if err!=nil{
		return err
	}

	couponbalance :=couponuser.Balance - couponTrack.Amount
	if couponbalance < 0 {

		return errors.New("优惠券余额不足")
	}
	tx,_ :=db.NewSession().Begin()
	defer func() {
		if err :=recover();err!=nil{
			tx.Rollback()
		}
	}()
	err =dao.NewCouponUser().UpdateBalanceWithCouponCodeTx(couponbalance,couponTrack.CouponCode,couponTrack.AppId,tx)
	if err!=nil{
		tx.Rollback()
		return err
	}
	err = couponTrack.UpdateTradeAndStatusWithTrackCodeTx(COUPON_STACK_STATUS_USED,0,callbackDto.SubTradeNo,couponTrack.TrackCode,tx)
	if err!=nil{
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}