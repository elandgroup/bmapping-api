package models

import (
	"fmt"
	"testing"

	"github.com/relax-space/go-kit/test"
)

func Test_IpayType_InsertMany(t *testing.T) {
	ipayType := IpayType{
		Description: "test insert many  by xiao",
	}
	ipayTypes := []IpayType{
		ipayType,
	}
	err := IpayType{}.InsertMany(ctx, &ipayTypes)
	test.Ok(t, err)
}

func Test_IpayType_GetById(t *testing.T) {
	has, ipayType, err := IpayType{}.GetById(ctx, 1)
	fmt.Println(has, ipayType)
	test.Ok(t, err)
}

func Test_IpayType_GetAll(t *testing.T) {
	count, items, err := IpayType{}.GetAll(ctx, nil, nil, 0, 1)
	fmt.Println(count, items)
	test.Ok(t, err)
}

func Test_IpayType_Update(t *testing.T) {
	ipayType := &IpayType{
		Description: "test insert one by xiao",
	}
	err := ipayType.Update(ctx, 1)
	test.Ok(t, err)
}

func Test_IpayType_InsertOne(t *testing.T) {
	ipayType := &IpayType{
		Description: "test insert one by xiao",
	}
	err := ipayType.InsertOne(ctx)
	test.Ok(t, err)
}

func Test_IpayType_Delete(t *testing.T) {
	err := IpayType{}.Delete(ctx, 2)
	test.Ok(t, err)
}
