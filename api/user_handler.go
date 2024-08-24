package api

import (
	"github.com/labstack/echo/v4"
	"image_processing/api/middleware"
	"image_processing/db"
	"image_processing/types"
	"net/http"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleRegister(c echo.Context) error {
	var params types.CreateUserParams
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	err := params.Validate()
	if err != nil {
		return echo.NewHTTPError(err.Code, err.Err)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return echo.NewHTTPError(err.Code, err.Err)
	}
	user, err = h.userStore.InsertUser(c.Request().Context(), user)
	if err != nil {
		return echo.NewHTTPError(err.Code, err.Err)
	}
	token, err := middleware.CreateToken(user)
	if err != nil {
		return echo.NewHTTPError(err.Code, err.Err)
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{"message": "user created", "token": token})

}
