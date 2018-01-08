package models

import (
	"bmapping-api/factory"
	"context"
	"errors"

	"github.com/go-xorm/xorm"
)

type IpayType struct {
	Id          int64
	Description string
}

func (IpayType) InsertMany(ctx context.Context, ipayTypes *[]IpayType) (err error) {
	row, err := factory.DB(ctx).Insert(ipayTypes)
	if int(row) == 0 {
		err = errors.New("no data has changed.")
		return
	}
	return
}

func (c *IpayType) InsertOne(ctx context.Context) (err error) {
	row, err := factory.DB(ctx).InsertOne(c)
	if int(row) == 0 {
		err = errors.New("no data has changed.")
		return
	}
	return
}

func (IpayType) GetById(ctx context.Context, id int64) (has bool, ipayType *IpayType, err error) {
	ipayType = &IpayType{}
	has, err = factory.DB(ctx).ID(id).Get(ipayType)
	return
}

func (IpayType) GetAll(ctx context.Context, sortby, order []string, offset, limit int) (totalCount int64, items []IpayType, err error) {
	queryBuilder := func() *xorm.Session {
		q := factory.DB(ctx)
		if err := setSortOrder(q, sortby, order); err != nil {
			factory.Logger(ctx).Error(err)
		}
		return q
	}

	errc := make(chan error)
	go func() {
		v, err := queryBuilder().Count(&IpayType{})
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
func (c *IpayType) Update(ctx context.Context, id int64) (err error) {
	row, err := factory.DB(ctx).ID(id).Update(c)
	if int(row) == 0 {
		err = errors.New("no data has changed.")
		return
	}
	return
}

func (IpayType) Delete(ctx context.Context, id int64) (err error) {
	row, err := factory.DB(ctx).ID(id).Delete(&IpayType{})
	if int(row) == 0 {
		err = errors.New("no data has changed.")
		return
	}
	return
}
