package main

import (
	"flag"
	"fmt"
	tplFunc "go-template/src/go/template-functions"
	"go-template/src/go/utils"
	"io"
	"os"

	"dario.cat/mergo"
)

func usage() {
	fmt.Printf("Usage: %s [OPTIONS] template-file\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	valuesFilePath := flag.String("file", "", "Values file path")
	valuesData := flag.String("data", "", "Values data in JSON format")
	outputFilePath := flag.String("output", "", "Render file output")

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	var values map[string]any

	if len(*valuesFilePath) > 0 {
		var err error
		values, err = utils.ReadValuesFile(*valuesFilePath)
		utils.CheckErr(err)
	}

	if len(*valuesData) > 0 {
		argVal := utils.MustUnmarshalValues(flag.Arg(1))
		utils.CheckErr(mergo.Map(&values, argVal))
	}

	var out io.Writer
	if len(*outputFilePath) > 0 {
		fi, err := os.OpenFile(*outputFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		utils.CheckErr(err)
		defer func() {
			utils.CheckErr(fi.Close())
		}()
		out = fi
	} else {
		out = os.Stdout
	}

	utils.CheckErr(tplFunc.Render(flag.Arg(0), out, values))
}
