package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
)

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

func NewCouponUser() *CouponUser  {

	return &CouponUser{}
}

func (self *CouponUser) InsertTx(tx *dbr.Tx) error  {

	_,err :=tx.InsertInto("coupon_user").Columns("app_id","open_id","coupon_code","title","remark","amount","balance","is_one","use_status").Record(self).Exec()

	return err
}

func (self *CouponUser) TotalAmountWithOpenId(openId string,appId string) (float64,error)  {
	var amount float64
	err :=db.NewSession().Select("sum(balance)").From("coupon_user").Where("open_id=?",openId).Where("app_id=?",appId).LoadValue(&amount)
	return amount,err
}
