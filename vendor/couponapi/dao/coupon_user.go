package dao

import (
	"github.com/gocraft/dbr"
	"gitlab.qiyunxin.com/tangtao/utils/db"
	"strings"
	"strconv"
	"gitlab.qiyunxin.com/tangtao/utils/log"
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

	return nil,nil

}

func (self *CouponUser) UpdateAmountAndBalanceWithId(amount float64,balance float64,id int64) error  {

	_,err :=db.NewSession().Update("coupon_user").Set("amount",amount).Set("balance",balance).Where("id=?",id).Exec()

	return err
}

func (self *CouponUser) TotalAmountWithOpenId(openId string,appId string) (float64,error)  {
	var amount []byte
	err :=db.NewSession().Select("sum(balance)").From("coupon_user").Where("open_id=?",openId).Where("app_id=?",appId).LoadValue(&amount)
	if amount==nil{
		return 0,nil
	}
	log.Error(amount)
	famount,_ :=strconv.ParseFloat(string(amount),10)
	return famount,err
}
