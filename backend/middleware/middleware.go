package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var validate = validator.New()

func BindAndValidate(c echo.Context, dest interface{}) error {
	if err := c.Bind(dest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid input format")
	}

	if err := validate.Struct(dest); err != nil {
		errs := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			field := strings.ToLower(err.Field())
			switch err.Tag() {
			case "required":
				errs[field] = field + " is required"
			case "email":
				errs[field] = "invalid email format"
			case "min":
				errs[field] = field + " must be at least " + err.Param() + " characters"
			default:
				errs[field] = "invalid value"
			}
		}

		
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"errors": errs,
		})
	}

	return nil
}

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid authorization header")
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "user_id not found in token")
		}

		// Store user_id in the context
		c.Set("user_id", userID)

		return next(c)
	}
}