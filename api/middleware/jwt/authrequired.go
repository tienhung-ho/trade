package jwtmiddleware

import (
	"client/internal/common/apperrors"
	authbusiness "client/internal/service/auth"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func AuthRequire(db *gorm.DB, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, errAccess := c.Cookie("access_token")
		refreshToken, errRefresh := c.Cookie("refresh_token")

		if errAccess != nil && errRefresh == nil && refreshToken != "" {
			c.JSON(http.StatusUnauthorized, apperrors.TokenExpired("Access Token", errAccess))
			c.Abort()
			return
		}

		if errAccess != nil {
			c.JSON(http.StatusForbidden, apperrors.NewUnauthorized(errAccess, "This action requires login to perform", "ErrRequireLogin", "ACCESS_TOKEN"))
			c.Abort()
			return
		}

		secretKey := os.Getenv("JWT_SECRET_KEY")
		issuer := os.Getenv("JWT_ISSUER")
		audience := os.Getenv("JWT_AUDIENCE")

		jwtService := authbusiness.NewJwtService(secretKey, issuer, audience)

		claims, err := jwtService.ValidateToken(accessToken)

		if err != nil {
			switch err.Error() {
			case "token is expired":
				c.JSON(http.StatusUnauthorized, apperrors.TokenExpired("Access Token", err))
				return
			case "token signature is invalid":
				c.JSON(http.StatusUnauthorized, apperrors.NewUnauthorized(err, "Token signature is invalid", "ErrInvalidTokenSignature", "ACCESS_TOKEN"))
			case "invalid token issuer":
				c.JSON(http.StatusUnauthorized, apperrors.NewUnauthorized(err, "Token issuer is invalid", "ErrInvalidTokenIssuer", "ACCESS_TOKEN"))
			case "invalid token audience":
				c.JSON(http.StatusUnauthorized, apperrors.NewUnauthorized(err, "Token audience is invalid", "ErrInvalidTokenAudience", "ACCESS_TOKEN"))
			default:
				c.JSON(http.StatusUnauthorized, apperrors.ErrInternal(err))
			}
			c.Abort()
			return
		}

		id := claims.ID
		walletID := claims.WalletID
		email := claims.Email

		c.Set("id", id)
		c.Set("wallet", walletID)
		c.Set("email", email)

		c.Next()
	}
}
