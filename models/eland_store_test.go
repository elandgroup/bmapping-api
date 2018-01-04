package models

import (
	"fmt"
	"testing"
	"time"

	"github.com/pangpanglabs/goutils/test"
)

func TestElandStoreCreate(t *testing.T) {
	d := ElandStore {
	    Id:int64(1),
	    GroupCode:"test",
	    Code:"test",
	    CountryId:int64(1),
	    Name:"test",
	    CreatedAt:"test",
	    UpdatedAt:"test",

	}
	affected, err := d.Create(ctx)
	test.Ok(t, err)

	test.Equals(t, affected, int64(1))
	test.Equals(t, d.Id, int64(1))
	test.Equals(t, d.GroupCode, "test")
	test.Equals(t, d.Code, "test")
	test.Equals(t, d.CountryId, int64(1))
	test.Equals(t, d.Name, "test")
	test.Equals(t, d.CreatedAt, "test")
	test.Equals(t, d.UpdatedAt, "test")

}

func TestElandStoreGetAndUpdate(t *testing.T) {
	d, err := Discount{}.GetById(ctx, 1)
	test.Ok(t, err)
	test.Equals(t, d.Id, int64(1))
	test.Equals(t, d.GroupCode, "test")
	test.Equals(t, d.Code, "test")
	test.Equals(t, d.CountryId, int64(1))
	test.Equals(t, d.Name, "test")
	test.Equals(t, d.CreatedAt, "test")
	test.Equals(t, d.UpdatedAt, "test")

    d.GroupCode = "test"
    d.Code = "test"
	d.CountryId = int64(2)
	d.Name = "test"
    d.CreatedAt = "test"
    d.UpdatedAt = "test"

	err = d.Update(ctx)
	test.Ok(t, err)

	test.Equals(t, d.GroupCode, "test2")
	test.Equals(t, d.Code, "test2")
	test.Equals(t, d.CountryId, int64(2)test.Equals(t, d.Name, "test2")
	test.Equals(t, d.CreatedAt, "test2")
	test.Equals(t, d.UpdatedAt, "test2")

}

func TestElandStoreGetAll(t *testing.T) {
	totalCount, items, err := ElandStore{}.GetAll(ctx, []string{"test"}, []string{"test"}, 0, 10)
	test.Ok(t, err)
	test.Equals(t, totalCount, int64(1))
	test.Equals(t, items[0].Id, int64(1))
}

func TestXXX(t *testing.T) {
	at, err := time.Parse("2006-01-02", "2017-12-31")
	test.Ok(t, err)
	test.Equals(t, at.Year(), 2017)
	test.Assert(t, at.Month() == 12, "Month should be equals to 12")
	fmt.Println(at)
}
