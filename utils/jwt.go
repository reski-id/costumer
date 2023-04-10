package utils

import (
	"costumer/models"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var JwtKey = []byte("secret_key23e4")

type Claims struct {
	UserID int  `json:"user_id"`
	Role   bool `json:"is_admin"`
	jwt.StandardClaims
}

func GenerateToken(userID int, Role string) (string, error) {
	// function body
	info := jwt.MapClaims{}
	info["ID"] = userID
	info["role"] = Role
	auth := jwt.NewWithClaims(jwt.SigningMethodHS256, info)
	token, err := auth.SignedString(JwtKey)
	if err != nil {
		log.Fatal("cannot generate key")
		return "", nil
	}
	return token, err
}

func ExtractData(c *gin.Context) (int, string, error) {
	head := c.GetHeader("Authorization")
	if head == "" {
		return -1, "", fmt.Errorf("Authorization header not provided")
	}

	token := strings.Split(head, " ")

	res, err := jwt.Parse(token[len(token)-1], func(t *jwt.Token) (interface{}, error) {
		return []byte(JwtKey), nil
	})
	if err != nil {
		return -1, "", err
	}

	if res.Valid {
		resClaim := res.Claims.(jwt.MapClaims)
		parseID := int(resClaim["ID"].(float64))
		parseRole := resClaim["role"].(string)
		return parseID, parseRole, nil
	}

	return -1, "", fmt.Errorf("Invalid token")
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid Authorization header"})
			c.Abort()
			return
		}

		tokenString := authHeader[7:]
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(JwtKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: err.Error()})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			c.Set("user_id", claims.UserID)
			if !claims.Role {
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "You do not have permission to access this resource"})
				c.Abort()
				return
			}
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid token"})
			c.Abort()
			return
		}
	}
}
