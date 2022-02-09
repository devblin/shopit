package helpers

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
)

var err error

func GetImageName(unique uint, name string, extension string) string {
	return name + strconv.Itoa(int(unique)) + "." + extension
}

func GetImageNameType(name string) (string, string) {
	var lastIndex = strings.LastIndex(name, ".")
	var imageName = name[:lastIndex]
	var imageType = name[lastIndex+1:]

	return imageName, imageType
}

func GetResizedImages(size uint, file multipart.File, fileName string) bytes.Buffer {
	var _, imageType = GetImageNameType(fileName)
	var img, tempImage image.Image
	var imgBuf bytes.Buffer

	if imageType == "jpeg" || imageType == "jpg" {
		if img, err = jpeg.Decode(file); err != nil {
			log.Fatal(err)
		}
	} else if imageType == "png" {
		if img, err = png.Decode(file); err != nil {
			log.Fatal(err)
		}
	}

	tempImage = resize.Resize(0, size, img, resize.Lanczos3)

	if imageType == "jpeg" || imageType == "jpg" {
		if err = jpeg.Encode(&imgBuf, tempImage, nil); err != nil {
			log.Fatal(err)
		}
	} else if imageType == "png" {
		if err = png.Encode(&imgBuf, tempImage); err != nil {
			log.Fatal(err)
		}
	}

	return imgBuf
}
