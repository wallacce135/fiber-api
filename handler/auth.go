package handler

import (
	"api/config"
	"api/database"
	"api/model"
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getUserByEmail(email string) (*model.User, error) {
	db := database.DB
	var user model.User

	if err := db.Where(&model.User{Email: email}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func getUserByUsername(username string) (*model.User, error) {
	db := database.DB
	var user model.User

	if err := db.Where(&model.User{Username: username}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func RegisterHandler(context *fiber.Ctx) error {
	user := new(model.User)

	if err := context.BodyParser(user); err != nil {
		return context.Status(400).JSON(err.Error())
	}

	hash, err := hashPassword(user.Password)
	if err != nil {
		return context.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not hash password", "errors": err.Error()})
	}

	user.Password = hash
	u := database.DB.FirstOrCreate(&user)

	if u.Error != nil {
		return context.Status(400).JSON(fiber.Map{"message": u.Error.Error()})
	}

	if u.RowsAffected == 0 {
		return context.Status(400).JSON("User already exists")
	}

	return context.Status(200).JSON(user)

}

func LoginHandler(context *fiber.Ctx) error {

	type LoginInput struct {
		Indentity string `json:"indentity"`
		Password  string `json:"password"`
	}

	type UserData struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	input := new(LoginInput)

	var ud UserData

	if err := context.BodyParser(input); err != nil {
		return context.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error while logining in", "errors": err.Error()})
	}

	identity := input.Indentity
	pass := input.Password
	userModel, err := new(model.User), *new(error)

	if validateEmail(identity) {
		userModel, err = getUserByEmail(identity)
	} else {
		userModel, err = getUserByUsername(identity)
	}

	fmt.Sprintln(err)

	if err != nil {
		return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Internal Server error(err != nil)", "data": err})
	} else if userModel == nil {
		CheckPasswordHash(pass, "")
		return context.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid identity or password(userModel == nil)", "data": err})
	} else {
		ud = UserData{
			ID:       userModel.ID,
			Username: userModel.Username,
			Email:    userModel.Email,
			Password: userModel.Password,
		}
	}

	if !CheckPasswordHash(pass, ud.Password) {
		fmt.Println("pass ->", pass)
		fmt.Println("ud.Password ->", ud.Password)
		return context.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid identity or password(password hash)", "data": nil})
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["username"] = ud.Username
	claims["user_id"] = ud.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(config.Config("SECRET")))
	if err != nil {
		fmt.Println(err.Error())
		return context.SendStatus(fiber.StatusInternalServerError)
	}

	return context.JSON(fiber.Map{"status": "success", "message": "Success login", "data": t})

}
