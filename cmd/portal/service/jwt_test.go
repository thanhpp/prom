package service_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/thanhpp/prom/cmd/portal/service"

	"github.com/dgrijalva/jwt-go"
)

func TestGenerate(t *testing.T) {
	newJwtClaims := jwt.MapClaims{
		"userid":     0,
		"authorized": true,
		"exp":        time.Now().Add(time.Minute * 15).Unix(),
	}

	newJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, newJwtClaims)
	token, err := newJwt.SignedString([]byte("myjwtkey"))
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(token)
}

func TestValidate(t *testing.T) {
	var (
		tokenStr string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MTk5NTMyOTQsInVzZXJpZCI6MH0.D52AL60DZNociAUbtWLcFjwRZ7sQtBF6IG0P7oCcbJY"
	)

	// validate token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// validate alg
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected sign method: %v", token.Header["alg"])
		}

		// if validate ok, then return secret key
		return []byte("myjwtkey"), nil
	})
	if err != nil {
		t.Error(err)
		return
	}

	// get claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		fmt.Println(claims["userid"])
		return
	}

	t.Error("invalid token")
}

func TestToJWTClaims(t *testing.T) {
	var (
		iat      = time.Now()
		claimsDt = &service.ClaimsDetail{
			UUID:      "123",
			UserID:    0,
			IssuedAt:  iat.Unix(),
			ExpiredAt: iat.Add(time.Second * 10).Unix(),
		}
	)

	m, err := claimsDt.ToJWTClaims()
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(m)
}

func TestJWTService(t *testing.T) {
	var (
		jwtSrv = &service.JWTSrv{
			SigningMethod: jwt.SigningMethodHS256,
			SecretKey:     []byte("myjwtsecretkey"),
		}
		iat          = time.Now()
		claimsDetail = &service.ClaimsDetail{
			UUID:      "uuid",
			UserID:    0,
			IssuedAt:  iat.Unix(),
			ExpiredAt: iat.Add(time.Second * 10).Unix(),
		}
	)

	token, err := jwtSrv.CreateToken(claimsDetail)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(token)

	claims, err := jwtSrv.GetClaimsValidate(token)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(claims)
}
