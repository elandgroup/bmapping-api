package models

import (
	"context"

	"github.com/go-xorm/xorm"

	"bmapping-api/factory"
)

type ElandStore struct {
	Id	int64	`json:"id"`
	GroupCode	string	`json:"groupCode"`
	Code	string	`json:"code"`
	CountryId	int64	`json:"countryId"`
	Name	string	`json:"name"`
	CreatedAt	string	`json:"createdAt"`
	UpdatedAt	string	`json:"updatedAt"`

}

func (d *ElandStore) Create(ctx context.Context) (int64, error) {
	return factory.DB(ctx).Insert(d)
}
func (ElandStore) GetById(ctx context.Context, id int64) (*ElandStore, error) {
	var v ElandStore
	if has, err := factory.DB(ctx).ID(id).Get(&v); err != nil {
		return nil, err
	} else if !has {
		return nil, nil
	}
	return &v, nil
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
	_, err = factory.DB(ctx).ID(d.Id).Update(d)
	return
}

func (ElandStore) Delete(ctx context.Context, id int64) (err error) {
	_, err = factory.DB(ctx).ID(id).Delete(&ElandStore{})
	return
}
