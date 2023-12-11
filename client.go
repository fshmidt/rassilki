package rassilki

import "errors"

type Client struct {
	Id       int    `json:"-"        db:"id"`
	Phone    string `json:"phone"    db:"phone"    binding:"required"`
	Code     string `json:"code"     db:"code"     binding:"required"`
	Tag      string `json:"tag"      db:"tag"`
	Timezone string `json:"timezone" db:"timezone" binding:"required"`
}

type UpdateClient struct {
	Phone    *string `json:"phone"`
	Code     *string `json:"code"`
	Tag      *string `json:"tag"`
	Timezone *string `json:"timezone"`
}

func (i UpdateClient) Validate() error {
	if i.Tag == nil && i.Code == nil && i.Timezone == nil && i.Phone == nil {
		return errors.New("updating structure has no values")
	}
	return nil
}
