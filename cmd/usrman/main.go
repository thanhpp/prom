package main

import (
	"github.com/thanhpp/prom/cmd/usrman/boot"
)

func main() {
	if err := boot.Boot(); err != nil {
		panic(err)
	}
}
