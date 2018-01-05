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
	Description string    `josn:"description" query:"description"`
	CreatedAt   time.Time `xorm:"created" query:"createdAt" xorm:"created"`
	UpdatedAt   time.Time `xorm:"updated" query:"updateAt" xorm:"created"`
}

func GetEIdByCode(ctx context.Context, code string, ipayTypeId int64) (has bool, eId int64, err error) {
	var greenMapping GreenMappingStoreIpay
	has, err = factory.DB(ctx).Table("green_store").Alias("a").
		Select("b.e_id").
		Join("inner", []string{"green_mapping_store_ipay", "b"}, "a.id=b.store_id").
		Where("a.code=?", code).And("b.ipay_type_id=?", ipayTypeId).
		Get(&greenMapping)
	eId = greenMapping.EId
	return
}
