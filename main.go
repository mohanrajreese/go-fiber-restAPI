package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go-fiber-api/models"
	"go-fiber-api/storage"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Language  string `json:"language"`
	Publisher string `json:"publisher"`
	Price     int    `json:"price"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBook(c *fiber.Ctx) error {
	book := Book{}
	err := c.BodyParser(&book)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "Request Failed"})
		return err
	}
	err = r.DB.Create(&book).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not create the book"})
		return err
	}
	c.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "Book Created Successfully"})
	return nil
}

func (r *Repository) GetBooks(c *fiber.Ctx) error {
	bookModels := &[]models.Books{}
	err := r.DB.Find(bookModels).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Bad Request"})
		return err
	}
	c.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "Fetch Book success",
			"data": bookModels})
	return nil
}

func (r *Repository) DeleteBook(c *fiber.Ctx) error {
	bookModel := models.Books{}
	id := c.Params("id")
	if id == "" {
		err := c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		if err != nil {
			return err
		}
		return nil
	}
	err := r.DB.Delete(bookModel, id)
	if err.Error != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete the book",
		})
		return err.Error
	}
	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Book Deleted Successfully",
	})
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	v1 := app.Group("/v1/api")
	v1.Post("/create_book", r.CreateBook)
	v1.Delete("/delete/:id", r.DeleteBook)
	v1.Get("get_book/:id", r.GetBooksByID)
	v1.Get("/books", r.GetBooks)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, errr := storage.NewConnection(config)
	if errr != nil {
		log.Fatal("Could not load the database")
	}
	r := Repository{DB: db}

	app := fiber.New()
	r.SetupRoutes(app)
	err = app.Listen(":8080")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Hello Mohan")

}
