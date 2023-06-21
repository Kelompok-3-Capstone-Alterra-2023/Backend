package jitsis

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CompleteChallenge(c echo.Context)error{
	data := make(map[string]interface{})
    if err := c.Bind(&data); err != nil {
        return err
    }
	challenge := data["challenge"].(string)
	return c.JSON(http.StatusOK, challenge);
}
