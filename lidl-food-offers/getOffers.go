package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Offer struct {
	name       string
	price      string
	offerPrice string
	discount   string
}

type Offers []Offer

func (e Offers) Len() int {
	return len(e)
}

func (e Offers) Less(i, j int) bool {
	return e[i].discount > e[j].discount
}

func (e Offers) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func getLinkToWeekOffer() (string, error) {

	// Request the HTML page.
	res, err := http.Get("https://www.lidl.co.uk/food-offers")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	// Find the link

	links := doc.Find("a.navigation__link").Nodes

	for _, link := range links {
		for _, attr := range link.Attr {
			if attr.Key == "href" && strings.HasPrefix(attr.Val, "/c/pick-of-the-week/") {
				return "https://www.lidl.co.uk" + attr.Val, nil
			}
		}
	}

	return "", fmt.Errorf("link not exist")
}

func getOffers() (Offers, error) {

	url, err := getLinkToWeekOffer()
	if err != nil {
		return nil, err
	}

	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	// Find the link
	offers := Offers{}

	doc.Find("article.ret-o-card").Each(func(_ int, s *goquery.Selection) {
		o := Offer{}
		if s.Find("span.lidl-m-pricebox__discount-price").Text() != "" {
			o.name = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(s.Find("h3.ret-o-card__headline").Text()), "2"))
			priceString := strings.TrimPrefix(s.Find("span.lidl-m-pricebox__discount-price").Text(), "£ ")
			price, _ := strconv.ParseFloat(priceString, 64)
			o.price = priceString
			offerPriceString := strings.TrimPrefix(s.Find("span.lidl-m-pricebox__price").Text(), "£ ")
			offerPrice, _ := strconv.ParseFloat(offerPriceString, 64)
			o.offerPrice = offerPriceString
			discountFloat := ((price - offerPrice) / price) * 100
			o.discount = strconv.Itoa(int(discountFloat))
			offers = append(offers, o)
		}
	})
	sort.Sort(Offers(offers))
	return offers, nil
}
