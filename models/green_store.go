package models

import (
	"bmapping-api/factory"
	"context"
	"errors"
	"time"

	"github.com/go-xorm/xorm"
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

/**↓ greenStore CRUD ↓**/
func (GreenStore) InsertMany(ctx context.Context, greenStores *[]GreenStore) (err error) {
	row, err := factory.DB(ctx).Insert(greenStores)
	if int(row) == 0 {
		err = errors.New("no data has changed.")
		return
	}
	return
}

func (c *GreenStore) InsertOne(ctx context.Context) (err error) {
	row, err := factory.DB(ctx).InsertOne(c)
	if int(row) == 0 {
		err = errors.New("no data has changed.")
		return
	}
	return
}

func (GreenStore) GetAll(ctx context.Context, sortby, order []string, offset, limit int) (totalCount int64, items []GreenStore, err error) {
	queryBuilder := func() *xorm.Session {
		q := factory.DB(ctx)
		if err := setSortOrder(q, sortby, order); err != nil {
			factory.Logger(ctx).Error(err)
		}
		return q
	}

	errc := make(chan error)
	go func() {
		v, err := queryBuilder().Count(&GreenStore{})
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

func (GreenStore) GetById(ctx context.Context, id int64) (has bool, greenStore *GreenStore, err error) {
	greenStore = &GreenStore{}
	has, err = factory.DB(ctx).ID(id).Get(greenStore)
	return
}

func (t *GreenStore) Update(ctx context.Context, id int64) (err error) {
	row, err := factory.DB(ctx).ID(id).Update(t)
	if int(row) == 0 {
		err = errors.New("no data has changed.")
		return
	}
	return
}

func (GreenStore) Delete(ctx context.Context, id int64) (err error) {
	row, err := factory.DB(ctx).ID(id).Delete(&GreenStore{})
	if int(row) == 0 {
		err = errors.New("no data has changed.")
		return
	}
	return
}
/**  ↑GreenStore UCRD↑  **/

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

func (GreenStoreGroup) InsertMany(ctx context.Context, greenStoreGroup *[]GreenStoreGroup) (err error) {
	row, err := factory.DB(ctx).Insert(greenStoreGroup)
	if int(row) == 0 {
		err = errors.New("no data has changed.")
		return
	}
	return
}

func (t *GreenStoreGroup) InsertOne(ctx context.Context) (err error) {
	row, err := factory.DB(ctx).InsertOne(t)
	if int(row) == 0 {
		err = errors.New("no data has changed.")
		return
	}
	return
}

func (GreenStoreGroup) GetAll(ctx context.Context, sortby, order []string, offset, limit int) (totalCount int64, items []GreenStoreGroup, err error) {
	queryBuilder := func() *xorm.Session {
		q := factory.DB(ctx)
		if err := setSortOrder(q, sortby, order); err != nil {
			factory.Logger(ctx).Error(err)
		}
		return q
	}

	errc := make(chan error)
	go func() {
		v, err := queryBuilder().Count(&GreenStoreGroup{})
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

func (GreenStoreGroup) GetById(ctx context.Context, id int64) (has bool, greenStoreGroup *GreenStoreGroup, err error) {
	greenStoreGroup = &GreenStoreGroup{}
	has, err = factory.DB(ctx).ID(id).Get(greenStoreGroup)
	return
}

func (t *GreenStoreGroup) Update(ctx context.Context, id int64) (err error) {
	row, err := factory.DB(ctx).ID(id).Update(t)
	if int(row) == 0 {
		err = errors.New("no data has changed.")
		return
	}
	return
}

func (GreenStoreGroup) Delete(ctx context.Context, id int64) (err error) {
	row, err := factory.DB(ctx).ID(id).Delete(&GreenStoreGroup{})
	if int(row) == 0 {
		err = errors.New("no data has changed.")
		return
	}
	return
}
