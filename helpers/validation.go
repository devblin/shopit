package helpers

import (
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func ValidateString(value string) bool {
	var checkType = reflect.TypeOf(value).Kind()
	if checkType == reflect.String {
		var trimString = strings.TrimSpace(value)
		return len(trimString) != 0
	}
	return false
}

func ValidateItem(body bson.M) string {
	var checkCategory = false
	if body["Category"] != nil {
		value := int(body["Category"].(float64))
		checkCategory = ValidateString(Category[value])
	}

	var checkName = body["Name"] != nil && ValidateString(body["Name"].(string))
	var checkStock = body["Stock"] != nil && int(body["Stock"].(float64)) >= 0
	var checkPrice = body["Price"] != nil && body["Price"].(float64) > 0
	var checkImage = body["Image"] != nil && ValidateString(body["Image"].(string))

	if !checkCategory {
		return "Invalid category"
	}
	if !checkName {
		return "Invalid name"
	}
	if !checkStock {
		return "Invalid stock"
	}
	if !checkPrice {
		return "Invalid price"
	}
	if !checkImage {
		return "Invalid image"
	}

	return ""
}
