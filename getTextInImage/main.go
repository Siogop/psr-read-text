package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/google/uuid"
)

type RequestJSON struct {
	ImageBase64 string `json:"imageBase64"`
}

type ResponseJSON struct {
	URL string `json:"url"`
}

func handler(ctx context.Context, reqRaw events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	_session := session.Must(session.NewSession())
	_rekognition := rekognition.New(_session)
	_uploader := s3manager.NewUploader(_session)
	bucket := os.Getenv("Bucket")
	uidGen := uuid.New()
	uid := uidGen.String()
	key := fmt.Sprintf("%s.png", uid)

	var headers map[string]string
	headers = make(map[string]string)

	headers["Access-Control-Allow-Origin"] = "*"
	headers["Access-Control-Allow-Credentials"] = "true"

	var reqJSON RequestJSON
	err := json.Unmarshal([]byte(reqRaw.Body), &reqJSON)
	if err != nil {
		fmt.Println("error unmarshalling request")
		return events.APIGatewayProxyResponse{Headers: headers, Body: "error unmarshalling request", StatusCode: 400}, nil
	}

	base64dec, err := base64.StdEncoding.DecodeString(reqJSON.ImageBase64)
	if err != nil {
		return events.APIGatewayProxyResponse{Headers: headers, Body: "error decoding image base 64\n", StatusCode: 500}, nil
	}

	_, err = _uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(base64dec),
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		fmt.Println("error uploading to s3")
		return events.APIGatewayProxyResponse{Headers: headers, Body: "unable to upload to s3\n", StatusCode: 500}, nil
	}

	textIn := &rekognition.DetectTextInput{
		Image: &rekognition.Image{
			S3Object: &rekognition.S3Object{
				Bucket: aws.String(bucket),
				Name:   aws.String(key),
			},
		},
	}

	textRes, err := _rekognition.DetectText(textIn)
	if err != nil {
		return events.APIGatewayProxyResponse{Headers: headers, Body: "unable to detect text\n", StatusCode: 500}, nil
	}

	var textOut bytes.Buffer
	for _, text := range textRes.TextDetections {
		textOut.WriteString(*text.DetectedText)
		textOut.WriteString("\n")
	}

	//resp := events.APIGatewayProxyResponse{Headers: make(map[string]string), Body: textOut.String(), StatusCode: 200}
	//resp.Headers["Access-Control-Allow-Origin"] = "*"
	return events.APIGatewayProxyResponse{Headers: headers, Body: textOut.String(), StatusCode: 200}, nil
	//return resp, nil
}

func main() {
	lambda.Start(handler)
}
