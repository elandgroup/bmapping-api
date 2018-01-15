package controllers

import (
	"bmapping-api/factory"
	"bmapping-api/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	cache "github.com/patrickmn/go-cache"
	relaxmid "github.com/relax-space/go-kitt/echomiddleware"
	"github.com/sirupsen/logrus"
)

type ElandStoreApiController struct {
}

func (g ElandStoreApiController) Get(c echo.Context) error {
	status := c.QueryParam("status")
	switch status {
	case "all":
		return g.GetAll(c)
	default:
		ipayTypeId := c.QueryParam("ipayTypeId")
		switch ipayTypeId {
		case "1":
			return g.GetEIdOffline(c)
		case "2":
			return g.GetEIdOnline(c)
		case "3":
			return g.GetEIdOnline(c)
		}
	}
	return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, nil)
}
func (ElandStoreApiController) GetAll(c echo.Context) error {
	return nil
}
func (ElandStoreApiController) GetEIdOffline(c echo.Context) error {
	code := c.QueryParam("code")
	group_code := c.QueryParam("groupCode")
	country_id, err := strconv.ParseInt(c.QueryParam("countryId"), 10, 64)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	ipayTypeId, err := strconv.ParseInt(c.QueryParam("ipayTypeId"), 10, 64)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	key := fmt.Sprintf("eland%v|%v|%v|%v",
		code,
		group_code,
		c.QueryParam("countryId"),
		c.QueryParam("ipayTypeId"),
	)
	currentCache := (*relaxmid.Config(c.Request().Context()))["|cache"].(*cache.Cache)
	if eId, found := currentCache.Get(key); found {
		return ReturnApiSucc(c, http.StatusOK, eId)
	}
	has, eId, _, err := models.GetEIdByThrArgs(c.Request().Context(), group_code, code, country_id, ipayTypeId)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}
	if !has {
		return ReturnApiFail(c, http.StatusNotFound, ApiErrorNotFound, nil)
	}
	currentCache.Set(key, eId, cache.NoExpiration)
	return ReturnApiSucc(c, http.StatusOK, eId)
}

func (ElandStoreApiController) GetEIdOnline(c echo.Context) error {
	code := c.QueryParam("code")
	group_code := c.QueryParam("groupCode")
	country_id, err := strconv.ParseInt(c.QueryParam("countryId"), 10, 64)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	ipayTypeId, err := strconv.ParseInt(c.QueryParam("ipayTypeId"), 10, 64)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	key := fmt.Sprintf("eland%v|%v|%v|%v",
		code,
		group_code,
		c.QueryParam("countryId"),
		c.QueryParam("ipayTypeId"),
	)
	currentCache := (*relaxmid.Config(c.Request().Context()))["|cache"].(*cache.Cache)
	if eId, found := currentCache.Get(key); found {
		return ReturnApiSucc(c, http.StatusOK, eId)
	}
	has, eId, store, err := models.GetEIdByThrArgs(c.Request().Context(), group_code, code, country_id, ipayTypeId)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}
	if !has {
		return ReturnApiFail(c, http.StatusNotFound, ApiErrorNotFound, nil)
	}
	var bizStore = struct {
		*models.ElandStore
		EId int64 `json:"eId"`
	}{
		store,
		eId,
	}
	currentCache.Set(key, bizStore, cache.NoExpiration)
	return ReturnApiSucc(c, http.StatusOK, bizStore)
}

type ElandStoreGroupApiController struct {
}

func (c ElandStoreGroupApiController) Init(g *echo.Group) {
	g.GET("", c.GetAll)
	g.POST("", c.Create)
	g.GET("/:id", c.GetOne)
	g.PUT("/:id", c.Update)
}

func (e ElandStoreGroupApiController) Create(c echo.Context) error {
	status := c.QueryParam("status")
	switch status {
	case "batch":
		return e.InsertMany(c)
	default:
		return e.InsertOne(c)
	}
}

func (ElandStoreGroupApiController) InsertMany(c echo.Context) error {
	v := make([]models.ElandStoreGroup, 0)
	if err := c.Bind(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	if err := (models.ElandStoreGroup{}).InsertMany(c.Request().Context(), &v); err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}

	return ReturnApiSucc(c, http.StatusOK, make([]interface{}, 0))
}

func (ElandStoreGroupApiController) GetAll(c echo.Context) error {
	var v SearchInput
	if err := c.Bind(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	if v.MaxResultCount == 0 {
		v.MaxResultCount = DefaultMaxResultCount
	}

	factory.BehaviorLogger(c.Request().Context()).WithBizAttr("maxResultCount", v.MaxResultCount).Log("SearchElandStore")

	factory.Logger(c.Request().Context()).WithFields(logrus.Fields{
		"sortby":         v.Sortby,
		"order":          v.Order,
		"maxResultCount": v.MaxResultCount,
		"skipCount":      v.SkipCount,
	}).Info("SearchInput")

	totalCount, items, err := models.ElandStoreGroup{}.GetAll(c.Request().Context(), v.Sortby, v.Order, v.SkipCount, v.MaxResultCount)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}

	factory.BehaviorLogger(c.Request().Context()).
		WithCallURLInfo(http.MethodGet, "https://play.google.com/books", nil, 200).
		WithBizAttrs(map[string]interface{}{
			"totalCount": totalCount,
			"itemCount":  len(items),
		}).
		Log("SearchComplete")

	return ReturnApiSucc(c, http.StatusOK, ArrayResult{
		TotalCount: totalCount,
		Items:      items,
	})
}

func (ElandStoreGroupApiController) InsertOne(c echo.Context) error {
	var v models.ElandStoreGroup
	if err := c.Bind(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	if err := c.Validate(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	if err := (&v).InsertOne(c.Request().Context()); err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}

	return ReturnApiSucc(c, http.StatusOK, v)
}

func (ElandStoreGroupApiController) GetOne(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	has, v, err := models.ElandStoreGroup{}.GetById(c.Request().Context(), id)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}
	if !has {
		return ReturnApiFail(c, http.StatusNotFound, ApiErrorNotFound, nil)
	}

	return ReturnApiSucc(c, http.StatusOK, v)
}

func (ElandStoreGroupApiController) Update(c echo.Context) error {
	var v models.ElandStoreGroup
	if err := c.Bind(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	if err := c.Validate(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	v.Id = id
	if err := (&v).Update(c.Request().Context()); err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}

	return ReturnApiSucc(c, http.StatusOK, v)
}
