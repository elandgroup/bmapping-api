package models

import (
	"bmapping-api/factory"
	"context"
	"time"
)

type ElandStore struct {
	Id          int64
	GroupId     int64
	Code        string
	Name        string
	LogoUrl     string
	QrUrl       string
	Description string
	CreatedAt   time.Time `xorm:"created"`
	UpdatedAt   time.Time `xorm:"updated"`
}

func GetEIdByThrArgs(ctx context.Context, group_code string, code string, country_id int64, ipayTypeId int64) (has bool, eId int64, err error) {
	var ElandMapping ElandMappingStoreIpay
	has, err = factory.DB(ctx).Table("eland_mapping_store_ipay").Alias("a").
		Select("a.e_id").
		Join("inner", []string{"eland_store", "b"}, "a.store_id=b.id").
		Join("inner", []string{"eland_store_group", "c"}, "b.group_id=c.id").
		Where("b.code=?", code).And("c.country_id=?", country_id).And("c.code=?", group_code).And("a.ipay_type_id=?", ipayTypeId).
		Get(&ElandMapping)
	eId = ElandMapping.EId
	return
}
