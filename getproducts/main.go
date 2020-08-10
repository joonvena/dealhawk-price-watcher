package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws/session"
)

type MyEvent struct {
	Name string `json:"name"`
}

type Product struct {
	ID         string `json:"productId,omitempty"`
	Name       string `json:"name"`
	ImageURL   string `json:"imageURL"`
	Price      string `json:"price"`
	ProductURL string `json:"productURL,omitempty"`
}

type Data struct {
	UserID   string    `json:"userId"`
	Products []Product `json:"products"`
}

// HandleRequest is a function that get executed when Lambda is run.
func HandleRequest(ctx context.Context) ([]Product, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	result, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("products_table")),
	})

	if err != nil {
		fmt.Println(err.Error())
	}

	var product []Product

	if err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &product); err != nil {
		fmt.Println(err.Error())
	}

	return product, nil

}

func main() {
	lambda.Start(HandleRequest)
}
