package types

import (
	"github.com/labstack/echo/v4"
	"image_processing/errors"
	"strconv"
	"strings"
)

type FilterRequestBody struct {
	Value float64 `json:"value"`
}

type TransformRequestBody struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type TransformParams struct {
	ImageID       int
	TransformType string
	TransformBody TransformRequestBody
}

func NewTransformParams(c echo.Context) (*TransformParams, *errors.Error) {
	transformType := c.Param("transform")
	transformType = strings.ToLower(transformType)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {

		return nil, errors.ErrBadRequest(err.Error())
	}

	var reqBody TransformRequestBody
	params, err := c.FormParams()
	if err != nil {
		return nil, errors.ErrBadRequest(err.Error())
	}
	width := params.Get("width")
	height := params.Get("height")
	reqBody.Width, err = strconv.Atoi(width)
	if err != nil {
		return nil, errors.ErrBadRequest(err.Error())
	}
	reqBody.Height, err = strconv.Atoi(height)
	if err != nil {
		return nil, errors.ErrBadRequest(err.Error())
	}
	return &TransformParams{ImageID: id, TransformType: transformType, TransformBody: reqBody}, nil

}

type FilterParams struct {
	ImageID    int
	FilterType string
	FilterBody FilterRequestBody
}

func NewFilterParams(c echo.Context) (*FilterParams, *errors.Error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return nil, errors.ErrBadRequest(err.Error())
	}
	filterType := c.Param("filter")
	filterType = strings.ToLower(filterType)
	var reqBody FilterRequestBody
	param := c.FormValue("value")
	reqBody.Value, err = strconv.ParseFloat(strings.TrimSpace(param), 64)
	if err != nil {
		return nil, errors.ErrBadRequest(err.Error())
	}
	return &FilterParams{ImageID: id, FilterType: filterType, FilterBody: reqBody}, nil
}
