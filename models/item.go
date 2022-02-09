package models

import (
	"go.mongodb.org/mongo-driver/bson"
)

type Item struct {
	Name     string `json:"Name"`
	Price    int    `json:"Price"`
	Stock    int    `json:"Stock"`
	Category int    `json:"Category"`
	Image    string `json:"Image"`
	Sold     int    `json:"Sold"`
}

var ItemSchema = bson.M{
	"bsonType": "object",
	"properties": bson.M{
		"Name": bson.M{
			"bsonType":    "string",
			"description": "Item name",
		},
		"Price": bson.M{
			"bsonType":    "int",
			"minimum":     0,
			"description": "Price of item",
		},
		"Stock": bson.M{
			"bsonType":    "int",
			"minimum":     0,
			"description": "Available number of item",
		},
		"Category": bson.M{
			"bsonType":    "int",
			"minimum":     0,
			"description": "Category of item",
		},
		"Image": bson.M{
			"bsonType":    "string",
			"description": "Image of item",
		},
		"Sold": bson.M{
			"bsonType":    "int",
			"minimum":     0,
			"description": "Number of items sold",
		},
	},
}

var Validator = bson.M{
	"$jsonSchema": ItemSchema,
}
