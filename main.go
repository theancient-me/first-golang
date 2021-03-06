package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"math/rand"
	"net/http"
	"time"
)

const (
	username = ""
	password = ""
	hostname = ""
	dbname   = ""
)

func OpenConnection() *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(0)
	db.SetMaxIdleConns(500)

	return db
}

func dsn(dbName string) string {
    return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

type ResponseJson struct {
	Status int 
	Message string `json:"link" xml:"link"`
}

type ResponseStatJson struct {
	Status int
	Message int `json:"visit"`
}

type PostOriginalUrl struct {
	Url  string `json:"url" form:"url" query:"url"`
}

type Url struct {
	FullUrl string `json:"full_url"`
	ShortUrl string `json:"short_url"`
	Visit int `json:"visit"`
}

func findShortUrl(key string) bool{

	//db, err := sql.Open("mysql",dsn(dbname))
	//
	//if err != nil {
	//	panic(err.Error())
	//}

	db := OpenConnection()

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
			Status:  200,
			Message: "OK",
		}
		return c.JSON(http.StatusOK, res)
	})


	e.GET("/l/:refUrl", func(c echo.Context) error {
		var url Url
		refUrl := c.Param("refUrl")

		db := OpenConnection()
		defer db.Close()

		row := db.QueryRow(`SELECT full_url, short_url, visits FROM url WHERE short_url = ?`, refUrl)

		err := row.Scan(&url.FullUrl, &url.ShortUrl, &url.Visit)
		updated, _ := db.Query(`update url set visits = visits+1 where short_url = ?`, refUrl)
		updated.Close()
		if err != nil {
			panic(err.Error())
		}
		c.Response().Header().Set("Location", url.FullUrl)
		res := &ResponseJson{
			Status:  302,
			Message: url.FullUrl,
		}
		return c.JSON(http.StatusFound, res)

		return c.String(http.StatusOK,"http.StatusOK");

	})

	e.GET("/l/:refUrl/stat", func(c echo.Context) error {
		var url Url
		refUrl := c.Param("refUrl")

		//db, err := sql.Open("mysql",dsn(dbname))
		//
		//if err != nil {
		//	panic(err.Error())
		//}
		db := OpenConnection()

		result := db.QueryRow("select visits from url where short_url = ?", refUrl);
		err := result.Scan(&url.Visit)
		if err != nil {
			panic(err.Error())
		}

		res := &ResponseStatJson{
			Status: 200,
			Message: url.Visit,
		}

		return c.JSON(http.StatusOK, res);
	})

	e.POST("/link", func(c echo.Context) (err error) {
		body := new(PostOriginalUrl)

		if err = c.Bind(body); err != nil {
			res := &ResponseJson{
				Status:  422,
				Message: "url can't be empty",
			}
			return c.JSON(http.StatusUnprocessableEntity, res)
		}

		//db, err := sql.Open("mysql", dsn(dbname))
		//
		//if err != nil {
		//	panic(err.Error())
		//}

		db := OpenConnection()

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
			Status:  200,
			Message: "http://sh.b5.tnpl.me/l/" + short_url,
		}
		return c.JSON(http.StatusOK, res)
	})

	e.Logger.Fatal(e.Start(":3500"))
}