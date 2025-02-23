package userauthhandler

import (
	"client/internal/common/apperrors"
	"client/internal/common/appresponses"
	authbusiness "client/internal/service/auth"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func RefreshToken() func(c *gin.Context) {
	return func(c *gin.Context) {
		refreshToken, err := c.Cookie("refresh_token")

		if err != nil {
			c.JSON(http.StatusUnauthorized, apperrors.NewUnauthorized(err, "This action requires login to perform", "ErrRequireLogin", "REFRESH_TOKEN"))
			c.Abort()
			return
		}

		secretKey := os.Getenv("JWT_SECRET_KEY")
		issuer := os.Getenv("JWT_ISSUER")
		audience := os.Getenv("JWT_AUDIENCE")

		jwtService := authbusiness.NewJwtService(secretKey, issuer, audience)

		claims, err := jwtService.ValidateToken(refreshToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			c.Abort()
			return
		}

		biz := authbusiness.NewRefreshTokenBusiness(jwtService)

		tokens, err := biz.RefreshToken(c, claims)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			c.Abort()
			return
		}

		c.SetCookie("access_token", tokens.AccessToken, 3600, "/", "", true, true)

		// Lưu Refresh Token vào Cookie
		c.SetCookie("refresh_token", tokens.RefreshToken, 30*24*3600, "/", "", true, true)

		c.JSON(http.StatusOK, appresponses.NewReponseUserToken(tokens.AccessToken, tokens.RefreshToken, nil))

	}
}
