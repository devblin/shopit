package database

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	HELPERS "shopit/helpers"
	MODELS "shopit/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ITEM_COLLECTION_NAME = "items"
var ITEM_MODEL primitive.M
var ITEM_COLLECTION *mongo.Collection
var MESSAGE string
var DESCRIPTION string
var STATUS int = http.StatusOK
var DATA interface{}
var AWS_REGION string
var AWS_ACCESS_KEY string
var AWS_SECRET_KEY string
var AWS_S3_BUCKET string
var THUMBNAIL_HEIGHT = 64
var STANDARD_HEIGHT = 400
var AWS_CONFIG *aws.Config
var DEFAULT_ITEM_IMAGE_NAME string
var IMAGE_DIMENSIONS = []uint{uint(THUMBNAIL_HEIGHT), uint(STANDARD_HEIGHT)}

func init() {
	ITEM_MODEL = MODELS.ItemSchema
	ITEM_COLLECTION = DB.Collection(ITEM_COLLECTION_NAME)
	DEFAULT_ITEM_IMAGE_NAME = HELPERS.GetEnv("DEFAULT_ITEM_IMAGE_NAME")
	AWS_REGION = HELPERS.GetEnv("AWS_REGION")
	AWS_ACCESS_KEY = HELPERS.GetEnv("AWS_ACCESS_KEY")
	AWS_SECRET_KEY = HELPERS.GetEnv("AWS_SECRET_KEY")
	AWS_S3_BUCKET = HELPERS.GetEnv("AWS_S3_BUCKET")
	AWS_CONFIG = &aws.Config{
		Region:      aws.String(AWS_REGION),
		Credentials: credentials.NewStaticCredentials(AWS_ACCESS_KEY, AWS_SECRET_KEY, ""),
	}
}

func GetItemList(c *gin.Context) {
	DESCRIPTION = "Items list"
	STATUS = http.StatusInternalServerError
	DATA = []bson.M{}
	var cursor *mongo.Cursor

	if cursor, err = ITEM_COLLECTION.Find(CTX, bson.M{}); err != nil {
		MESSAGE = err.Error()
		STATUS = http.StatusInternalServerError
	} else {
		if err = cursor.All(CTX, &DATA); err != nil {
			MESSAGE = err.Error()
			STATUS = http.StatusInternalServerError
		} else {
			MESSAGE = "Success"
			STATUS = http.StatusOK
		}
	}

	c.JSON(STATUS, gin.H{
		"description": DESCRIPTION,
		"data":        DATA,
		"message":     MESSAGE,
	})
}

func AddItem(c *gin.Context) {
	DESCRIPTION = "Add item"
	STATUS = http.StatusInternalServerError
	DATA = nil
	var jsonBody bson.M = nil
	var rawBody, _ = ioutil.ReadAll(c.Request.Body)
	var validate string = ""
	var itemModel MODELS.Item = MODELS.Item{}
	var data *mongo.InsertOneResult

	if err = json.Unmarshal(rawBody, &jsonBody); err != nil {
		STATUS = http.StatusInternalServerError
		MESSAGE = err.Error()
	} else {
		validate = HELPERS.ValidateItem(jsonBody)

		if validate == "" {
			itemModel.Name = jsonBody["Name"].(string)
			itemModel.Category = int(jsonBody["Category"].(float64))
			itemModel.Price = int(jsonBody["Price"].(float64))
			itemModel.Stock = int(jsonBody["Stock"].(float64))
			itemModel.Sold = int(jsonBody["Sold"].(float64))
			if jsonBody["Image"].(string) != DEFAULT_ITEM_IMAGE_NAME {
				_, ext := HELPERS.GetImageNameType(jsonBody["Image"].(string))
				itemModel.Image = uuid.NewString() + "." + ext
			} else {
				itemModel.Image = jsonBody["Image"].(string)
			}
		}

		if validate != "" {
			STATUS = http.StatusBadRequest
			MESSAGE = validate
		} else if data, err = ITEM_COLLECTION.InsertOne(CTX, &itemModel); err != nil {
			STATUS = http.StatusInternalServerError
			MESSAGE = err.Error()
		} else {
			if itemModel.Image != DEFAULT_ITEM_IMAGE_NAME {
				jsonBody = bson.M{}
				jsonBody["InsertedID"] = data.InsertedID
				jsonBody["ImageID"] = itemModel.Image
				DATA = jsonBody
			} else {
				DATA = data
			}
			STATUS = http.StatusOK
			MESSAGE = "Item added successfully"
		}
	}

	c.JSON(STATUS, gin.H{
		"description": DESCRIPTION,
		"message":     MESSAGE,
		"data":        DATA,
	})
}

func GetItemDetails(c *gin.Context) {
	DESCRIPTION = "Item details with given id"
	STATUS = http.StatusInternalServerError
	DATA = nil
	var data bson.M = nil
	var itemId, _ = c.Params.Get("itemId")

	if !HELPERS.ValidateString(itemId) {
		STATUS = http.StatusBadRequest
		MESSAGE = "Invalid item id"
	} else {
		if bsonItemId, err := primitive.ObjectIDFromHex(itemId); err != nil {
			STATUS = http.StatusInternalServerError
			MESSAGE = err.Error()
		} else if err = ITEM_COLLECTION.FindOne(CTX, bson.M{"_id": bsonItemId}).Decode(&data); err != nil {
			STATUS = http.StatusInternalServerError
			MESSAGE = err.Error()
		} else {
			DATA = data
			STATUS = http.StatusOK
			MESSAGE = "Success"
		}
	}

	c.JSON(STATUS, gin.H{
		"description": DESCRIPTION,
		"data":        DATA,
		"message":     MESSAGE,
	})
}

func DeleteItem(c *gin.Context) {
	DESCRIPTION = "Delete item with given id"
	DATA = nil
	STATUS = http.StatusInternalServerError
	err = nil
	var itemId primitive.ObjectID
	var data *mongo.DeleteResult = nil
	var itemData bson.M = nil
	var body = bson.M{}
	var imageFullName string

	if err = c.BindJSON(&body); err != nil {
		STATUS = http.StatusInternalServerError
		MESSAGE = err.Error()
	} else {
		if itemId, err = primitive.ObjectIDFromHex(body["id"].(string)); err != nil {
			STATUS = http.StatusInternalServerError
			MESSAGE = err.Error()
		} else {
			if err = ITEM_COLLECTION.FindOne(CTX, bson.M{"_id": itemId}).Decode(&itemData); err != nil {
				STATUS = http.StatusInternalServerError
				MESSAGE = err.Error()
			} else {
				imageFullName = itemData["image"].(string)

				if data, err = ITEM_COLLECTION.DeleteOne(CTX, bson.M{"_id": itemId}); err != nil {
					STATUS = http.StatusInternalServerError
					MESSAGE = err.Error()
				} else {
					DATA = data
					if data != nil && data.DeletedCount == 0 {
						STATUS = http.StatusAlreadyReported
						MESSAGE = "Item doesn't exists"
					} else {
						STATUS = http.StatusOK
						MESSAGE = "Item successfully deleted"

						if imageFullName != DEFAULT_ITEM_IMAGE_NAME {
							if err = deleteObjectInS3(imageFullName); err != nil {
								STATUS = http.StatusInternalServerError
								MESSAGE = err.Error()
							}
						}
					}
				}
			}
		}
	}

	c.JSON(STATUS, gin.H{
		"description": DESCRIPTION,
		"data":        DATA,
		"message":     MESSAGE,
	})
}

func UpdateItem(c *gin.Context) {
	DESCRIPTION = "Update item with given id"
	DATA = nil
	var itemData bson.M
	var update bson.D
	var data *mongo.UpdateResult
	var jsonBody bson.M = nil
	var rawBody, _ = ioutil.ReadAll(c.Request.Body)
	var validate string = ""
	var itemId, imageName, currImageName string

	if err = json.Unmarshal(rawBody, &jsonBody); err != nil {
		STATUS = http.StatusFailedDependency
		MESSAGE = err.Error()
	} else {
		validate = HELPERS.ValidateItem(jsonBody)

		if validate != "" {
			STATUS = http.StatusBadRequest
			MESSAGE = validate
		} else {
			itemId = jsonBody["id"].(string)
			imageName = jsonBody["Image"].(string)

			if bsonItemId, err := primitive.ObjectIDFromHex(itemId); err != nil {
				STATUS = http.StatusInternalServerError
				MESSAGE = err.Error()
			} else if err := ITEM_COLLECTION.FindOne(CTX, bson.M{"_id": bsonItemId}).Decode(&itemData); err != nil {
				STATUS = http.StatusInternalServerError
				MESSAGE = err.Error()
			} else {
				currImageName = itemData["image"].(string)

				if imageName != DEFAULT_ITEM_IMAGE_NAME {
					_, ext := HELPERS.GetImageNameType(imageName)

					if currImageName != DEFAULT_ITEM_IMAGE_NAME && imageName != currImageName {
						// update currImageName ext with imageName ext
						name, _ := HELPERS.GetImageNameType(currImageName)
						if err = deleteObjectInS3(currImageName); err != nil {
							STATUS = http.StatusInternalServerError
							MESSAGE = err.Error()
						}
						currImageName = name + "." + ext
					} else if currImageName == DEFAULT_ITEM_IMAGE_NAME {
						// generate new currImageName with imageName ext
						name := uuid.NewString()
						currImageName = name + "." + ext
					}
					// Save currImageName and send the currImageName
					jsonBody["Image"] = currImageName
				} else if currImageName != DEFAULT_ITEM_IMAGE_NAME {
					// delete currImageName from S3
					if err = deleteObjectInS3(currImageName); err != nil {
						STATUS = http.StatusInternalServerError
						MESSAGE = err.Error()
					}
				}

				update = bson.D{{"$set", bson.D{
					{"name", jsonBody["Name"]},
					{"price", jsonBody["Price"]},
					{"stock", jsonBody["Stock"]},
					{"image", jsonBody["Image"]},
					{"category", jsonBody["Category"]},
				}}}

				if data, err = ITEM_COLLECTION.UpdateOne(CTX, bson.M{"_id": bsonItemId}, update); err != nil {
					STATUS = http.StatusInternalServerError
					MESSAGE = err.Error()
				} else {
					jsonBody = bson.M{}
					jsonBody["ImageID"] = currImageName
					jsonBody["MatchedCount"] = data.MatchedCount
					jsonBody["ModifiedCount"] = data.ModifiedCount
					STATUS = http.StatusOK
					MESSAGE = "Item successfully updated"
				}
			}
		}
	}

	c.JSON(STATUS, gin.H{
		"description": DESCRIPTION,
		"data":        jsonBody,
		"message":     MESSAGE,
	})
}

func UploadImage(c *gin.Context) {
	DESCRIPTION = "Upload item image"
	var awsSession *session.Session

	if err = c.Request.ParseMultipartForm(32 << 20); err != nil {
		STATUS = http.StatusRequestEntityTooLarge
		MESSAGE = err.Error()
	} else if image, imageHeader, err := c.Request.FormFile("image"); err != nil {
		STATUS = http.StatusFailedDependency
		MESSAGE = err.Error()
	} else {
		if awsSession, err = session.NewSession(AWS_CONFIG); err != nil {
			STATUS = http.StatusFailedDependency
			MESSAGE = err.Error()
		} else {
			count := len(IMAGE_DIMENSIONS)

			for i := 0; i < count; i++ {
				if err = addItemImagesInS3(awsSession, image, imageHeader.Filename, IMAGE_DIMENSIONS[i]); err != nil {
					STATUS = http.StatusFailedDependency
					MESSAGE = err.Error()
					break
				} else {
					STATUS = http.StatusOK
					MESSAGE = "Success"
				}
				image.Seek(0, io.SeekStart)
			}
		}
	}

	c.JSON(STATUS, gin.H{
		"description": DESCRIPTION,
		"message":     MESSAGE,
	})
}

func putObjectInS3(s3Client *s3.S3, fileName string, file io.ReadSeeker) error {
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(AWS_S3_BUCKET),
		Key:    aws.String(fileName),
		Body:   file,
		ACL:    aws.String("private"),
	})

	return err
}

func addItemImagesInS3(s *session.Session, file multipart.File, fileName string, size uint) error {
	var s3Client = s3.New(s)
	var imageName, imageExt = HELPERS.GetImageNameType(fileName)
	var buffer = HELPERS.GetResizedImages(size, file, fileName)

	tempName := HELPERS.GetImageName(size, imageName, imageExt)
	err = putObjectInS3(s3Client, tempName, bytes.NewReader(buffer.Bytes()))

	return err
}

func deleteObjectInS3(imageFullName string) error {
	var imageName, imageExt = HELPERS.GetImageNameType(imageFullName)
	var awsSession *session.Session
	if awsSession, err = session.NewSession(AWS_CONFIG); err != nil {
		return err
	}
	var count = len(IMAGE_DIMENSIONS)
	var s3Client = s3.New(awsSession)

	for i := 0; i < count; i++ {
		_, err = s3Client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(AWS_S3_BUCKET),
			Key:    aws.String(HELPERS.GetImageName(IMAGE_DIMENSIONS[i], imageName, imageExt)),
		})
	}

	return err
}
