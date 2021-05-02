package service

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// -------------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- JWT SERVICE ----------------------------------------------------------

type JWTSrv struct {
	SigningMethod jwt.SigningMethod
	SecretKey     []byte
}

type iJWTSrv interface {
	CreateToken(claims *ClaimsDetail) (tokenStr string, err error)
	GetClaimsValidate(tokenStr string) (claims *ClaimsDetail, err error)
}

var _ iJWTSrv = (*JWTSrv)(nil)

// -----------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- FUNCTIONS ----------------------------------------------------------

func (j JWTSrv) CreateToken(claims *ClaimsDetail) (tokenStr string, err error) {
	claimsDt, err := claims.ToJWTClaims()
	if err != nil {
		return "", err
	}

	newJWT := jwt.NewWithClaims(j.SigningMethod, claimsDt)
	tokenStr, err = newJWT.SignedString(j.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (j JWTSrv) GetClaimsValidate(tokenStr string) (claims *ClaimsDetail, err error) {
	// validate token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// validate alg
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected sign method: %v", token.Header["alg"])
		}

		// if validate ok, then return secret key
		return j.SecretKey, nil
	})
	if err != nil {
		return nil, err
	}

	// get claims
	claimsJWT, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("Invalid token")
	}

	claims = new(ClaimsDetail)
	if err = claims.ImportData(&claimsJWT); err != nil {
		return nil, err
	}

	return claims, nil
}
