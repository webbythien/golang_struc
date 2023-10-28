package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
)

type Comment struct {
	Id     uint   `json:"id"`
	PostId uint   `json:"post_id"`
	Text   string `json:"text"`
}

func main() {
	db, err := gorm.Open(mysql.Open("root:rootroot@tcp(127.0.0.1:3306)/comments_ms"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(Comment{})

	app := fiber.New()
	app.Use(cors.New())

	app.Get("/api/posts/:id/comments",
		func(c *fiber.Ctx) error {
			var comments []Comment

			db.Find(&comments, "post_id = ?", c.Params("id"))
			return c.JSON(comments)
		})

	app.Post("/api/comments", func(c *fiber.Ctx) error {
		var comment Comment
		if err := c.BodyParser(&comment); err != nil {
			return err
		}
		db.Create(&comment)

		if rand.Intn(9) <= 8 {

			body, _ := json.Marshal(map[string]string{
				"text": comment.Text,
			})
			url := fmt.Sprintf("http://localhost:8000/api/posts/%d/comments", comment.PostId)
			http.Post(url, "application/json", bytes.NewBuffer(body))

		}
		return c.JSON(comment)

	})

	app.Listen(":8001")
}
