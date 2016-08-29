package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
)

type CouponTrack struct  {
	Id int64
	AppId string
	TradeNo string
	TradeType int
	TrackCode string
	OpenId string
	CouponCode string
	Title string
	Remark string
	Amount float64
	TrackType int
	CouponAmount float64
	Status int
}

func NewCouponTrack() *CouponTrack  {

	return &CouponTrack{}
}

func (self *CouponTrack) InsertTx(tx *dbr.Tx) error  {

	_,err :=tx.InsertInto("coupon_track").Columns("trade_no","trade_type","track_code","open_id","coupon_code","title","remark","amount","track_type","coupon_amount","status","app_id").Record(self).Exec()

	return err
}

func (self *CouponTrack) WithTrackCode(trackCode string) (*CouponTrack,error)  {

	var model *CouponTrack
	_,err :=db.NewSession().Select("*").From("coupon_track").Where("track_code=?",trackCode).LoadStructs(&model)

	return model,err
}

func (self *CouponTrack) UpdateTradeAndStatusWithTrackCodeTx(status int,tradeType int,tradeNo string,trackCode string,tx *dbr.Tx) error  {

	_,err :=tx.Update("coupon_track").Set("trade_no",tradeNo).Set("trade_type",tradeType).Set("status",status).Where("trade_no=?",trackCode).Exec()

	return err
}