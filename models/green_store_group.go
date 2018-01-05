package models

type GreenStoreGroup struct {
	Id        int64  `json:"id" query:"id" xorm:"pk autoincr 'id'"`
	Code      string `json:"code" query:"code" xorm:"VARCHAR(4) unique 'code'"`
	Name      string `json:"name" query:"name"`
	CreatedAt string `json:"createdAt" query:"createdAt" xorm:"created"`
	UpdatedAt string `json:"updatedAt" query:"updatedAt" xorm:"updated"`
}
