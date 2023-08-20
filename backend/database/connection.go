package database

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"

	"go.mongodb.org/mongo-driver/mongo"

	"shopit/helpers"

	"github.com/aws/aws-sdk-go/aws/session"
)

var CTX context.Context = nil
var CANCEL context.CancelFunc
var NUMBER_DB_SESSIONS int = 0
var DB *mongo.Database = nil
var AWS_SESSION *session.Session
var AWS_DYNAMO_DB *dynamodb.DynamoDB
var AWS_S3 *s3.S3

func init() {
	var ENDPOINT string = ""
	if helpers.GetEnv("ENV") == helpers.DEV {
		ENDPOINT = helpers.GetEnv("LOCALSTACK_BASE_URL")
	}

	var err error = nil
	var AWS_REGION = helpers.GetEnv("AWS_REGION")
	var AWS_ACCESS_KEY_ID = helpers.GetEnv("AWS_ACCESS_KEY_ID")
	var AWS_SECRET_ACCESS_KEY = helpers.GetEnv("AWS_SECRET_ACCESS_KEY")
	var AWS_CONFIG = &aws.Config{
		Region:           aws.String(AWS_REGION),
		Credentials:      credentials.NewStaticCredentials(AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, ""),
		Endpoint:         aws.String(ENDPOINT),
		S3ForcePathStyle: aws.Bool(true),
	}

	AWS_SESSION, err = session.NewSession(AWS_CONFIG)
	if err != nil {
		log.Fatal(err)
	}
	AWS_DYNAMO_DB = dynamodb.New(AWS_SESSION)
	AWS_S3 = s3.New(AWS_SESSION)
	log.Print("CONNECTED TO DATABASE...")
}
