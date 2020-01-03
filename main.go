package main // only need mysql OR sqlite
// both are included here for reference
import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"gopkg.in/tylerb/graceful.v1"
)

var db *gorm.DB

type Person struct {
	gorm.Model `json:"model"`
	FirstName  string `json:"firstname"`
	LastName   string `json:"lastname"`
	City       string `json:"city"`
}

func handlerFunc(msg string) func(echo.Context) error {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, msg)
	}
}

func getHome(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, world!")
	}
}

func getPeople(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		var people []Person
		db.Find(&people)
		fmt.Println("{}", people)
		return c.JSON(http.StatusOK, people)
	}
}

func getPerson(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		id := c.Param("id")
		var person Person
		db.Where("id = ?", id).First(&person)
		return c.JSON(http.StatusOK, person)
	}
}

func createPerson(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		firstname := c.QueryParam("firstname")
		lastname := c.QueryParam("lastname")
		city := c.QueryParam("city")

		db.Create(&Person{FirstName: firstname, LastName: lastname, City: city})
		return c.String(http.StatusOK, firstname+" person successfully created")
	}
}

func updatePerson(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		id := c.Param("id")
		firstname := c.QueryParam("firstname")
		lastname := c.QueryParam("lastname")
		city := c.QueryParam("city")

		var person Person
		db.Where("id = ?", id).First(&person)
		person.FirstName = firstname
		person.LastName = lastname
		person.City = city
		db.Save(&person)
		return c.String(http.StatusOK, lastname+" person successfully updated")
	}
}

func deletePerson(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		var person Person
		id := c.Param("id")
		d := db.Where("id = ?", id).Delete(&person)
		fmt.Println(d)
		return c.JSON(http.StatusOK, id+" deleted")
	}
}

func handleRequest(db *gorm.DB) {
	e := echo.New()

	e.GET("/", getHome(db))
	e.GET("/people", getPeople(db))
	e.GET("/people/:id", getPerson(db))
	e.POST("/people", createPerson(db))
	e.PUT("/people/:id", updatePerson(db))
	e.DELETE("/people/:id", deletePerson(db))

	port := os.Args[1]
	e.Server.Addr = ":" + port
	graceful.ListenAndServe(e.Server, 5*time.Second)
}

func initialMigration(db *gorm.DB) {
	db.AutoMigrate(&Person{})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbServer := os.Getenv("DB_SERVER")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbDatabase := os.Getenv("DB_DATABASE")
	dbPort := os.Getenv("DB_PORT")

	db, err := gorm.Open("mysql", dbUser+":"+dbPassword+"@tcp("+dbServer+":"+dbPort+")/"+dbDatabase+"?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}

	defer db.Close()
	initialMigration(db)
	handleRequest(db)
}
