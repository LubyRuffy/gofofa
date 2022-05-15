package main

import (
	"flag"
	"log"

	"github.com/lubyruffy/gofofa"
)

func main() {
	configURL := flag.String("fofaURL", gofofa.FofaURLFromEnv(), "format: <url>/?email=<email>&key=<key>&version=<v2>")
	flag.Parse()

	cli, err := gofofa.NewClient(*configURL)
	if err != nil {
		panic(err)
	}
	log.Println(cli.Account)
}
