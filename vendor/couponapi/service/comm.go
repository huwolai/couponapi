package service

import (
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"gitlab.qiyunxin.com/tangtao/utils/network"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"errors"
	"gitlab.qiyunxin.com/tangtao/utils/util"
)

type OrderDetailDto struct  {
	Id int64 `json:"id"`
	No string `json:"no"`
	PayapiNo string `json:"payapi_no"`
	OpenId string `json:"open_id"`
	Name string `json:"name,omitempty"`
	Mobile string `json:"mobile,omitempty"`
	AddressId int64 `json:"address_id"`
	Address string `json:"address"`
	AppId string `json:"app_id"`
	Title string `json:"title"`
	Price float64 `json:"price"`
	RealPrice float64 `json:"real_price"`
	PayPrice float64 `json:"pay_price"`
	OmitMoney float64 `json:"omit_money"`
	RejectCancelReason string `json:"reject_cancel_reason"`
	CancelReason string `json:"cancel_reason"`
	OrderStatus int `json:"order_status"`
	PayStatus int `json:"pay_status"`
	Items []*OrderItemDetailDto `json:"items"`
	Json string `json:"json"`
	CreateTime string `json:"create_time"`

}

type OrderItemDetailDto struct  {
	Id int64 `json:"id"`
	No string `json:"no"`
	AppId string `json:"app_id"`
	OpenId string `json:"open_id"`
	//商户名称
	MerchantName string `json:"merchant_name"`
	//商户ID
	MerchantId int64 `json:"merchant_id"`
	//商品cover 封面图 url
	ProdCoverImg string `json:"prod_coverimg"`
	ProdTitle string `json:"prod_title"`
	ProdId int64 `json:"prod_id"`
	Num int `json:"num"`
	OfferUnitPrice float64 `json:"offer_unit_price"`
	OfferTotalPrice float64 `json:"offer_total_price"`
	BuyUnitPrice float64 `json:"buy_unit_price"`
	BuyTotalPrice float64 `json:"buy_total_price"`
	Json string `json:"json"`
}


//获取订单详情
func GetOrderDetail(orderNo string,openId string) (*OrderDetailDto,error)  {

	shopUrl :=config.GetValue("shopapi_url").ToString()
	shopappId :=config.GetValue("shopapi_appid").ToString()

	data,err := network.GetJson(shopUrl+"/order/"+orderNo+"/detail",map[string]string{
		"open_id":openId,
		"app_id":shopappId,
	},nil)
	if err!=nil{
		log.Error(err)
		return nil,errors.New("调用shopapi获取订单详情失败!")
	}
	if data!=nil{
		var  dto *OrderDetailDto
		err :=util.ReadJsonByByte(data,&dto)
		if err!=nil{
			log.Error(string(data))
			log.Error(err)
			return nil,err
		}
		return dto,err
	}

	return nil,nil

}
