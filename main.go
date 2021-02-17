package main

import "C"
import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"math/rand"
	"net/http"
	"time"
)


type ResponseJson struct {
	Status int 
	Message string `json:"message" xml:"message"`
}

type PostOriginalUrl struct {
	Url  string `json:"url" form:"url" query:"url"`
}

type Url struct {
	FullUrl string `json:"full_url"`
	ShortUrl string `json:"short_url"`
}

func findShortUrl(key string) bool{

	db, err := sql.Open("mysql","root:@tcp(127.0.0.1:3306)/testsck")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close();

	result, err := db.Query("select * from url where short_url = ?", key)

	if err != nil{
		panic(err.Error())
	}

	var db_key = true;

	for result.Next(){
		var url Url

		err = result.Scan(&url.FullUrl, &url.ShortUrl)
		if err != nil {
			panic(err.Error())
		}

		db_key = false
	}

	return db_key
}

func generateShortUrl(n int) string {

	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	var b string;
	for i := 0 ; i<n ;i++ {
		b = b + string(letters[rand.Intn(len(letters))])
	}
	return b
}

func main() {
	e := echo.New()

	rand.Seed(time.Now().UnixNano())

	e.GET("/", func(c echo.Context) error {
		  res := &ResponseJson{
			  Status : 200,
			  Message : "OK",
		  }
		  return c.JSON(http.StatusOK, res)
	})

	e.GET("/:refUrl", func(c echo.Context) error {
		refUrl := c.Param("refUrl")
		db, err := sql.Open("mysql","root:@tcp(127.0.0.1:3306)/testsck")
		if err != nil {
			panic(err.Error())
		}

		defer db.Close()

		result, err := db.Query("select * from url where short_url = ?", refUrl)

		if err != nil{
			panic(err.Error())
		}

		var fullUrl string;

		for result.Next(){
			var url Url
			err = result.Scan(&url.FullUrl, &url.ShortUrl)
			if err != nil {
				panic(err.Error())
			}
			fullUrl = url.FullUrl
		}

		res := &ResponseJson{
			Status: 200,
			Message: fullUrl,
		}

		return c.JSON(http.StatusOK, res)
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

		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/testsck")

		if err != nil {
			panic(err.Error())
		}

		defer db.Close()

		fmt.Println(body.Url)
		var status = false
		var short_url string;
		for !status {
			short_url = generateShortUrl(4) //generate short_url
			//short_url = "b5.ntpl.me/"+short_url
			fmt.Println(short_url)
			status = findShortUrl(short_url)
		}

		fmt.Println("can insert")

		insert, err := db.Query(`INSERT INTO url (full_url, short_url) VALUES ( ?, ? )`, body.Url, short_url)

		if err != nil {
			panic(err.Error())
		}

		defer insert.Close()

		res := &ResponseJson{
			Status: 200,
			Message: "b5.tnpl.me/"+short_url,
		}
		return c.JSON(http.StatusOK, res )
	})
	
	e.Logger.Fatal(e.Start(":3500"))
}