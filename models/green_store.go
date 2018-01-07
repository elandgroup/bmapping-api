package models

import (
	"bmapping-api/factory"
	"context"
	"time"
)

type GreenStore struct {
	Id          int64     `json:"id" query:"id" xorm:"pk autoincr 'id'"`
	GroupId     int64     `json:"groupId" query:"groupId" `
	Code        string    `json:"code" query:"code" xorm:"VARCHAR(4) unique 'code'"`
	Name        string    `json:"name" query:"name"`
	LogoUrl     string    `json:"logoUrl" query:"logoUrl"`
	QrUrl       string    `json:"qrUry" query:"qrUry"`
	Description string    `json:"description" query:"description"`
	CreatedAt   time.Time `json:"createAt" query:"createdAt" xorm:"created"`
	UpdatedAt   time.Time `json:"updateAt" query:"updateAt" xorm:"updated"`
}

func GetEIdByCode(ctx context.Context, code string, ipayTypeId int64) (has bool, eId int64, greenStore *GreenStore, err error) {
	var bizStore struct {
		GreenStore *GreenStore `json:"greenStore" xorm:"extends"`
		EId        int64       `json:"eId"`
	}
	has, err = factory.DB(ctx).Table("green_store").Alias("a").
		Select(`b.e_id,
			a.name,
			a.logo_url,
			a.qr_url,
			a.code
			`).
		Join("inner", []string{"green_mapping_store_ipay", "b"}, "a.id=b.store_id").
		Where("a.code=?", code).And("b.ipay_type_id=?", ipayTypeId).
		Get(&bizStore)
	eId = bizStore.EId
	greenStore = bizStore.GreenStore
	return
}

type GreenMappingStoreIpay struct {
	Id         int64
	StoreId    int64
	IpayTypeId int64
	EId        int64
	CreatedAt  time.Time `xorm:"created"`
	UpdatedAt  time.Time `xorm:"updated"`
}

type GreenStoreGroup struct {
	Id        int64     `json:"id" query:"id" xorm:"pk autoincr 'id'"`
	Code      string    `json:"code" query:"code" xorm:"VARCHAR(4) unique 'code'"`
	Name      string    `json:"name" query:"name"`
	CreatedAt time.Time `json:"createdAt" query:"createdAt" xorm:"created"`
	UpdatedAt time.Time `json:"updatedAt" query:"updatedAt" xorm:"updated"`
}
