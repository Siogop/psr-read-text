package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

type TextInImage struct {
	ID           string
	DetectedText string
}

func handler(ctx context.Context, event events.SQSEvent) (bool, error) {
	_session := session.Must(session.NewSession())
	_rekognition := rekognition.New(_session)
	_dynamodb := dynamodb.New(_session)

	var eventIn events.S3Event

	for _, record := range event.Records {
		err := json.Unmarshal([]byte(record.Body), &eventIn)
		if err != nil {
			return false, err
		}

		for _, s3Record := range eventIn.Records {
			textIn := &rekognition.DetectTextInput{
				Image: &rekognition.Image{
					S3Object: &rekognition.S3Object{
						Bucket: aws.String(s3Record.S3.Bucket.Name),
						Name:   aws.String(s3Record.S3.Object.Key),
					},
				},
			}

			textRes, err := _rekognition.DetectText(textIn)
			if err != nil {
				return false, err
			}

			var textOut bytes.Buffer
			for _, text := range textRes.TextDetections {
				textOut.WriteString(*text.DetectedText)
				textOut.WriteString("\n")
			}

			table := os.Getenv("Table")

			item := TextInImage{
				ID:           s3Record.S3.Object.Key,
				DetectedText: textOut.String(),
			}
			fmt.Println(item)
			av, err := dynamodbattribute.MarshalMap(item)
			if err != nil {
				fmt.Println("Got error marshalling new item:")
				fmt.Println(err.Error())
				os.Exit(1)
			}

			input := &dynamodb.PutItemInput{
				Item:      av,
				TableName: aws.String(table),
			}

			_, err = _dynamodb.PutItem(input)
			if err != nil {
				fmt.Println("Got error calling PutItem:")
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}
	}

	return true, nil
}

func main() {
	lambda.Start(handler)
}
