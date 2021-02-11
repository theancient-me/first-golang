package main

import "C"
import (
	"github.com/labstack/echo/v4"
	"net/http"
)


type ResponseJson struct {
	Status int 
	Message string `json:"message" xml:"message"`
}

type PostOriginalUrl struct {
	Url  string `json:"url" form:"url" query:"url"`
}

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		  res := &ResponseJson{
			  Status : 200,
			  Message : "OK",
		  }
		  return c.JSON(http.StatusOK, res)
	})

	e.GET("/:refUrl", func(c echo.Context) error {
		refUrl := c.Param("refUrl")
		return c.String(http.StatusOK, refUrl)
	})
	
	e.POST("/shortUrl", func(c echo.Context) (err error) {
		body := new(PostOriginalUrl)

		if err = c.Bind(body); err != nil {
			res := &ResponseJson{
				Status: 422,
				Message: "url can't be empty",
			}
			return c.JSON(http.StatusUnprocessableEntity,res)
		}
		return c.JSON(http.StatusOK, body)
	})
	
	e.Logger.Fatal(e.Start(":3500"))
}