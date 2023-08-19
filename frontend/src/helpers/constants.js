const BASE_API_URL = process.env.REACT_APP_ENV === "dev" ? process.env.REACT_APP_API : ''
const BASE_S3_URL = process.env.REACT_APP_S3_URL.replace(/\/$/, "");

const
	PRIMARY_COLOR = '#6c5af2',
	IMAGE_BASE_URL = BASE_S3_URL,
	DEFAULT_IMAGE_NAME = process.env.REACT_APP_DEFAULT_ITEM_IMAGE_NAME,
	DEFAULT_IMAGE_URL = {
		SIZE400: `${BASE_S3_URL}/default-product400.jpg`,
		SIZE64: `${BASE_S3_URL}/default-product64.jpg`
	},
	PAGES = {
		HOME: "/",
		ITEMS: "/",
		INVENTORY: "/inventory",
	},
	API = {
		ITEMS: `${BASE_API_URL}/api/item/list`,
		ITEM_DETAILS: (id) => `${BASE_API_URL}/api/item/` + id,
		ITEM: `${BASE_API_URL}/api/item/`,
		ITEM_IMAGE: `${BASE_API_URL}/api/item/image`,
	},
	DEFAULT_ALERT_DATA = {
		open: false,
		message: "",
		severity: "success"
	},
	SEVERITY = {
		INFO: "info",
		ERROR: "error",
		SUCCESS: "success",
		WARN: "warning",
	},
	MESSAGES = {
		ERROR_FETCHING: "Some error occured while fetching resources",
		SERVER_ERROR: "Server error, please try again later",
		NO_ITEMS: "No items",
		INVENTORY_EMPTY: "Inventory empty, to add item(s) please click on \"+ Add Item\" button."
	}

const ITEMS_CATEGORY = {
	0: "Appliances",
	1: "Arts, Crafts & Sewing ",
	2: "Automotive",
	3: "Baby",
	4: "Beauty",
	5: "Books",
	6: "Collectibles & Fine Arts",
	7: "Electronics",
	8: "Clothing, Shoes & Jewelry",
	9: "Gift Cards",
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

export {
	PRIMARY_COLOR,
	PAGES,
	DEFAULT_ALERT_DATA,
	SEVERITY,
	API,
	MESSAGES,
	DEFAULT_IMAGE_URL,
	ITEMS_CATEGORY,
	IMAGE_BASE_URL,
	DEFAULT_IMAGE_NAME
};