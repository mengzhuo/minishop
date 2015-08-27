// Package main provides ...
package main

import (
	"flag"
	"minishop"
)

var (
	address = flag.String("address", ":8765", "")
	db_path = flag.String("db", "localhost", "")
)

func main() {
	flag.Parse()
	minishop.Serve(*address, *db_path)
}
