package main

import (
	"encoding/json"
	"fmt"

	"github.com/doctordesh/yrm"
)

func main() {
	res, err := yrm.ParseFile("config.yrm")
	if err != nil {
		panic(err)
	}

	b, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", b)
}
