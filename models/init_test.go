package models

import (
	"context"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/pangpanglabs/goutils/echomiddleware"
)

var ctx context.Context

func init() {
	xormEngine, err := xorm.NewEngine("mysql", os.Getenv("BMAPPING_ENV"))
	if err != nil {
		panic(err)
	}
	xormEngine.ShowSQL(true)
	ctx = context.WithValue(context.Background(), echomiddleware.ContextDBName, xormEngine.NewSession())
}
