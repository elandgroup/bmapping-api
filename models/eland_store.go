package models

import (
	"context"
	"errors"

	"github.com/go-xorm/xorm"

	"bmapping-api/factory"
)

type ElandStore struct {
	Id        int64  `json:"id" query:"id" xorm:"pk autoincr 'id'"`
	GroupCode string `json:"groupCode" query:"groupCode" xorm:"VARCHAR(4) unique(group_code_code_country_id) 'group_code'"`
	Code      string `json:"code" query:"code" xorm:"VARCHAR(4) unique(group_code_code_country_id) 'code'"`
	CountryId int64  `json:"countryId" query:"countryId" xorm:"unique(group_code_code_country_id) 'country_id'"`
	Name      string `json:"name" query:"name"`
	CreatedAt string `json:"createdAt" query:"createdAt" xorm:"created"`
	UpdatedAt string `json:"updatedAt" query:"updatedAt" xorm:"updated"`
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
	row, err := factory.DB(ctx).InsertOne(d)
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
		v, err := queryBuilder().Count(&ElandStore{})
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
func (d *ElandStore) Update(ctx context.Context) (err error) {
	row, err := factory.DB(ctx).ID(d.Id).Update(d)
	if int(row) == 0 {
		err = errors.New("no data has changed.")
		return
	}
	return
}

func (ElandStore) Delete(ctx context.Context, id int64) (err error) {
	row, err := factory.DB(ctx).ID(id).Delete(&ElandStore{})
	if int(row) == 0 {
		err = errors.New("no data has changed.")
		return
	}
	return
}
