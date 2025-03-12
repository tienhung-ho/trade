package jwtmiddleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		jwtSecret := os.Getenv("WEB3_SECRET_KEY")
		// Lấy token từ header Authorization
		tokenString, err := c.Cookie("wall_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token cookie is required"})
			c.Abort()
			return
		}

		// Parse và xác thực token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Đảm bảo thuật toán đúng
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(jwtSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token about web3: " + err.Error()})
			c.Abort()
			return
		}

		// Lấy thông tin từ claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			address, ok := claims["address"].(string)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims in token about web3"})
				c.Abort()
				return
			}

			// Lưu thông tin người dùng vào context để các handler khác có thể sử dụng
			c.Set("walletAddress", address)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims about web3"})
			c.Abort()
			return
		}
	}
}
