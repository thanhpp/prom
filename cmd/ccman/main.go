package main

import (
	"flag"

	"github.com/thanhpp/prom/cmd/ccman/boot"
)

var (
	shardID = flag.Int64("shardID", -1, "Shard id for discovery")
)

func main() {
	flag.Parse()
	if err := boot.Boot(*shardID); err != nil {
		panic(err)
	}
}
