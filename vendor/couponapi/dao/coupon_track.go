package dao

import "github.com/gocraft/dbr"

type CouponTrack struct  {
	Id int64
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

	_,err :=tx.InsertInto("coupon_track").Columns("trade_no","trade_type","track_code","open_id","coupon_code","title","remark","amount","track_type","coupon_amount","status").Record(self).Exec()

	return err
}