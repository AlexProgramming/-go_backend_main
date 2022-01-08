package controllers

import (
	"db"
	"github.com/gofiber/fiber/v2"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"models"
	"net/mail"
	"strconv"
	"time"
	"tokens"
)

// parse the email address for validity
func valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func Register(c *fiber.Ctx) error {
	userMap := make(map[string]string)
	if err := c.BodyParser(&userMap); err != nil {
		return err
	}

	if userMap["password"] != userMap["password_confirmation"] {
		c.Status(fiber.StatusUnauthorized)

		return c.JSON(fiber.Map{
			"response": "password mismatch",
		})
	}

	// check entropy
	entropy := passwordvalidator.GetEntropy(userMap["password"])

	// optionally warn user in response of insecure password
	/*
		if entropy <= 50 {
	} */

	// check email for valid syntax
	if isValid := valid(userMap["email"]); !isValid {
		c.Status(fiber.StatusForbidden)

		return c.JSON(fiber.Map{
			"response": "email format invalid",
		})
	}

	profile := new(models.Profile)

	// hash the password that is stored in db
	if err := profile.SetPassword([]byte(userMap["password"])); err != nil {
		c.Status(fiber.StatusInternalServerError)

		return c.JSON(fiber.Map{
			"response": "Unable to process request, please try again",
		})
	}

	profile.FirstName = userMap["first_name"]
	profile.LastName = userMap["last_name"]
	profile.Email = userMap["email"]
	profile.PasswordEntropy = entropy

	// verify that the email has not already been entered into system before continuing
	if err := db.Instance.Create(&profile); err != nil {
		c.Status(fiber.StatusUnauthorized)

		return c.JSON(fiber.Map{
			"response": "email already registered in system",
		})
	}

	// successful profile add
	c.Status(fiber.StatusCreated)
	return c.JSON(profile)
}

func Login(c *fiber.Ctx) error {
	userMap := make(map[string]string)
	if err := c.BodyParser(&userMap); err != nil {
		return err
	}

	profile := new(models.Profile)

	// check db for existence of user email
	db.Instance.Where("email = ?", userMap["email"]).First(&profile)

	// if Id is 0, account not created
	if profile.Id == 0 {
		c.Status(fiber.StatusNotFound)

		return c.JSON(fiber.Map{
			"response": userMap["email"] + " not found in system",
		})
	}

	// email found in system, authenticate the password with the hash
	if err := profile.VerifyPassword([]byte(userMap["password"])); err != nil {
		c.Status(fiber.StatusUnauthorized)

		return c.JSON(fiber.Map{
			"response": "Incorrect password provided",
		})
	}

	// note jwt are not encrypted and capable of being read, do not use secret unencrypted data
	// login authenticated, provide a json web token to the user
	// for now just sending id unencrypted but add an encryption function

	token, err := tokens.CreateJWT(strconv.FormatUint(profile.Id, 10))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)

		return c.JSON(fiber.Map{
			"response": "Unable to issue web token, please try again",
		})
	}

	cookie := fiber.Cookie{
		Name:     "rest_cookie",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 48), // align with web token creation, param for hours
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	// return success
	c.Status(fiber.StatusOK)

	return c.JSON(fiber.Map{
		"response": "Successfully logged in",
	})
}

func getUser() {

}

func Logout() {

}

func updateUserInfo() {

}

func updatePassword() {

}
