package main

import (
	"fmt"
	"strconv"
	"strings"
)

// ConvertToNumber is used to convert string price coming from scraper to actual number for comparison
func ConvertToNumber(number string) float64 {
	convertedString := strings.Replace(number, ",", ".", -1)
	n, err := strconv.ParseFloat(convertedString, 32)
	if err != nil {
		fmt.Println(err.Error())
	}
	return n
}
