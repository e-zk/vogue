package vogue

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	baseUrl = "https://graphql.vogue.com"
)

var Timeout time.Duration = time.Second * 6
var Headers map[string][]string = map[string][]string{
	"Host":         []string{"graphql.vogue.com"},
	"User-Agent":   []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"},
	"Content-Type": []string{"application/json"},
}

type Brand struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Season struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
	Year int    `json:"year,omitempty"`
}

type brandsResp struct {
	Data struct {
		AllBrands struct {
			Brands []Brand `json:"Brand"`
		} `json:"allBrands"`
	} `json:"data"`
}

type seasonsResp struct {
	Data struct {
		AllSeasons struct {
			Seasons []Season `json:"Season"`
		} `json:"allSeasons"`
	} `json:"data"`
}

type ShowImage struct {
	Id      string `json:"id,omitempty"`
	Url     string `json:"url"`
	Caption string `json:"caption,omitempty"`
	Credit  string `json:"title,omitempty"`
	Width   int    `json:"width,omitempty"`
	Height  int    `json:"height,omitempty"`
}

type ShowSlide struct {
	Id       string    `json:"id"`
	Type     string    `json:"type"`
	Title    string    `json:"title,omitempty"`
	Image    ShowImage `json:"photosTout"`
	TypeName string    `json:"__typename,omitempty"`
}

type ShowGallery struct {
	Title  string `json:"string"`
	Slides struct {
		Slides []ShowSlide `json:"slide"`
	} `json:"slidesV2"`
}

// could this be a map?
type ShowGalleries struct {
	Collection *ShowGallery `json:"collection,omitempty"`
	Atmosphere *ShowGallery `json:"atmosphere,omitempty"`
	Beauty     *ShowGallery `json:"beauty,omitempty"`
	Detail     *ShowGallery `json:"detail,omitempty"`
	FrontRow   *ShowGallery `json:"frontRow,omitempty"`
}

// fashionShowV2
type Show struct {
	PublishedGMT time.Time         `json:"GMTPubDate"`
	Url          string            `json:"url"`
	Title        string            `json:"title"`
	FullSlug     string            `json:"slug"`
	Id           string            `json:"id"`
	City         map[string]string `json:"city"`
	Brand        Brand             `json:"brand"`
	Season       Season            `json:"season"`
	HeroImage    ShowImage         `json:"photosTout"`
	Galleries    ShowGalleries     `json:"galleries"`
	Video        *interface{}      `json:"video,omitempty"` // TODO
}

type showResp struct {
	Data struct {
		FashionShowV2 Show `json:"fashionShowV2"`
	} `json:"data"`
}

func graphQ(query string) ([]byte, error) {
	c := http.Client{
		Timeout: Timeout,
	}

	// base url
	reqUrl, _ := url.JoinPath(baseUrl, "graphql")

	// set query param (url.Values.Encode does not work properly for graphql query over http)
	reqUrl = reqUrl + "?query=" + url.PathEscape(query)

	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return []byte{}, err
	}

	// set headers
	req.Header = Headers

	// exec request
	resp, err := c.Do(req)
	if err != nil {
		return []byte{}, err
	}
	if resp.Body == nil {
		return []byte{}, fmt.Errorf("no data returned")
	}

	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	fmt.Printf("respb: %s\n", respB)

	return respB, nil
}

func GetBrands() ([]Brand, error) {
	b, err := graphQ(qBrands)
	if err != nil {
		return []Brand{}, err
	}

	bs := brandsResp{}
	err = json.Unmarshal(b, &bs)
	if err != nil {
		return []Brand{}, err
	}

	return bs.Data.AllBrands.Brands, nil
}

func GetSeasons() ([]Season, error) {
	b, err := graphQ(qSeasons)
	if err != nil {
		return []Season{}, err
	}

	ss := seasonsResp{}
	err = json.Unmarshal(b, &ss)
	if err != nil {
		return []Season{}, err
	}

	return ss.Data.AllSeasons.Seasons, nil
}

// fashionShowV2
// fullSlug = {season-slug}/{brand-slug}
func GetShow(fullSlug string) (Show, error) {
	b, err := graphQ(fmt.Sprintf(qFashionShow, fullSlug))
	if err != nil {
		return Show{}, err
	}

	fs := showResp{}
	err = json.Unmarshal(b, &fs)
	if err != nil {
		return Show{}, err
	}

	return fs.Data.FashionShowV2, nil
}
