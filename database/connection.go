package database

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/mongo"

	HELPERS "shopit/helpers"

	"go.mongodb.org/mongo-driver/mongo/options"
)

var CTX context.Context = nil
var CANCEL context.CancelFunc
var err error = nil
var CLIENT *mongo.Client
var CLIENT_IP string = ""
var NUMBER_DB_SESSIONS int = 0
var DB *mongo.Database = nil

func init() {
	var DATABASE_URL = HELPERS.GetEnv("DATABASE_URL")
	var DATABASE_NAME = HELPERS.GetEnv("DATABASE_NAME")
	var CLIENT_OPTIONS = options.Client().ApplyURI(DATABASE_URL)

	CTX = context.Background()
	CLIENT, err = mongo.Connect(CTX, CLIENT_OPTIONS)
	if err != nil {
		log.Fatal(err)
	}

	err = CLIENT.Ping(CTX, nil)
	if err != nil {
		log.Fatal(err)
	}

	DB = CLIENT.Database(DATABASE_NAME)
	CreateCollection(ITEM_COLLECTION_NAME)
	log.Print("CONNECTED TO DATABASE")
}

func CheckConnection(c *gin.Context) {
	if err = CLIENT.Ping(CTX, nil); err != nil {
		MESSAGE = err.Error()
		STATUS = http.StatusBadGateway
	} else {
		CLIENT_IP = c.ClientIP()
		NUMBER_DB_SESSIONS = CLIENT.NumberSessionsInProgress()
		MESSAGE = "Ok"
		STATUS = http.StatusOK
	}

	c.JSON(STATUS, gin.H{
		"message":               MESSAGE,
		"client":                CLIENT_IP,
		"number_of_connections": NUMBER_DB_SESSIONS,
	})
}
