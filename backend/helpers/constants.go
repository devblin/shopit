package helpers

type Collections string

var Category = map[int]string{
	0:  "Appliances",
	1:  "Arts, Crafts & Sewing ",
	2:  "Automotive",
	3:  "Baby",
	4:  "Beauty",
	5:  "Books",
	6:  "Collectibles & Fine Arts",
	7:  "Electronics",
	8:  "Clothing, Shoes & Jewelry",
	9:  "Gift Cards",
	10: "Grocery & Gourmet Food",
	11: "Handmade",
	12: "Health & Personal Care",
	13: "Home & Kitchen",
	14: "Industrial & Scientific",
	15: "Patio, Lawn & Garden",
	16: "Luggage & Travel Gear",
	17: "Apps & Games",
	18: "Movies & TV",
	19: "Office Products",
	20: "Computers",
	21: "Software",
	22: "Sports & Outdoors",
	23: "Others",
}

var AllowedImageTypes = map[string]bool{
	"jpeg": true,
	"jpg":  true,
	"png":  true,
}
