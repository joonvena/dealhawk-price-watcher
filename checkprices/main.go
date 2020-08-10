package main

import (
	"context"
	"strings"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

type Product struct {
	ID         string `json:"productId,omitempty"`
	Name       string `json:"name"`
	ImageURL   string `json:"imageURL"`
	Price      string `json:"price"`
	OldPrice   string `json:"old_price"`
	ProductURL string `json:"productURL,omitempty"`
}

type Output struct {
	Products []Product `json:"products"`
	Update   bool      `json:"update"`
}

var wg sync.WaitGroup

var updateProducts []Product

func Contains(a []Product, x string) bool {
	for _, n := range a {
		if x == n.ID {
			return true
		}
	}
	return false
}

func getPrice(c *colly.Collector, url string) string {
	item := Product{}
	c.OnHTML("#main", func(e *colly.HTMLElement) {
		price := strings.Join(strings.Fields(e.ChildText(".price-tag-content__price-tag-price--current")), "")
		item.Price = price
	})
	c.Visit(url)
	return string(item.Price)
}

func checkProducts(c *colly.Collector, product Product) {
	defer wg.Done()
	price := getPrice(c, product.ProductURL)
	if ConvertToNumber(price) < ConvertToNumber(product.Price) {
		product.OldPrice = product.Price
		product.Price = price
		if !Contains(updateProducts, product.ID) {
			updateProducts = append(updateProducts, product)
		}
	}
}

func HandleRequest(ctx context.Context, products []Product) (Output, error) {

	c := colly.NewCollector()
	extensions.RandomUserAgent(c)

	wg.Add(len(products))
	for _, p := range products {
		go checkProducts(c, p)
	}
	wg.Wait()

	result := Output{Products: updateProducts, Update: true}

	if result.Products == nil {
		result.Update = false
	}

	return result, nil
}

func main() {
	lambda.Start(HandleRequest)
}
