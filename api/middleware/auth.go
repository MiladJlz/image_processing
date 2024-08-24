package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"image_processing/errors"
	"image_processing/types"
	"net/http"
	"os"
	"time"
)

var secretKey = os.Getenv("SECRET")

type JwtCustomClaims struct {
	ID       int    `json:"ID"`
	Name     string `json:"name"`
	Password string `json:"password"`
	jwt.RegisteredClaims
}

func CreateToken(user *types.User) (*string, *errors.Error) {

	claims := &JwtCustomClaims{user.ID,
		user.Username,
		user.EncryptedPassword,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 5)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secretKey))
	if err != nil {

		return nil, errors.ErrServer(err.Error())
	}

	return &t, nil
}

func JWTAuthentication(c echo.Context, err error) error {

	token, ok := c.Request().Header["Jwt"]
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "token not present in the header")
	}
	claims, err2 := validateToken(token[0])
	if err2 != nil {
		return echo.NewHTTPError(err2.Code, err2.Err)

	}

	c.Set("userID", claims.ID)
	return nil
}
func validateToken(tokenStr string) (*JwtCustomClaims, *errors.Error) {

	token, err := jwt.ParseWithClaims(tokenStr, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, errors.ErrUnAuthorized(err.Error())
	}
	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.ErrUnAuthorized(err.Error())
	}
}
