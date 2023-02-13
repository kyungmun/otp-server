package controller

import (
	//"api-fiber-gorm/database"
	"net/http"
	"strconv"

	"github.com/kyungmun/otp-server/middleware"
	"github.com/kyungmun/otp-server/models"
	"github.com/kyungmun/otp-server/service"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	svc *service.UserServices
}

func NewUserController(s *service.UserServices) *UserController {
	return &UserController{svc: s}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func validToken(t *jwt.Token, id string) bool {
	n, err := strconv.Atoi(id)
	if err != nil {
		return false
	}

	claims := t.Claims.(jwt.MapClaims)
	uid := int(claims["user_id"].(float64))

	return uid == n
}

func (c *UserController) validUser(id string, p string) bool {
	user, err := c.svc.GetByID(id)
	if err != nil {
		return false
	}
	if user.Username == "" {
		return false
	}
	if !CheckPasswordHash(p, user.Password) {
		return false
	}
	return true
}

func (c *UserController) GetAll(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 0
	}
	pageSize, err := strconv.Atoi(ctx.Query("page_size"))
	if err != nil {
		pageSize = 0
	}
	users, err := c.svc.GetAll(page, pageSize)
	if err != nil {
		ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"message": err.Error(),
			"result":  false,
		})
		return err
	}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "fiber engine, user all get successfully",
		"count":   len(*users),
		"data":    users,
		"result":  true,
	})
	return nil
}

// GetUser get a user
func (c *UserController) GetUser(ctx *fiber.Ctx) error {
	userid := ctx.Params("user_id")
	if userid == "" {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "id cannot be empty",
			"result":  false,
		})
		return nil
	}
	user, err := c.svc.GetByID(userid)
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
			"id":      userid,
			"result":  false,
		})
		return nil
	}

	if user.Username == "" {
		return ctx.Status(404).JSON(fiber.Map{"status": "error", "message": "No user found with ID", "data": nil})
	}

	return ctx.JSON(fiber.Map{"status": "success", "message": "User found", "data": user})
}

// CreateUser new user
func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
	type NewUser struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	user := new(models.User)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	hash, err := hashPassword(user.Password)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't hash password", "data": err})
	}

	user.Password = hash
	user, err = c.svc.CreateRecord(user)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create user", "data": err})
	}

	newUser := NewUser{
		Email:    user.Email,
		Username: user.Username,
	}

	return ctx.JSON(fiber.Map{"status": "success", "message": "Created user", "data": newUser})
}

/*
// UpdateUser update user
func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	type UpdateUserInput struct {
		Names string `json:"names"`
	}
	var uui UpdateUserInput
	if err := c.BodyParser(&uui); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id) {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})
	}

	db := database.DB
	var user models.User

	db.First(&user, id)
	user.Names = uui.Names
	db.Save(&user)

	return c.JSON(fiber.Map{"status": "success", "message": "User successfully updated", "data": user})
}
*/

// DeleteUser delete user
func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
	type PasswordInput struct {
		Password string `json:"password"`
	}
	var pi PasswordInput
	if err := ctx.BodyParser(&pi); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}
	id := ctx.Params("user_id")
	token := ctx.Locals("user").(*jwt.Token)

	if !validToken(token, id) {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})
	}

	if !c.validUser(id, pi.Password) {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": "Not valid user", "data": nil})
	}

	err := c.svc.DeleteRecord(id)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"status": "error", "message": "Delete user fail", "data": err})
	}

	return ctx.JSON(fiber.Map{"status": "success", "message": "User successfully deleted", "data": nil})
}

func (c *UserController) SetupRoutes(app *fiber.App) {
	api := app.Group("/user")

	api.Post("/", c.CreateUser)
	api.Get("/", c.GetAll)
	api.Get("/:user_id", c.GetUser)
	api.Delete("/:user_id", middleware.Protected(), c.DeleteUser)
}
