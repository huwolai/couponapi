package dao

type Coupon struct  {
	Id int64
	AppId string
	VCode string
	Amount float64
	PublishNum int
	PublishedNum int
	IsOne int
	Status int
	Title string
	Remark string
}

func (self *Coupon) Insert() error  {

	return nil
}