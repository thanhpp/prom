package main

import (
	"github.com/thanhpp/prom/cmd/ccman/boot"
)

func main() {
	if err := boot.Boot(); err != nil {
		panic(err)
	}
}
