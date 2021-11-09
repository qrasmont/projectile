package main

import (
	"flag"
	"fmt"
)

func main() {
	path := flag.String("p", "", "The project's path.")
	flag.Parse()

	action := flag.Args()

	fmt.Printf("path: %s, action: %s\n", *path, action)
}
