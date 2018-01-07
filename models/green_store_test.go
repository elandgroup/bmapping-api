package models

import (
	"fmt"
	"testing"

	"github.com/relax-space/go-kit/test"
)

func Test_GreenStore_GetEIdByCode(t *testing.T) {
	has, eId, store, err := GetEIdByCode(ctx, "AA01", 1)
	fmt.Println(has, eId, *store)
	test.Ok(t, err)
}
