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

type GreenStoreApiController struct {
}

func (c GreenStoreApiController) Init(g *echo.Group) {
	g.GET("", c.Get)
	g.GET("", c.GetAll)
	g.POST("", c.Create)
	g.GET("/:id", c.GetOne)
	g.PUT("/:id", c.Update)
}

func (e GreenStoreApiController) Create(c echo.Context) error {
	status := c.QueryParam("status")
	switch status {
	case "batch":
		return e.InsertMany(c)
	default:
		return e.InsertOne(c)
	}
}
func (e GreenStoreApiController) GetAll(c echo.Context) error {
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

	totalCount, items, err := models.GreenStore{}.GetAll(c.Request().Context(), v.Sortby, v.Order, v.SkipCount, v.MaxResultCount)
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

func (GreenStoreApiController) InsertOne(c echo.Context) error {
	var v models.GreenStore
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

func (GreenStoreApiController) InsertMany(c echo.Context) error {
	v := make([]models.GreenStore, 0)
	if err := c.Bind(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	if err := (models.GreenStore{}).InsertMany(c.Request().Context(), &v); err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}

	return ReturnApiSucc(c, http.StatusOK, make([]interface{}, 0))
}

func (GreenStoreApiController) GetOne(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	has, v, err := models.GreenStoreGroup{}.GetById(c.Request().Context(), id)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}
	if !has {
		return ReturnApiFail(c, http.StatusNotFound, ApiErrorNotFound, nil)
	}

	return ReturnApiSucc(c, http.StatusOK, v)
}

func (GreenStoreApiController) Update(c echo.Context) error {
	var v models.GreenStoreGroup
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
	if err := (&v).Update(c.Request().Context(), id); err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}

	return ReturnApiSucc(c, http.StatusOK, v)
}

func (GreenStoreApiController) Delete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	err = models.GreenStore{}.Delete(c.Request().Context(), id)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorParameter, err)
	}
	return ReturnApiSucc(c, http.StatusOK, nil)
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

//green store group

type GreenStoreGroupApiController struct {
}

func (c GreenStoreGroupApiController) Init(g *echo.Group) {
	g.GET("", c.GetAll)
	g.POST("", c.Create)
	g.GET("/:id", c.GetOne)
	g.PUT("/:id", c.Update)
}

func (e GreenStoreGroupApiController) Create(c echo.Context) error {
	status := c.QueryParam("status")
	switch status {
	case "batch":
		return e.InsertMany(c)
	default:
		return e.InsertOne(c)
	}
}

func (GreenStoreGroupApiController) InsertMany(c echo.Context) error {
	v := make([]models.GreenStoreGroup, 0)
	if err := c.Bind(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	if err := (models.GreenStoreGroup{}).InsertMany(c.Request().Context(), &v); err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}

	return ReturnApiSucc(c, http.StatusOK, make([]interface{}, 0))
}

func (GreenStoreGroupApiController) GetAll(c echo.Context) error {
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

	totalCount, items, err := models.GreenStoreGroup{}.GetAll(c.Request().Context(), v.Sortby, v.Order, v.SkipCount, v.MaxResultCount)
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

func (GreenStoreGroupApiController) InsertOne(c echo.Context) error {
	var v models.GreenStoreGroup
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

func (GreenStoreGroupApiController) GetOne(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	has, v, err := models.GreenStoreGroup{}.GetById(c.Request().Context(), id)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}
	if !has {
		return ReturnApiFail(c, http.StatusNotFound, ApiErrorNotFound, nil)
	}

	return ReturnApiSucc(c, http.StatusOK, v)
}

func (GreenStoreGroupApiController) Update(c echo.Context) error {
	var v models.GreenStoreGroup
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
	if err := (&v).Update(c.Request().Context(), id); err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}

	return ReturnApiSucc(c, http.StatusOK, v)
}
