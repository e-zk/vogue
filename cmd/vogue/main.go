package main

import (
	"flag"
	"fmt"
	"log"

	"go.zakaria.org/vogue"
)

var (
	lsSeasons bool
	lsBrands  bool
	fullSlug  string
)

func init() {
	flag.BoolVar(&lsSeasons, "seasons", false, "list all seasons")
	flag.BoolVar(&lsBrands, "brands", false, "list all brands")
	flag.StringVar(&fullSlug, "show", "", "print all collection urls from show")
	flag.Parse()
}

func main() {
	if fullSlug != "" {
		ss, err := vogue.GetShow(fullSlug)
		if err != nil {
			log.Fatal(err)
		}

		for _, slide := range ss.Galleries.Collection.Slides.Slides {
			fmt.Printf("%s\n", slide.Image.Url)
		}
		return
	}

	if lsSeasons {
		ssns, err := vogue.GetSeasons()
		if err != nil {
			log.Fatal(err)
		}

		for _, s := range ssns {
			fmt.Printf("%s\t%s\n", s.Slug, s.Name)
		}
	} else if lsBrands {
		brands, err := vogue.GetBrands()
		if err != nil {
			log.Fatal(err)
		}

		for _, b := range brands {
			fmt.Printf("%s\t%s\n", b.Slug, b.Name)
		}
	} else {
		log.Fatal("huh? what?")
	}
}
