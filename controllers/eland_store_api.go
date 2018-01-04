package controllers

import (
	"net/http"
	"strconv"

	"bmapping-api/factory"
	"bmapping-api/models"

	"github.com/labstack/echo"
    "github.com/sirupsen/logrus"
)

type ElandStoreApiController struct {
}

func (c ElandStoreApiController) Init(g *echo.Group) {
	g.GET("", c.GetAll)
	g.POST("", c.Create)
	g.GET("/:id", c.GetOne)
	g.PUT("/:id", c.Update)
}
func (ElandStoreApiController) GetAll(c echo.Context) error {
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

	totalCount, items, err := models.ElandStore{}.GetAll(c.Request().Context(), v.Sortby, v.Order, v.SkipCount, v.MaxResultCount)
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

func (ElandStoreApiController) Create(c echo.Context) error {
	var v ElandStoreInput
	if err := c.Bind(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	if err := c.Validate(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	eland_store, err := v.ToModel()
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	if _, err := eland_store.Create(c.Request().Context()); err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}

	return ReturnApiSucc(c, http.StatusOK, eland_store)
}

func (ElandStoreApiController) GetOne(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	v, err := models.ElandStore{}.GetById(c.Request().Context(), id)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}
	if v == nil {
		return ReturnApiFail(c, http.StatusNotFound, ApiErrorNotFound, nil)
	}

	return ReturnApiSucc(c, http.StatusOK, v)
}

func (ElandStoreApiController) Update(c echo.Context) error {
	var v ElandStoreInput
	if err := c.Bind(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	if err := c.Validate(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	eland_store, err := v.ToModel()
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	eland_store.Id = id
	if err := eland_store.Update(c.Request().Context()); err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}

	return ReturnApiSucc(c, http.StatusOK, eland_store)
}
