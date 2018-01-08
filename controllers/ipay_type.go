package controllers

import (
	"bmapping-api/factory"
	"bmapping-api/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type IPayTypeApiController struct {
}

func (c IPayTypeApiController) Init(g *echo.Group) {
	g.GET("", c.GetAll)
	g.POST("", c.Create)
	g.GET("/:id", c.GetOne)
	g.PUT("/:id", c.Update)
}

func (e IPayTypeApiController) Create(c echo.Context) error {
	status := c.QueryParam("status")
	switch status {
	case "batch":
		return e.InsertMany(c)
	default:
		return e.InsertOne(c)
	}
}

func (IPayTypeApiController) InsertMany(c echo.Context) error {
	v := make([]models.IpayType, 0)
	if err := c.Bind(&v); err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	if err := (models.IpayType{}).InsertMany(c.Request().Context(), &v); err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}

	return ReturnApiSucc(c, http.StatusOK, make([]interface{}, 0))
}

func (IPayTypeApiController) GetAll(c echo.Context) error {
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

	totalCount, items, err := models.IpayType{}.GetAll(c.Request().Context(), v.Sortby, v.Order, v.SkipCount, v.MaxResultCount)
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

func (IPayTypeApiController) InsertOne(c echo.Context) error {
	var v models.IpayType
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

func (IPayTypeApiController) GetOne(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
	has, v, err := models.IpayType{}.GetById(c.Request().Context(), id)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}
	if !has {
		return ReturnApiFail(c, http.StatusNotFound, ApiErrorNotFound, nil)
	}

	return ReturnApiSucc(c, http.StatusOK, v)
}

func (IPayTypeApiController) Update(c echo.Context) error {
	var v models.IpayType
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
