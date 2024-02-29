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

var Headers map[string][]string = map[string][]string{
	"Host":         []string{"graphql.vogue.com"},
	"User-Agent":   []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"},
	"Content-Type": []string{"application/json"},
	// "Referrer-Policy": []string{"origin"},
}

type Brand struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type brandsResp struct {
	Data struct {
		AllBrands struct {
			Brands []Brand `json:"Brand"`
		} `json:"allBrands"`
	} `json:"data"`
}

func graphQ(query string) ([]byte, error) {
	c := http.Client{
		Timeout: time.Second * 5,
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
	b, err := graphQ("query{allBrands{Brand{name slug}}}")
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
