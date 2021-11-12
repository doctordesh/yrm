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
production: false
`

func main() {
	res, err := yrm.Parse(input)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	b, _ := json.MarshalIndent(res, "", "    ")

	fmt.Printf("%s\n", b)
}
