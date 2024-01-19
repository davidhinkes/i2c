package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/davidhinkes/i2c"
)

func handleErr(e error) {
	if e == nil {
		return
	}
	log.Fatal(e)
}

func main() {
	flag.Parse()
	fmt.Println("program example, copywrite 2024")
	i, err := i2c.Make("/dev/i2c-1")
	handleErr(err)
	ret := make([]byte, 4)
	err = i.Read(ret, 0x1e, 0)
	handleErr(err)
	fmt.Printf("Got %v bytes: %x", len(ret), ret)
}
