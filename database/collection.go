package database

import (
	"log"

	MODELS "shopit/models"

	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

var COLLECTIONS []string

func checkCollection(name string) bool {
	COLLECTIONS, err = DB.ListCollectionNames(CTX, bson.M{})

	if len(COLLECTIONS) > 0 {
		for _, collection := range COLLECTIONS {
			if collection == name {
				return true
			}
		}
	}

	return false
}

func CreateCollection(name string) {
	if check := checkCollection(name); !check {
		opts := options.CreateCollection().SetValidator(MODELS.Validator)
		if err := DB.CreateCollection(CTX, name, opts); err != nil {
			log.Println(err.Error())
		}
	} else {
		log.Printf("%s: COLLECTTION ALREADY EXISTS", name)
	}
}
