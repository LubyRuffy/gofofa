package main

import (
	"log"

	"github.com/lubyruffy/gofofa"
)

func main() {
	cli, err := gofofa.NewClient("")
	if err != nil {
		panic(err)
	}
	log.Println(cli.Account)
}
