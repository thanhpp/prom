package service_test

import (
	"fmt"
	"testing"
	"time"

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
		token string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MTk5NDI0MTgsInVzZXJpZCI6MH0.bzT6tycENcckikQj6rmn5RiAKvfu1sqWRPZob1CYrt0"
	)
	
	return
}