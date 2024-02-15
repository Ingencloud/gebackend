package controller

import (
	"fmt"
	"ge/database"
	"ge/models"
	"ge/utils"
	"math"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	User := models.User{
		FullName: data["name"],
		Number:   data["number"],
	}
	database.Database.Db.Create(&User)

	return c.JSON(User)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	if result := database.Database.Db.Where("number = ?", data["number"]).First(&user); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"message": "User not found",
			})
		}
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Could not login",
		})
	}

	// Generate JWT token
	token, err := utils.GenerateJwt(strconv.Itoa(int(user.ID)))
	if err != nil {
		fmt.Println("Token Generation Error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not login",
		})
	}
	fmt.Println("Generated Token:", token)

	// Set JWT token as a cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	// fmt.Println("Set Cookie:", cookie)
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Success",
		"token":   token,
	})
}

func AllUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit := 10
	offset := (page - 1) * limit
	var total int64
	var users []models.User

	database.Database.Db.Offset(offset).Limit(limit).Find(&users)
	database.Database.Db.Model(&models.User{}).Count(&total)
	result := database.Database.Db.Preload("PlaylistDetails").Preload("PlaylistDetails.Playlist").Find(&users)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
	}

	return c.JSON(fiber.Map{
		"data": users,
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"last_page": math.Ceil(float64(int(total) / limit)),
		},
	})
}
