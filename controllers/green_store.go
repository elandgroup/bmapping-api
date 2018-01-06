package controllers

import (
	"bmapping-api/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
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
	has, eId, err := models.GetEIdByCode(c.Request().Context(), code, ipayTypeId)
	if err != nil {
		return ReturnApiFail(c, http.StatusInternalServerError, ApiErrorDB, err)
	}
	if !has {
		return ReturnApiFail(c, http.StatusNotFound, ApiErrorNotFound, nil)
	}
	return ReturnApiSucc(c, http.StatusOK, eId)
}
