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

// Headers to pass to http.Client when executing requests to the Vogue GraphQL API.
var Headers map[string][]string = map[string][]string{
	"Host":         []string{"graphql.vogue.com"},
	"User-Agent":   []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"},
	"Content-Type": []string{"application/json"},
}

// Timeout to pass to http.Client when executing requests.
// Defaults to 6 seconds.
var Timeout time.Duration = time.Second * 6

// Brand describes a brand.
type Brand struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// Season describes a fashion season.
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

// ShowGalleries holds the each gallery for a fashion show.
type ShowGalleries struct {
	Collection *ShowGallery `json:"collection,omitempty"`
	Atmosphere *ShowGallery `json:"atmosphere,omitempty"`
	Beauty     *ShowGallery `json:"beauty,omitempty"`
	Detail     *ShowGallery `json:"detail,omitempty"`
	FrontRow   *ShowGallery `json:"frontRow,omitempty"`
}

// Show represents a full fashion show.
type Show struct {
	PublishedGMT time.Time         `json:"GMTPubDate"`
	Url          string            `json:"url"`
	Title        string            `json:"title"`
	FullSlug     string            `json:"slug"`
	Id           string            `json:"id"`
	City         map[string]string `json:"city,omitempty"`
	Brand        Brand             `json:"brand"`
	Season       Season            `json:"season"`
	HeroImage    ShowImage         `json:"photosTout"`
	Galleries    *ShowGalleries    `json:"galleries,omitempty"`
	Video        *interface{}      `json:"video,omitempty"` // TODO
}

type showResp struct {
	Data struct {
		FashionShowV2 Show `json:"fashionShowV2"`
	} `json:"data"`
}

type contentResp struct {
	Data struct {
		AllContent struct {
			Content []Show `json:"Content"`
		} `json:"allContent"`
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

	return respB, nil
}

// GetBrands returns a list of Brand structs.
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

// GetSeasons returns a list of Season structs.
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

// GetShow returns the Show struct for the given fashion show.
// Runway/fashion shows are distinguished by their fullSlug this is of the format: '{season-slug}/{brand-slug}'.
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

// List all shows for given brand.
// Does not return show galleries. Use GetShow() for this.
func GetBrandShows(brandSlug string) ([]Show, error) {
	b, err := graphQ(fmt.Sprintf(qBrandShows, brandSlug))
	if err != nil {
		return []Show{}, err
	}

	c := contentResp{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		return []Show{}, err
	}

	return c.Data.AllContent.Content, nil
}

// List all shows for a given season.
// Does not return show galleries. Use GetShow() for this.
func GetSeasonShows(seasonSlug string) ([]Show, error) {
	b, err := graphQ(fmt.Sprintf(qSeasonShows, seasonSlug))
	if err != nil {
		return []Show{}, err
	}

	c := contentResp{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		return []Show{}, err
	}

	return c.Data.AllContent.Content, nil
}
