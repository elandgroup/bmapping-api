package controllers

import (
	"bmapping-api/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	cache "github.com/patrickmn/go-cache"
	relaxmid "github.com/relax-space/go-kitt/echomiddleware"
)

type GreenStoreApiController struct {
}

func (c GreenStoreApiController) Init(g *echo.Group) {
	g.GET("", c.GetAll)
}

func (e GreenStoreApiController) GetAll(c echo.Context) error {
	status := c.QueryParam("status")
	switch status {
	case "all":

	default:
		return e.GetEIdByCode(c)
	}
	return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, nil)
}

func (GreenStoreApiController) GetEIdByCode(c echo.Context) error {
	code := c.QueryParam("code")
	ipayTypeId, err := strconv.ParseInt(c.QueryParam("ipay_type_id"), 10, 64)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	key := fmt.Sprintf("green%v|%v",
		code,
		c.QueryParam("ipay_type_id"),
	)
	currentCache := (*relaxmid.Config(c.Request().Context()))["|cache"].(*cache.Cache)
	if eId, found := currentCache.Get(key); found {
		return ReturnApiSucc(c, http.StatusOK, eId)
	}
	has, eId, err := models.GetEIdByCode(c.Request().Context(), code, ipayTypeId)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}
	if !has {
		return ReturnApiFail(c, http.StatusNotFound, ApiErrorNotFound, nil)
	}
	currentCache.Set(key, eId, cache.NoExpiration)
	return ReturnApiSucc(c, http.StatusOK, eId)
}
