package main

import (
	"encoding/csv"
	"flag"
	"os"
	"strings"

	"github.com/lubyruffy/gofofa"
)

func main() {
	configURL := flag.String("fofaURL", gofofa.FofaURLFromEnv(), "format: <url>/?email=<email>&key=<key>&version=<v2>")
	deductMode := flag.String("deductMode", "DeductModeFree", "DeductModeFree or DeductModeFCoin")
	fieldString := flag.String("fields", "ip,port", "more info, visit fofa website")
	flag.Parse()

	cli, err := gofofa.NewClient(*configURL)
	if err != nil {
		panic(err)
	}
	if len(*deductMode) > 0 {
		cli.DeductMode = gofofa.ParseDeductMode(*deductMode)
	}

	//log.Println(cli.Account)

	res, err := cli.HostSearch("port=80", 100, strings.Split(*fieldString, ","))
	if err != nil {
		panic(err)
	}
	writer := csv.NewWriter(os.Stdout)
	if err = writer.WriteAll(res); err != nil {
		panic(err)
	}
}
