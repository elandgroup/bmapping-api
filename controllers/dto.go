package controllers

import (
	"bmapping-api/models"
)

const (
	DefaultMaxResultCount = 30
)

type SearchInput struct {
	Sortby         []string `query:"sortby"`
	Order          []string `query:"order"`
	SkipCount      int      `query:"skipCount"`
	MaxResultCount int      `query:"maxResultCount"`
}

type ElandStoreInput struct {
	Id	int64     `json:"id"`
	GroupCode	string     `json:"groupCode"`
	Code	string     `json:"code"`
	CountryId	int64     `json:"countryId"`
	Name	string     `json:"name"`
	CreatedAt	string     `json:"createdAt"`
	UpdatedAt	string     `json:"updatedAt"`
}


func (d *ElandStoreInput) ToModel() (*models.ElandStore, error) {
	return &models.ElandStore{
		Id:      d.Id,
		GroupCode:      d.GroupCode,
		Code:      d.Code,
		CountryId:      d.CountryId,
		Name:      d.Name,
		CreatedAt:      d.CreatedAt,
		UpdatedAt:      d.UpdatedAt,
	}, nil
}

