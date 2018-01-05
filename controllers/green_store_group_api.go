package controllers

import (
	"bmapping-api/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type GreenStoreGroupApiController struct {
}

func (c GreenStoreGroupApiController) Init(g *echo.Group) {
	g.GET("", c.GetAll)
}

func (e GreenStoreGroupApiController) GetAll(c echo.Context) error {
	status := c.QueryParam("status")
	var err error
	switch status {
	case "e_id":
		return e.GetEIdByCode(c)
	default:
		return ReturnApiFail(c, http.StatusBadRequest, ApiErrorParameter, err)
	}
}

func (GreenStoreGroupApiController) GetEIdByCode(c echo.Context) error {
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
