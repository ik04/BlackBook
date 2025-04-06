package middleware

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)
var validate = validator.New()

func BindAndValidate(c echo.Context, dest interface{}) error {
	if err := c.Bind(dest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid input format")
	}
	if err := validate.Struct(dest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
