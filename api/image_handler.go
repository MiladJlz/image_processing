package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"image"
	"image_processing/db"
	"image_processing/types"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type fileReader struct {
	*bytes.Reader
	buf *bytes.Buffer
}

func NewFileReader(b []byte) *fileReader {
	return &fileReader{bytes.NewReader(b), bytes.NewBuffer(b)}
}

type ImageHandler struct {
	imageStore db.ImageStore
	linkStore  db.LinkStore
}

func NewImageHandler(imageStore db.ImageStore, linkStore db.LinkStore) *ImageHandler {
	return &ImageHandler{
		imageStore: imageStore,
		linkStore:  linkStore,
	}
}

func (h *ImageHandler) HandleUpload(c echo.Context) error {
	userID := c.Get("userID").(int)
	file, _, err := c.Request().FormFile("myFile")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	fr := NewFileReader(data)
	reader := bytes.NewReader(fr.buf.Bytes())
	_, format, err := image.Decode(reader)
	if err != nil {
		echo.NewHTTPError(http.StatusBadRequest, err)
		log.Fatal(err)
	}
	defer file.Close()
	uid, err := uuid.NewUUID()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)

	}
	dst, err := os.Create("./storage/" + uid.String() + "." + format)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	defer dst.Close()
	if _, err := io.Copy(dst, fr.Reader); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	res, err2 := h.imageStore.InsertImage(c.Request().Context(), &types.Image{Name: uid, UserID: userID, Format: format})
	if err2 != nil {
		return echo.NewHTTPError(err2.Code, err2.Err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "file uploaded successfully", "imageID": &res.ID})
}

func generateUniqueLink(s *types.Image) string {
	data := fmt.Sprint(s.Name.String(), time.Now().Unix())
	hash := sha256.Sum256([]byte(data))
	hashStr := hex.EncodeToString(hash[:])
	return hashStr[:7]
}

func (h *ImageHandler) HandleGenerate(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)

	}
	img := NewProcessedImage(c, h, id)
	res, err2 := img.getImageFromDB()
	if err2 != nil {
		return echo.NewHTTPError(err2.Code, err2.Err)
	}
	link := generateUniqueLink(res)
	err2 = h.linkStore.InsertLink(c.Request().Context(), res.ID, link)
	if err2 != nil {
		return echo.NewHTTPError(err2.Code, err2.Err)

	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "link generated successfully", "ling": "localhost:8080/image/" + link})

}
func (h *ImageHandler) HandleImage(c echo.Context) error {
	link := c.Param("link")
	id, err := h.linkStore.GetLinkID(c.Request().Context(), link)
	if err != nil {
		return echo.NewHTTPError(err.Code, err.Err)

	}
	img := NewProcessedImage(c, h, id)

	file, err := img.getImageFile()
	if err != nil {
		return echo.NewHTTPError(err.Code, err.Err)
	}
	if err = img.sendImage(file); err != nil {
		return echo.NewHTTPError(err.Code, err.Err)
	}
	return nil
}

func (h *ImageHandler) HandleTransform(c echo.Context) error {
	transformParams, err := types.NewTransformParams(c)
	if err != nil {
		return echo.NewHTTPError(err.Code, err.Err)
	}
	img := NewProcessedImage(c, h, transformParams.ImageID)
	file, err := img.getImageFile()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	var res *image.NRGBA
	switch transformParams.TransformType {
	case "resize":
		res = imaging.Resize(file, transformParams.TransformBody.Width, transformParams.TransformBody.Height, imaging.NearestNeighbor)
	case "fit":
		res = imaging.Fit(file, transformParams.TransformBody.Width, transformParams.TransformBody.Height, imaging.NearestNeighbor)
	case "fill":
		res = imaging.Fill(file, transformParams.TransformBody.Width, transformParams.TransformBody.Height, imaging.Center, imaging.NearestNeighbor)

	default:
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid transform type")
	}

	if err = img.sendImage(res); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return nil
}
func (h *ImageHandler) HandleFilter(c echo.Context) error {
	filterParams, err := types.NewFilterParams(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Err)
	}
	img := NewProcessedImage(c, h, filterParams.ImageID)
	file, err := img.getImageFile()
	if err != nil {
		return echo.NewHTTPError(err.Code, err.Err)

	}
	var res *image.NRGBA
	switch filterParams.FilterType {
	case "blur":
		res = imaging.Blur(file, filterParams.FilterBody.Value)
	case "sharpening":
		res = imaging.Sharpen(file, filterParams.FilterBody.Value)
	case "gamma":
		res = imaging.AdjustGamma(file, filterParams.FilterBody.Value)
	case "contrast":
		res = imaging.AdjustContrast(file, filterParams.FilterBody.Value)
	case "brightness ":
		res = imaging.AdjustBrightness(file, filterParams.FilterBody.Value)
	case "saturation":
		res = imaging.AdjustSaturation(file, filterParams.FilterBody.Value)
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid filter type")
	}

	if err = img.sendImage(res); err != nil {
		return echo.NewHTTPError(err.Code, err.Err)

	}
	return nil

}
