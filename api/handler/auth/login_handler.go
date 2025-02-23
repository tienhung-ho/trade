package userauthhandler

import (
	"client/internal/common/appresponses"
	usermodel "client/internal/model/mysql/user"
	userrepo "client/internal/repository/mysql/user"
	authbusiness "client/internal/service/auth"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LoginWeb2Handler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var data usermodel.UserLoginWeb2

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusInternalServerError, err)
			c.Abort()
			return
		}

		secretKey := os.Getenv("JWT_SECRET_KEY")
		issuer := os.Getenv("JWT_ISSUER")
		audience := os.Getenv("JWT_AUDIENCE")

		store := userrepo.NewMySQLUser(db)
		jwtService := authbusiness.NewJwtService(secretKey, issuer, audience)
		biz := authbusiness.NewAuthLoginBiz(store, jwtService)

		user, token, err := biz.LoginWeb2(c.Request.Context(), &data)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			c.Abort()
			return
		}
		c.SetCookie("access_token", token.AccessToken, 3600, "/", "", true, true)
		// Lưu Refresh Token vào Cookie
		c.SetCookie("refresh_token", token.RefreshToken, 30*24*3600, "/", "", true, true)

		c.JSON(http.StatusOK, appresponses.NewDataResponse(user, "login successfully!"))
	}
}

func LoginWeb3Handler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var data usermodel.UserLoginWeb3

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusInternalServerError, err)
			c.Abort()
			return
		}

		secretKey := os.Getenv("JWT_SECRET_KEY")
		issuer := os.Getenv("JWT_ISSUER")
		audience := os.Getenv("JWT_AUDIENCE")

		store := userrepo.NewMySQLUser(db)
		jwtService := authbusiness.NewJwtService(secretKey, issuer, audience)
		biz := authbusiness.NewAuthLoginBiz(store, jwtService)

		user, token, err := biz.LoginWeb3(c.Request.Context(), &data)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			c.Abort()
			return
		}

		c.SetCookie("access_token", token.AccessToken, 3600, "/", "", true, true)
		// Lưu Refresh Token vào Cookie
		c.SetCookie("refresh_token", token.RefreshToken, 30*24*3600, "/", "", true, true)

		c.JSON(http.StatusOK, appresponses.NewDataResponse(user, "login successfully!"))
	}
}
