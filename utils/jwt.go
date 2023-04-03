package utils

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("secret_key23e4")

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
	token, err := auth.SignedString(jwtKey)
	if err != nil {
		log.Fatal("cannot generate key")
		return "", nil
	}
	return token, err
}

func ExtractData(c *gin.Context) (int, string) {
	head := c.GetHeader("Authorization")
	token := strings.Split(head, " ")

	res, _ := jwt.Parse(token[len(token)-1], func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if res.Valid {
		resClaim := res.Claims.(jwt.MapClaims)
		parseID := int(resClaim["ID"].(float64))
		parseRole := resClaim["role"].(string)
		return parseID, parseRole
	}
	return -1, ""
}

func UseJWT(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("ID", int(claims["ID"].(float64)))
			c.Set("role", claims["role"].(string))
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		}
	}
}
