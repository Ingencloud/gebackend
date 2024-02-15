package main

import (
	"ge/controller"
	"ge/database"
	"ge/middleware"
	"ge/models"
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

const (
	codeLength   = 8
	verification = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func generateInviteCode() string {
	rand.Seed(time.Now().UnixNano())
	code := make([]byte, codeLength)
	for i := 0; i < codeLength; i++ {
		code[i] = verification[rand.Intn(len(verification))]
	}
	return string(code)
}
func main() {
	app := fiber.New()
	database.ConnectDb()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:8080/,http://localhost:8080,https://engage.becomingthetackies.site/,https://engage.becomingthetackies.site",
		AllowMethods:     "GET,POST,PUT,DELETE",
	}))

	// Register route
	app.Post("/register", controller.Register)

	// Login route
	app.Post("/login", controller.Login)
	app.Get("/users", controller.AllUsers)

	app.Get("/generate", func(c *fiber.Ctx) error {
		var codes []string
		for i := 0; i < 500; i++ {
			code := generateInviteCode()
			// Insert the generated code into the database
			database.Database.Db.Create(&models.InviteCode{Code: code, IsUsed: false})
			codes = append(codes, code)
		}
		return c.JSON(fiber.Map{"codes": codes})
	})

	// Endpoint to verify and use an invite code
	app.Post("/verify", func(c *fiber.Ctx) error {
		var req struct {
			Code string `json:"code"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
		}

		// Check if the code exists in the database and is not used
		var inviteCode models.InviteCode
		if err := database.Database.Db.Where("code = ? AND is_used = ?", req.Code, false).First(&inviteCode).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or used invite code"})
		}

		// Mark the invite code as used
		database.Database.Db.Model(&inviteCode).Update("is_used", true)

		return c.JSON(fiber.Map{"message": "Invite code verified successfully"})
	})

	// Endpoint to create playlist
	app.Post("/playlist", middleware.IsAuthenticated, func(c *fiber.Ctx) error {
		var playlistUser models.PlaylistDetails
		if err := c.BodyParser(&playlistUser); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
		}

		// Get UserID from the request context
		userID, ok := c.Locals("UserID").(string)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to retrieve user ID"})
		}

		// Parse UserID to integer
		parsedUserID, err := strconv.Atoi(userID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
		}

		playlistUser.UserID = uint(parsedUserID)
		if err := database.Database.Db.Create(&playlistUser).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error adding playlist"})
		}

		return c.JSON(fiber.Map{"message": "Playlist created successfully", "playlist": playlistUser})
	})

	// })

	// Endpoint to list all playlists
	app.Get("/playlists", func(c *fiber.Ctx) error {
		var playlists []models.Playlist
		// Retrieve all playlists from the database
		database.Database.Db.Find(&playlists)

		return c.JSON(fiber.Map{"playlists": playlists})
	})

	// Run the server
	app.Listen(":3000")
}
