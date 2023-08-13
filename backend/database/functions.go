package database

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"reflect"
	"shopit/helpers"
	"shopit/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/mongo"
)

var ITEM_COLLECTION *mongo.Collection
var MESSAGE string
var DESCRIPTION string
var STATUS int = http.StatusOK
var DATA interface{}
var AWS_S3_BUCKET string
var AWS_DYNAMO_DB_ITEM string
var THUMBNAIL_HEIGHT = 64
var STANDARD_HEIGHT = 400
var DEFAULT_ITEM_IMAGE_NAME string
var IMAGE_DIMENSIONS = []uint{uint(THUMBNAIL_HEIGHT), uint(STANDARD_HEIGHT)}

func init() {
	DEFAULT_ITEM_IMAGE_NAME = helpers.GetEnv("DEFAULT_ITEM_IMAGE_NAME")
	AWS_S3_BUCKET = helpers.GetEnv("AWS_S3_BUCKET")
	AWS_DYNAMO_DB_ITEM = helpers.GetEnv("AWS_DYNAMO_DB_ITEM")
}

func GetItemList(c *gin.Context) {
	DESCRIPTION = "Items list"
	STATUS = http.StatusOK
	DATA = []models.Item{}
	MESSAGE = ""
	var err error = nil
	var scanOutput *dynamodb.ScanOutput

	if scanOutput, err = AWS_DYNAMO_DB.Scan(&dynamodb.ScanInput{
		TableName: aws.String(AWS_DYNAMO_DB_ITEM),
	}); err != nil {
		MESSAGE = err.Error()
		STATUS = http.StatusInternalServerError
	} else if err = dynamodbattribute.UnmarshalListOfMaps(scanOutput.Items, &DATA); err != nil {
		MESSAGE = err.Error()
		STATUS = http.StatusInternalServerError
	}

	c.JSON(STATUS, gin.H{
		"description": DESCRIPTION,
		"data":        DATA,
		"message":     MESSAGE,
	})
}

func AddItem(c *gin.Context) {
	DESCRIPTION = "Add item"
	STATUS = http.StatusOK
	MESSAGE = "Item added successfully"
	var itemOut = models.Item{}
	var itemInput = models.Item{}
	var rawBody, _ = ioutil.ReadAll(c.Request.Body)

	if err := json.Unmarshal(rawBody, &itemInput); err != nil {
		STATUS = http.StatusInternalServerError
		MESSAGE = err.Error()
	} else {
		var validate = helpers.ValidateItem(itemInput)

		if validate != "" {
			STATUS = http.StatusBadRequest
			MESSAGE = validate
		} else {
			itemInput.Id = uuid.NewString()
			if itemInput.Image != DEFAULT_ITEM_IMAGE_NAME {
				_, ext := helpers.GetImageNameType(itemInput.Image)
				itemInput.Image = uuid.NewString() + "." + ext
			}

			if itemAttr, err := dynamodbattribute.MarshalMap(itemInput); err != nil {
				STATUS = http.StatusInternalServerError
				MESSAGE = err.Error()
			} else if putItemOut, err := AWS_DYNAMO_DB.PutItem(&dynamodb.PutItemInput{
				Item:      itemAttr,
				TableName: aws.String(AWS_DYNAMO_DB_ITEM),
			}); err != nil {
				STATUS = http.StatusInternalServerError
				MESSAGE = err.Error()
			} else if err := dynamodbattribute.UnmarshalMap(putItemOut.Attributes, &itemOut); err != nil {
				STATUS = http.StatusInternalServerError
				MESSAGE = err.Error()
			}
		}
	}

	c.JSON(STATUS, gin.H{
		"description": DESCRIPTION,
		"message":     MESSAGE,
		"data":        itemInput,
	})
}

func GetItemDetails(c *gin.Context) {
	DESCRIPTION = "Item details with given id"
	STATUS = http.StatusOK
	DATA = nil
	var Id, _ = c.Params.Get("itemId")
	var itemOut models.Item

	if !helpers.ValidateString(Id) {
		STATUS = http.StatusBadRequest
		MESSAGE = "Invalid item id"
	} else if getItemOut, err := AWS_DYNAMO_DB.GetItem(&dynamodb.GetItemInput{
		Key:       map[string]*dynamodb.AttributeValue{"Id": {S: aws.String(Id)}},
		TableName: aws.String(AWS_DYNAMO_DB_ITEM),
	}); err != nil {
		STATUS = http.StatusInternalServerError
		MESSAGE = err.Error()
	} else if err := dynamodbattribute.UnmarshalMap(getItemOut.Item, &itemOut); err != nil {
		STATUS = http.StatusInternalServerError
		MESSAGE = err.Error()
	} else {
		DATA = itemOut
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
	MESSAGE = "Item successfully deleted"
	STATUS = http.StatusOK
	var itemOut = models.Item{}
	var itemInput = models.Item{}
	var imageFullName string

	if err := c.BindJSON(&itemInput); err != nil {
		STATUS = http.StatusInternalServerError
		MESSAGE = err.Error()
	} else {
		if getItemOut, err := AWS_DYNAMO_DB.GetItem(&dynamodb.GetItemInput{
			TableName: aws.String(AWS_DYNAMO_DB_ITEM),
			Key:       map[string]*dynamodb.AttributeValue{"Id": {S: aws.String(itemInput.Id)}},
		}); err != nil {
			STATUS = http.StatusInternalServerError
			MESSAGE = err.Error()
		} else if err = dynamodbattribute.UnmarshalMap(getItemOut.Item, &itemOut); err != nil {
			STATUS = http.StatusInternalServerError
			MESSAGE = err.Error()
		} else if _, err := AWS_DYNAMO_DB.DeleteItem(&dynamodb.DeleteItemInput{
			TableName: aws.String(AWS_DYNAMO_DB_ITEM),
			Key:       map[string]*dynamodb.AttributeValue{"Id": {S: aws.String(itemInput.Id)}},
		}); err != nil {
			STATUS = http.StatusInternalServerError
			MESSAGE = err.Error()
		} else {
			imageFullName = itemOut.Image

			if imageFullName != DEFAULT_ITEM_IMAGE_NAME {
				if err = deleteObjectInS3(imageFullName); err != nil {
					STATUS = http.StatusInternalServerError
					MESSAGE = err.Error()
				}
			}

			DATA = itemOut
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
	STATUS = http.StatusOK
	MESSAGE = "Item successfully updated"
	var itemInput = models.Item{}
	var rawBody, _ = ioutil.ReadAll(c.Request.Body)
	var validate string = ""
	var itemId, imageName, currImageName string
	var err error = nil
	var itemOut models.Item

	if err = json.Unmarshal(rawBody, &itemInput); err != nil {
		STATUS = http.StatusFailedDependency
		MESSAGE = err.Error()
	} else {
		validate = helpers.ValidateItem(itemInput)

		if validate != "" {
			STATUS = http.StatusBadRequest
			MESSAGE = validate
		} else {
			itemId = itemInput.Id
			imageName = itemInput.Image
			_, ext := helpers.GetImageNameType(imageName)

			if getItemOut, err := AWS_DYNAMO_DB.GetItem(&dynamodb.GetItemInput{
				Key:       map[string]*dynamodb.AttributeValue{"Id": {S: aws.String(itemId)}},
				TableName: aws.String(AWS_DYNAMO_DB_ITEM),
			}); err != nil {
				STATUS = http.StatusNotFound
				MESSAGE = err.Error()
			} else if err := dynamodbattribute.UnmarshalMap(getItemOut.Item, &itemOut); err != nil {
				STATUS = http.StatusInternalServerError
				MESSAGE = err.Error()
			} else {
				currImageName = itemOut.Image

				if imageName != DEFAULT_ITEM_IMAGE_NAME {
					if currImageName != DEFAULT_ITEM_IMAGE_NAME && imageName != currImageName {
						// update currImageName ext with imageName ext
						name, _ := helpers.GetImageNameType(currImageName)
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
					// save currImageName and send the currImageName
					itemInput.Image = currImageName
				} else if currImageName != DEFAULT_ITEM_IMAGE_NAME {
					// delete currImageName from S3
					if err = deleteObjectInS3(currImageName); err != nil {
						STATUS = http.StatusInternalServerError
						MESSAGE = err.Error()
					}
				}

				if STATUS == http.StatusOK {
					itemInputType := reflect.TypeOf(itemInput)
					itemInputValue := reflect.ValueOf(itemInput)
					attributeUpdates := make(map[string]*dynamodb.AttributeValueUpdate)
					for i := 0; i < itemInputType.NumField(); i++ {
						key := itemInputType.Field(i).Name
						value := itemInputValue.Field(i).Interface()

						if key == "Id" {
							continue
						}

						if attributeValue, err := dynamodbattribute.Marshal(value); err != nil {
							STATUS = http.StatusInternalServerError
							MESSAGE = err.Error()
						} else {
							attributeUpdates[key] = &dynamodb.AttributeValueUpdate{
								Value:  attributeValue,
								Action: aws.String(dynamodb.AttributeActionPut),
							}
						}
					}

					if _, err := AWS_DYNAMO_DB.UpdateItem(&dynamodb.UpdateItemInput{
						TableName:        aws.String(AWS_DYNAMO_DB_ITEM),
						Key:              map[string]*dynamodb.AttributeValue{"Id": {S: aws.String(itemId)}},
						ReturnValues:     aws.String(dynamodb.ReturnValueAllNew),
						AttributeUpdates: attributeUpdates,
					}); err != nil {
						STATUS = http.StatusInternalServerError
						MESSAGE = err.Error()
					}
				}
			}
		}
	}

	c.JSON(STATUS, gin.H{
		"description": DESCRIPTION,
		"data":        itemInput,
		"message":     MESSAGE,
	})
}

func UploadImage(c *gin.Context) {
	var err error = nil
	DESCRIPTION = "Upload item image"
	STATUS = http.StatusOK
	MESSAGE = "Success"

	if err = c.Request.ParseMultipartForm(32 << 20); err != nil {
		STATUS = http.StatusRequestEntityTooLarge
		MESSAGE = err.Error()
	} else if image, imageHeader, err := c.Request.FormFile("image"); err != nil {
		STATUS = http.StatusFailedDependency
		MESSAGE = err.Error()
	} else {
		count := len(IMAGE_DIMENSIONS)

		for i := 0; i < count; i++ {
			if err = addItemImagesInS3(image, imageHeader.Filename, IMAGE_DIMENSIONS[i]); err != nil {
				STATUS = http.StatusFailedDependency
				MESSAGE = err.Error()
				break
			}
			image.Seek(0, io.SeekStart)
		}
	}

	c.JSON(STATUS, gin.H{
		"description": DESCRIPTION,
		"message":     MESSAGE,
	})
}

func putObjectInS3(fileName string, file io.ReadSeeker) error {
	var err error = nil
	_, err = AWS_S3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(AWS_S3_BUCKET),
		Key:    aws.String(fileName),
		Body:   file,
	})

	return err
}

func addItemImagesInS3(file multipart.File, fileName string, size uint) error {
	var imageName, imageExt = helpers.GetImageNameType(fileName)
	if buffer, err := helpers.GetResizedImages(size, file, fileName); err != nil {
		return err
	} else {
		tempName := helpers.GetImageName(size, imageName, imageExt)

		err := putObjectInS3(tempName, bytes.NewReader(buffer.Bytes()))
		return err
	}
}

func deleteObjectInS3(imageFullName string) error {
	var err error = nil
	var imageName, imageExt = helpers.GetImageNameType(imageFullName)
	var count = len(IMAGE_DIMENSIONS)

	for i := 0; i < count; i++ {
		_, err = AWS_S3.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(AWS_S3_BUCKET),
			Key:    aws.String(helpers.GetImageName(IMAGE_DIMENSIONS[i], imageName, imageExt)),
		})
	}

	return err
}
