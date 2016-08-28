package dao

import "github.com/gocraft/dbr"

type CouponUser struct  {
	Id int64
	AppId string
	OpenId string
	CouponCode string
	Title string
	Remark string
	Amount float64
	Balance float64
	IsOne int
	UseStatus int
}

func (self *CouponUser) InsertTx(tx *dbr.Tx) error  {

	_,err :=tx.InsertInto("coupon_user").Columns("app_id","open_id","coupon_code","title","remark","amount","balance","is_one","use_status").Exec()

	return err
}
