package controllers

import (
	"net/http"

	"github.com/labstack/echo"
)

type TenancyStoreApiController struct {
}

func (c TenancyStoreApiController) Init(g *echo.Group) {
	g.GET("", c.Get)
}

func (g TenancyStoreApiController) Get(c echo.Context) error {
	tenancy := c.QueryParam("tenancy")
	switch tenancy {
	case "green":
		return GreenStoreApiController{}.Get(c)
	case "eland":
		return ElandStoreApiController{}.Get(c)
	}
	return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, nil)
}
