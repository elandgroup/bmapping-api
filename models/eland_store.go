package models

import (
	"bmapping-api/factory"
	"context"
	"errors"
	"time"

	"github.com/go-xorm/xorm"
)

type ElandStore struct {
	Id          int64     `json:"id"`
	GroupId     int64     `json:"groupId"`
	Code        string    `json:"codeId"`
	Name        string    `json:"name"`
	LogoUrl     string    `json:"logoUrl"`
	QrUrl       string    `json:"qrUrl"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt" xorm:"created"`
	UpdatedAt   time.Time `json:"updatedAt" xorm:"updated"`
}

func GetEIdByThrArgs(ctx context.Context, group_code string, code string,
	country_id int64, ipayTypeId int64) (has bool, eId int64, elandStore *ElandStore, err error) {
	var bizStore struct {
		ElandStore *ElandStore `json:"elandStore" xorm:"extends"`
		EId        int64       `json:"eId"`
	}
	has, err = factory.DB(ctx).Table("eland_mapping_store_ipay").Alias("a").
		Select(`
			a.e_id,
			b.name,
			b.logo_url,
			b.qr_url,
			b.code
			`).
		Join("inner", []string{"eland_store", "b"}, "a.store_id=b.id").
		Join("inner", []string{"eland_store_group", "c"}, "b.group_id=c.id").
		Where("b.code=?", code).And("c.country_id=?", country_id).And("c.code=?", group_code).And("a.ipay_type_id=?", ipayTypeId).
		Get(&bizStore)
	eId = bizStore.EId
	elandStore = bizStore.ElandStore
	return
}

func (ElandStore) InsertMany(ctx context.Context, elandStore *[]ElandStore) (err error) {
	row, err := factory.DB(ctx).Insert(elandStore)
	if int(row) == 0 {
		err = errors.New("no data has changed.")
		return
	}
	return
}

func (d *ElandStore) InsertOne(ctx context.Context) (err error) {
	row, err := factory.DB(ctx).Insert(d)
	if int(row) == 0 {
		err = errors.New("no data has changed.")
		return
	}
	return
}

func (ElandStore) GetById(ctx context.Context, id int64) (has bool, elandStore *ElandStore, err error) {
	elandStore = &ElandStore{}
	has, err = factory.DB(ctx).ID(id).Get(elandStore)
	return
}

func (ElandStore) GetAll(ctx context.Context, sortby, order []string, offset, limit int) (totalCount int64, items []ElandStore, err error) {
	queryBuilder := func() *xorm.Session {
		q := factory.DB(ctx)
		if err := setSortOrder(q, sortby, order); err != nil {
			factory.Logger(ctx).Error(err)
		}
		return q
	}

	errc := make(chan error)
	go func() {
		v, err := queryBuilder().Count(&ElandStore{})
		if err != nil {
			errc <- err
			return
		}
		totalCount = v
		errc <- nil
	}()

	go func() {
		if err := queryBuilder().Limit(limit, offset).Find(&items); err != nil {
			errc <- err
			return
		}
		errc <- nil
	}()

	if err := <-errc; err != nil {
		return 0, nil, err
	}
	if err := <-errc; err != nil {
		return 0, nil, err
	}
	return
}

func (d *ElandStore) Update(ctx context.Context, id int64) (err error) {
	row, err := factory.DB(ctx).ID(id).Update(d)
	if int(row) == 0 {
		err = errors.New("no data has changed.")
		return err
	}
	return
}

func (ElandStore) Delete(ctx context.Context, id int64) (err error) {
	row, err := factory.DB(ctx).ID(id).Delete(&ElandStore{})
	if int(row) == 0 {
		err = errors.New("no data has changed.")
		return err
	}
	return
}

type ElandStoreGroup struct {
	Id        int64     `json:"id" query:"id" xorm:"pk autoincr 'id'"`
	Code      string    `json:"code" query:"code" xorm:"VARCHAR(4) unique(code_country_id) 'code'"`
	CountryId int64     `json:"countryId" query:"countryId" xorm:"unique(code_country_id) 'country_id'"`
	Name      string    `json:"name" query:"name"`
	CreatedAt time.Time `json:"createdAt" query:"createdAt" xorm:"created"`
	UpdatedAt time.Time `json:"updatedAt" query:"updatedAt" xorm:"updated"`
}

func (ElandStoreGroup) InsertMany(ctx context.Context, elandStoreGroup *[]ElandStoreGroup) (err error) {
	row, err := factory.DB(ctx).Insert(elandStoreGroup)
	if int(row) == 0 {
		err = errors.New("no data has changed.")
		return
	}
	return
}
func (d *ElandStoreGroup) InsertOne(ctx context.Context) (err error) {
	row, err := factory.DB(ctx).InsertOne(d)
	if int(row) == 0 {
		err = errors.New("no data has changed.")
		return
	}
	return
}

func (ElandStoreGroup) GetById(ctx context.Context, id int64) (has bool, elandStoreGroup *ElandStoreGroup, err error) {
	elandStoreGroup = &ElandStoreGroup{}
	has, err = factory.DB(ctx).ID(id).Get(elandStoreGroup)
	return
}
func (ElandStoreGroup) GetAll(ctx context.Context, sortby, order []string, offset, limit int) (totalCount int64, items []ElandStoreGroup, err error) {
	queryBuilder := func() *xorm.Session {
		q := factory.DB(ctx)
		if err := setSortOrder(q, sortby, order); err != nil {
			factory.Logger(ctx).Error(err)
		}
		return q
	}

	errc := make(chan error)
	go func() {
		v, err := queryBuilder().Count(&ElandStoreGroup{})
		if err != nil {
			errc <- err
			return
		}
		totalCount = v
		errc <- nil

	}()

	go func() {
		if err := queryBuilder().Limit(limit, offset).Find(&items); err != nil {
			errc <- err
			return
		}
		errc <- nil
	}()

	if err := <-errc; err != nil {
		return 0, nil, err
	}
	if err := <-errc; err != nil {
		return 0, nil, err
	}
	return
}
func (d *ElandStoreGroup) Update(ctx context.Context) (err error) {
	row, err := factory.DB(ctx).ID(d.Id).Update(d)
	if int(row) == 0 {
		err = errors.New("no data has changed.")
		return
	}
	return
}

func (ElandStoreGroup) Delete(ctx context.Context, id int64) (err error) {
	row, err := factory.DB(ctx).ID(id).Delete(&ElandStoreGroup{})
	if int(row) == 0 {
		err = errors.New("no data has changed.")
		return
	}
	return
}

type ElandMappingStoreIpay struct {
	Id         int64
	StoreId    int64
	IpayTypeId int64
	EId        int64
	CreatedAt  time.Time `xorm:"created"`
	UpdatedAt  time.Time `xorm:"updated"`
}
