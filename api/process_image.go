package api

import (
	"bytes"
	"github.com/disintegration/imaging"
	"github.com/labstack/echo/v4"
	"image"
	"image/jpeg"
	"image/png"
	"image_processing/errors"
	"image_processing/types"
	"os"
)

const storagePath = "./storage/"

type ProcessedImage struct {
	echo.Context
	*ImageHandler
	id int
}

func NewProcessedImage(c echo.Context, h *ImageHandler, id int) *ProcessedImage {
	return &ProcessedImage{c, h, id}
}
func (i *ProcessedImage) getImageFile() (image.Image, *errors.Error) {

	dbRes, err := i.getImageFromDB()
	if err != nil {
		return nil, err
	}
	imgFile, err := searchFileAndConvert(storagePath, dbRes.Name.String()+"."+dbRes.Format)

	if err != nil {
		return nil, err
	}
	return imgFile, nil
}

func searchFileAndConvert(dirPath, fileName string) (image.Image, *errors.Error) {

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, errors.ErrServer(err.Error())
	}

	for _, file := range files {
		if file.Name() == fileName {
			src, err := imaging.Open(dirPath + fileName)
			if err != nil {
				return nil, errors.ErrServer("failed to open image: " + err.Error())

			}
			return src, nil
		}
	}
	return nil, errors.ErrServer("error file not found")
}
func (i *ProcessedImage) getImageFromDB() (*types.Image, *errors.Error) {
	image, err := i.imageStore.GetImageByID(i.Request().Context(), i.id)
	if err != nil {
		return nil, err
	}
	i.Set("format", image.Format)
	return image, nil
}

func (i *ProcessedImage) sendImage(img image.Image) *errors.Error {
	buf := new(bytes.Buffer)
	format := i.Get("format")
	if format == "jpeg" {
		err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 100})
		if err != nil {
			return errors.ErrServer("error while encoding image to jpeg" + err.Error())
		}
	} else {
		err := png.Encode(buf, img)
		if err != nil {
			return errors.ErrServer("error while encoding image to png " + err.Error())
		}
	}
	_, err := i.Response().Write(buf.Bytes())
	if err != nil {
		return errors.ErrServer("error sending image " + err.Error())
	}
	return nil
}
