package models

import (
	"context"
	"errors"

	"github.com/go-xorm/xorm"

	"bmapping-api/factory"
)

type ElandStoreGroup struct {
	Id        int64  `json:"id" query:"id" xorm:"pk autoincr 'id'"`
	Code      string `json:"code" query:"code" xorm:"VARCHAR(4) unique(code_country_id) 'code'"`
	CountryId int64  `json:"countryId" query:"countryId" xorm:"unique(code_country_id) 'country_id'"`
	Name      string `json:"name" query:"name"`
	CreatedAt string `json:"createdAt" query:"createdAt" xorm:"created"`
	UpdatedAt string `json:"updatedAt" query:"updatedAt" xorm:"updated"`
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
		v, err := queryBuilder().Count(&ElandStoreGroup{})
		if err != nil {
			errc <- err
			return
		}
		totalCount = v

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
