package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Product struct {
	ID         string `json:"productId,omitempty"`
	Name       string `json:"name"`
	ImageURL   string `json:"imageURL"`
	Price      string `json:"price"`
	OldPrice   string `json:"old_price"`
	ProductURL string `json:"productURL,omitempty"`
}

type Input struct {
	Products []Product `json:"products"`
	Foo      int       `json:"foo"`
}

func HandleRequest(ctx context.Context, input Input) ([]Product, error) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	for _, product := range input.Products {
		test := &dynamodb.UpdateItemInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":p": {
					S: aws.String(product.Price),
				},
				":op": {
					S: aws.String(product.OldPrice),
				},
			},
			Key: map[string]*dynamodb.AttributeValue{
				"productId": {
					S: aws.String(product.ID),
				},
			},
			ReturnValues:     aws.String("UPDATED_NEW"),
			UpdateExpression: aws.String("set price = :p, old_price = :op"),
			TableName:        aws.String("products"),
		}

		_, err := svc.UpdateItem(test)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return input.Products, nil
}

func main() {
	lambda.Start(HandleRequest)
}
