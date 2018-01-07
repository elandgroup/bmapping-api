package models

import (
	"fmt"
	"testing"

	"github.com/relax-space/go-kit/test"
)

func Test_ElandStore_GetEId(t *testing.T) {
	has, eId, store, err := GetEIdByThrArgs(ctx, "AA01", "CR00", 2, 1)
	fmt.Println(has, eId, *store)
	test.Ok(t, err)
}

//group

func Test_ElandStore_InsertMany(t *testing.T) {
	d := ElandStoreGroup{
		Code:      "test_01",
		CountryId: int64(1),
		Name:      "test_01",
	}
	var stores []ElandStoreGroup
	stores = append(stores, d)

	err := (ElandStoreGroup{}).InsertMany(ctx, &stores)
	test.Ok(t, err)
}

func Test_ElandStore_GetById(t *testing.T) {
	has, elandStoreGroup, err := ElandStoreGroup{}.GetById(ctx, 1)
	fmt.Println(has, elandStoreGroup)
	test.Ok(t, err)
}

func Test_ElandStore_GetAll(t *testing.T) {
	count, items, err := ElandStoreGroup{}.GetAll(ctx, nil, nil, 0, 2)
	fmt.Println(count, items)
	test.Ok(t, err)
}

func Test_ElandStore_Update(t *testing.T) {
	d := &ElandStoreGroup{
		Id:   1,
		Name: "test_03",
	}
	err := d.Update(ctx)
	test.Ok(t, err)
}
func Test_ElandStore_InsertOne(t *testing.T) {
	d := &ElandStoreGroup{
		Code:      "test_01",
		CountryId: int64(2),
		Name:      "test_05",
	}
	err := d.InsertOne(ctx)
	test.Ok(t, err)
}
func Test_ElandStore_Delete(t *testing.T) {
	err := ElandStoreGroup{}.Delete(ctx, 3)
	test.Ok(t, err)
}
