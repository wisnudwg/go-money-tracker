package controllers

import (
	"go-money-tracker/initializers"
	"go-money-tracker/models"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	// Get the email/pass from body
	var body struct {
		Email    string
		Name     string
		Password string
		Phone    string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Create new user
	user := models.User{Name: body.Name, Email: body.Email, Password: string(hash), Phone: body.Phone}
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		// if mysqlErr, ok := result.Error.(*mysql.MysqlError); ok {
		//
		// }
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error,
		})
		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"message": "User created",
	})
}

func ReadUser(c *gin.Context) {
	// Extract uid from params
	uid, err := strconv.Atoi(c.Param("uid"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse id from param",
		})
		return
	}

	// Find user based on ID
	var user models.User
	initializers.DB.First(&user, uid)
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Can't find this user",
		})
		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func Login(c *gin.Context) {
	// Get email/pass from body
	var body struct {
		Email    string
		Name     string
		Password string
		Phone    string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Look up requested user
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Compare sent in pass with saved user pass hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Generate jwt-token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign in and get the complete encoded token as a string using the secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	// Send it back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": tokenString,
	})
}

func ValidateToken(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}

func UpdateUser(c *gin.Context) {
	// Get updated user body
	var body struct {
		ID       int
		Email    string
		Name     string
		Password string
		Phone    string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Look up requested user
	var user models.User
	initializers.DB.First(&user, "id = ?", body.ID)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Can't find this user",
		})
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Update user data
	result := initializers.DB.Model(&user).Updates(models.User{Name: body.Name, Email: body.Email, Password: string(hash), Phone: body.Phone})
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update user",
		})
		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"message": "User data updated",
	})
}

func DeleteUser(c *gin.Context) {
	// Extract id from params
	uid, err := strconv.Atoi(c.Param("uid"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse id from param",
		})
		return
	}

	// Find user based on ID
	var user models.User
	initializers.DB.First(&user, uid)
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Can't find this user",
		})
		return
	}

	// Delete user data
	result := initializers.DB.Delete(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error,
		})
		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted",
	})
}
