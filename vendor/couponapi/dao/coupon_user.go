package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"strings"
)

type CouponUser struct  {
	Id int64
	AppId string
	OpenId string
	CouponCode string
	Title string
	Remark string
	Amount float64
	Flag string
	Balance float64
	IsOne int
	UseStatus int
}

func NewCouponUser() *CouponUser  {

	return &CouponUser{}
}

func (self *CouponUser) InsertTx(tx *dbr.Tx) error  {

	_,err :=tx.InsertInto("coupon_user").Columns("app_id","open_id","coupon_code","title","remark","amount","balance","is_one","use_status","flag").Record(self).Exec()

	return err
}



func (self *CouponUser) WithCodesOrFlag(openId string,codes []string,flag string,appId string) ([]*CouponUser,error)  {
	var list []*CouponUser
	if codes!=nil&&len(codes)>0 {
		_,err :=db.NewSession().SelectBySql("select * from coupon_user where open_id=? and app_id=? and coupon_code in ?  ORDER BY INSTR(',?,',CONCAT(',',coupon_code,','))",openId,appId,codes,strings.Join(codes,",")).LoadStructs(&list)
		return list,err
	}

	if flag!=""{
		_,err :=db.NewSession().SelectBySql("select * from coupon_user where open_id=? and app_id=? and flag=?",openId,appId,flag).LoadStructs(&list)
		return list,err
	}

	_,err :=db.NewSession().SelectBySql("select * from coupon_user where open_id=? and app_id=?",openId,appId).LoadStructs(&list)
	return list,err

}

//根据优惠券代码查询用户优惠券
func (self *CouponUser) WithCouponCode(openId string,code string,appId string) (*CouponUser,error)  {
	var model *CouponUser
	_,err :=db.NewSession().Select("*").From("coupon_user").Where("coupon_code=?",code).Where("open_id=?",openId).Where("app_id=?",appId).LoadStructs(&model)
	return model,err
}

//更新优惠券面值跟余额
func (self *CouponUser) UpdateAmountAndBalanceWithId(amount float64,balance float64,id int64) error  {

	_,err :=db.NewSession().Update("coupon_user").Set("amount",amount).Set("balance",balance).Where("id=?",id).Exec()

	return err
}

func (self *CouponUser) UpdateBalanceWithCouponCodeTx(balance float64,couponCode string,appId string,tx *dbr.Tx) error  {

	_,err :=tx.Update("coupon_user").Set("balance",balance).Where("coupon_code=?",couponCode).Where("app_id=?",appId).Exec()

	return err
}

func (self *CouponUser) TotalAmountWithOpenId(status int,openId string,appId string) (float64,error)  {
	var amount dbr.NullFloat64
	err :=db.NewSession().Select("sum(balance)").From("coupon_user").Where("open_id=?",openId).Where("app_id=?",appId).LoadValue(&amount)

	return amount.Float64,err
}
