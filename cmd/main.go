package main

import (
	"go.zakaria.org/vogue"
	"log"
)

func main() {
	log.Printf("starting bruh...")
	brands, err := vogue.GetBrands()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", brands)
}
