package middleware

import (
	"fmt"
	"go-money-tracker/initializers"
	"go-money-tracker/models"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
)

func ExtractUid(c *gin.Context) int {
	// Find Uid in request headers
	_uid := c.Request.Header["Uid"]
	if _uid == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse Uid from request headers",
		})
		return 0
	}

	// Parse Uid as int
	uid, err := strconv.Atoi(_uid[0])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse Uid in header",
		})
		return 0
	}

	// Lookup if user exists in DB
	var user models.User
	initializers.DB.First(&user, "id = ?", uid)
	if user.ID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Can't find this user",
		})
		return 0
	}

	// Output
	return uid
}

func CheckAuthorization(c *gin.Context) {
	// Get Authorization/tokenString from req headers
	_tokenString := c.Request.Header["Authorization"]
	if _tokenString == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to parse Authorization from header",
		})
		return
	}

	tokenString := _tokenString[0]

	// Decode/validate it
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv(("SECRET_KEY"))), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorizatin token expired, please sign in again",
			})
		}

		// Find the user with token sub
		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "This user is not authorized",
			})
		}

		// Attach to req
		c.Set("user", user)

		// Continue
		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid authorization token",
		})
	}
}

func RequireAuth(c *gin.Context) {
	// Check authorization last due to function behavior
	uid := ExtractUid(c)
	if uid != 0 {
		CheckAuthorization(c)
	}
}
