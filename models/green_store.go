package models

import "time"

type GreenStore struct {
	Id          int64
	GroupId     int64
	Code        string
	Name        string
	LogoUrl     string
	QrUrl       string
	Description string
	CreatedAt   time.Time `xorm:"created"`
	UpdatedAt   time.Time `xorm:"updated"`
}
