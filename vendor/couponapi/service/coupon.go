package service

import (
	"couponapi/dao"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"fmt"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"couponapi/comm"
	"errors"
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"couponapi/constant"
)

type CouponUser struct  {

	Title string
	Remark string
	//优惠代码
	CouponCode string
	//优惠凭证
	CouponToken string
	//优惠金额
	CouponAmount float64
	TrackCode string
	//订单号
	OrderNo string
	OpenId string
	//appid
	AppId string
}

const FLAG_ACCOUNT_RECHARGE  = "ACCOUNT_RECHARGE"

//充值获取优惠券
func RechargeCoupon(openId string,subTradeNo string,amount float64,appId string) error  {
	couponUser :=dao.NewCouponUser()
	// 获取用户优惠券信息
	couponUsers ,err :=couponUser.WithCodesOrFlag(openId,nil,FLAG_ACCOUNT_RECHARGE,appId)
	if err!=nil{
		return err
	}

	//如果存在此优惠券信息 那么将钱充值到此券中
	if couponUsers!=nil&&len(couponUsers)>0{
		cpuser :=couponUsers[0]
		err = couponUser.UpdateAmountAndBalanceWithId(cpuser.Amount+amount,cpuser.Balance+amount,cpuser.Id)
		return err
	}

	couponUser.Amount = amount
	couponUser.OpenId = openId
	couponUser.CouponCode = util.GenerUUId()
	couponUser.Balance = amount
	couponUser.IsOne = 0
	couponUser.Title="冲多少送多少活动"
	couponUser.Remark=fmt.Sprintf("冲%.2f送%.2f",amount,amount)
	couponUser.UseStatus = 1
	couponUser.Flag = FLAG_ACCOUNT_RECHARGE
	couponUser.AppId = appId
	tx,_ :=db.NewSession().Begin()
	defer func() {
		if err := recover();err != nil {
			log.Error(err)
			tx.Rollback()
		}
	}()
	err =couponUser.InsertTx(tx)
	if err!=nil{
		log.Error(err)
		tx.Rollback()
		return err
	}

	couponTrack := dao.NewCouponTrack()
	couponTrack.CouponCode = couponUser.CouponCode
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
	err =tx.Commit()
	if err!=nil{
		log.Error(err)
		return err
	}
	return nil
}

//获取优惠券金额
func CouponAmount(openId string,appId string) (float64,error)  {

	return dao.NewCouponUser().TotalAmountWithOpenId(constant.COUPON_USER_STATUS_ACTIVED,openId,appId)
}

//分发优惠券
func CouponDistribute(openId string,orderNo string,flag string,codes []string,appId string) (*CouponUser,error) {

	couponuser :=dao.NewCouponUser()
	couponuserList,err :=couponuser.WithCodesOrFlag(openId,codes,flag,appId)
	if err!=nil{
		log.Error(err)
		return nil,err
	}
	if couponuserList==nil||len(couponuserList)<=0 {
		log.Info("用户没有优惠券！")
		return nil,nil
	}
	orderDetail,err := GetOrderDetail(orderNo,appId)
	if err!=nil{
		log.Error(err)
		return nil,err
	}
	if orderDetail==nil{

		return nil,errors.New("订单信息没找到!")
	}
	//目前暂时只支持一种优惠使用 不支持多种优惠同时使用
	couponuser = couponuserList[0]
	if couponuser.Balance<=0 {
		log.Info("券余额不足！")
		return nil,nil
	}

	if orderDetail.RealPrice<=0.01 {
		log.Info("商品真实金额小于或等于0.01！")
		return nil,nil
	}

	couponAmount := getCouponAmount(orderDetail)
	if couponAmount> couponuser.Balance {
		couponAmount = couponuser.Balance
	}

	trackCode :=util.GenerUUId()
	result :=&CouponUser{}
	result.Title = couponuser.Title
	result.Remark = couponuser.Remark
	result.AppId = appId
	result.OpenId = openId
	result.CouponAmount = couponAmount
	result.CouponCode = couponuser.CouponCode
	result.OrderNo = orderDetail.No
	result.TrackCode = trackCode
	jwtauth := comm.InitJWTAuthenticationBackend()
	notifyUrl :=config.GetValue("coupon_notify_url").ToString()
	token,err :=jwtauth.GenerateCouponToken(openId,result.CouponCode,trackCode,result.OrderNo,result.CouponAmount,notifyUrl,appId)
	if err!=nil{
		log.Error(err)
		return nil,err
	}
	result.CouponToken = token

	tx,_ := db.NewSession().Begin()
	couponTrack := dao.NewCouponTrack()
	couponTrack.OpenId = result.OpenId
	couponTrack.AppId = result.AppId
	couponTrack.CouponAmount = result.CouponAmount
	couponTrack.Amount = orderDetail.RealPrice
	couponTrack.Title=couponuser.Title
	couponTrack.Remark = couponuser.Remark
	couponTrack.CouponCode = result.CouponCode
	couponTrack.Status = 0
	couponTrack.TrackType = 1
	couponTrack.TrackCode = trackCode
	couponTrack.TradeNo = orderNo
	err =couponTrack.InsertTx(tx)
	if err!=nil{
		log.Error(err)
		tx.Rollback()
		return nil,err
	}
	tx.Commit()

	return result,nil
}

func getCouponAmount(order *OrderDetailDto) float64  {
	items := order.Items
	var totalCouponAmount float64
	if items!=nil{
		for _,item :=range items {
			//只有标记为此商品的才进行优惠
			if item.ProdFlag=="merchant_default" {
				totalCouponAmount += comm.Round(item.BuyTotalPrice/2.0,2)
			}
		}
	}

	return totalCouponAmount
}