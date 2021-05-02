package service

import (
	"encoding/json"

	"github.com/dgrijalva/jwt-go"
)

type ClaimsDetail struct {
	UUID      string `json:"uuid"`
	UserID    uint32 `json:"userid"`
	IssuedAt  int64  `json:"iat"`
	ExpiredAt int64  `json:"exp"`
}

func (c *ClaimsDetail) ToJWTClaims() (claims *jwt.MapClaims, err error) {
	var data = new(jwt.MapClaims)

	dataB, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(dataB, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (c *ClaimsDetail) ImportData(claims *jwt.MapClaims) (err error) {
	dataB, err := json.Marshal(claims)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(dataB, c); err != nil {
		return err
	}

	return nil
}
