package helpers

import (
	"reflect"
	"shopit/models"
	"strings"
)

func ValidateString(value string) bool {
	var checkType = reflect.TypeOf(value).Kind()
	if checkType == reflect.String {
		var trimString = strings.TrimSpace(value)
		return len(trimString) != 0
	}
	return false
}

func ValidateItem(body models.Item) string {
	var checkCategory = false
	if body.Category >= 0 {
		value := body.Category
		checkCategory = ValidateString(Category[value])
	}

	var checkName = body.Name != "" && ValidateString(body.Name)
	var checkStock = body.Stock >= 0
	var checkPrice = body.Price > 0
	var checkImage = body.Image != ""
	var _, imageExt = GetImageNameType(body.Image)

	if _, err = IsAllowedImageExt(imageExt); err != nil {
		return err.Error()
	}
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
