package models

import "time"

type GreenMappingStoreIpay struct {
	Id         int64
	StoreId    int64
	IpayTypeId int64
	EId        int64
	CreatedAt  time.Time `xorm:"created"`
	UpdatedAt  time.Time `xorm:"updated"`
}
