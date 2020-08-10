package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/service/ses"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
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

type User struct {
	ID       string   `json:"userId"`
	Email    string   `json:"email"`
	Products []string `json:"products"`
}

type Response struct {
	Products []Product `json:"products"`
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func ComposeMail(response Response) (mailBody *ses.SendEmailInput) {

	var tpl bytes.Buffer

	t := template.Must(template.New("mail-tmpl").Parse(`
		<h1>These product prices have changed: </h1>
		<div>
		{{ range .Products }}
			<div style="display: flex; align-items: center;">
				<div style="flex: 1;">
				<img src="{{.ImageURL}}" alt="{{.Name}}" style="max-width: 100%">
				</div>
				<div style="flex: 3; margin-left: 10px; margin-right: 5px;">
				<a href="{{.ProductURL}}">{{.Name}}</a>
				</div>
				<div>
				<h4 style="margin: 0; padding: 0;">{{.Price}}</h4>
				<h4 style="margin: 0; padding: 0; float: right; text-decoration: line-through; color: rgb(160, 160, 160);">{{.OldPrice}}</h4>
				</div>
			</div>
			
		{{end}}
		<div>
	`))

	err := t.Execute(&tpl, response)

	if err != nil {
		log.Println(err.Error())
	}

	mailMessage := fmt.Sprintf(`%v`, tpl.String())

	mailBody = &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(os.Getenv("email")),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(mailMessage),
				},
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String("This is email send from AWS SES."),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String("Product prices changed"),
			},
		},
		Source: aws.String("joonas@gingerdev.net"),
	}
	return
}

func HandleRequest(ctx context.Context, input []Product) (string, error) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	emailsvc := ses.New(sess)

	var response Response

	response.Products = input

	re, err := emailsvc.SendEmail(ComposeMail(response))

	if err != nil {
		log.Println(err.Error())
	}

	log.Println(re)

	return "Notify send succesfully", nil
}

func main() {
	lambda.Start(HandleRequest)
}
