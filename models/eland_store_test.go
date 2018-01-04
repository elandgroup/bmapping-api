package models

import (
	"fmt"
	"testing"

	"github.com/pangpanglabs/goutils/test"
)

func Test_ElandStore_InsertMany(t *testing.T) {
	d := ElandStore{
		GroupCode: "test_01",
		Code:      "test_01",
		CountryId: int64(1),
		Name:      "test_01",
	}
	var stores []ElandStore
	stores = append(stores, d)

	err := (ElandStore{}).InsertMany(ctx, &stores)
	test.Ok(t, err)
}

func Test_ElandStore_GetById(t *testing.T) {
	has, elandStore, err := ElandStore{}.GetById(ctx, 1)
	fmt.Println(has, elandStore)
	test.Ok(t, err)
}

func Test_ElandStore_GetAll(t *testing.T) {
	count, items, err := ElandStore{}.GetAll(ctx, nil, nil, 0, 2)
	fmt.Println(count, items)
	test.Ok(t, err)
}

func Test_ElandStore_Update(t *testing.T) {
	d := &ElandStore{
		Id:   1,
		Name: "test_03",
	}
	err := d.Update(ctx)
	test.Ok(t, err)
}
func Test_ElandStore_InsertOne(t *testing.T) {
	d := &ElandStore{
		GroupCode: "test_01",
		Code:      "test_01",
		CountryId: int64(2),
		Name:      "test_05",
	}
	err := d.InsertOne(ctx)
	test.Ok(t, err)
}
func Test_ElandStore_Delete(t *testing.T) {
	err := ElandStore{}.Delete(ctx, 3)
	test.Ok(t, err)
}
