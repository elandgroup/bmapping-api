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

func (e GreenStoreApiController) Get(c echo.Context) error {
	status := c.QueryParam("status")
	switch status {
	case "all":

	default:
		return e.GetEId(c)
	}
	return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, nil)
}

func (g GreenStoreApiController) GetEId(c echo.Context) error {
	ipayTypeId := c.QueryParam("ipayTypeId")
	switch ipayTypeId {
	case "1":
		return g.GetEIdOffline(c)
	case "2":
		return g.GetEIdOnline(c)
	}
	return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, nil)
}

func (GreenStoreApiController) GetEIdOffline(c echo.Context) error {
	code := c.QueryParam("code")
	ipayTypeId, err := strconv.ParseInt(c.QueryParam("ipayTypeId"), 10, 64)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	key := fmt.Sprintf("green%v|%v",
		code,
		c.QueryParam("ipayTypeId"),
	)
	currentCache := (*relaxmid.Config(c.Request().Context()))["|cache"].(*cache.Cache)
	if eId, found := currentCache.Get(key); found {
		return ReturnApiSucc(c, http.StatusOK, eId)
	}
	has, eId, _, err := models.GetEIdByCode(c.Request().Context(), code, ipayTypeId)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}
	if !has {
		return ReturnApiFail(c, http.StatusNotFound, ApiErrorNotFound, nil)
	}
	currentCache.Set(key, eId, cache.NoExpiration)
	return ReturnApiSucc(c, http.StatusOK, eId)
}

func (GreenStoreApiController) GetEIdOnline(c echo.Context) error {
	code := c.QueryParam("code")
	ipayTypeId, err := strconv.ParseInt(c.QueryParam("ipayTypeId"), 10, 64)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	key := fmt.Sprintf("green%v|%v",
		code,
		c.QueryParam("ipayTypeId"),
	)

	currentCache := (*relaxmid.Config(c.Request().Context()))["|cache"].(*cache.Cache)
	if store, found := currentCache.Get(key); found {
		return ReturnApiSucc(c, http.StatusOK, store)
	}
	has, eId, store, err := models.GetEIdByCode(c.Request().Context(), code, ipayTypeId)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}
	if !has {
		return ReturnApiFail(c, http.StatusNotFound, ApiErrorNotFound, nil)
	}
	var bizStore = struct {
		*models.GreenStore
		EId int64 `json:"eId"`
	}{
		store,
		eId,
	}
	currentCache.Set(key, &bizStore, cache.NoExpiration)
	return ReturnApiSucc(c, http.StatusOK, bizStore)
}
