package main

import (
	"encoding/json"
	"fmt"

	"github.com/doctordesh/yrm"
)

var input = `
// a comment
host: "localhost"
ports:
	http: 8888
	grpc: 9999
startup_delay: 5.5
env: "production"
verbose: true
lorem: false
l:
	t: true
	b:
		c: false
	a: a
`

func main() {
	// res, err := yrm.ParseFile("config.yrm")
	// if err != nil {
	// 	panic(err)
	// }

	res, err := yrm.Parse(input)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	b, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", b)
}
