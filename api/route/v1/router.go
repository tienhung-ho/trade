package routerv1

import (
	userauthhandler "client/api/handler/auth"
	cosmosmodel "client/internal/model/cosmos"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, rdb *redis.Client, appCtx *cosmosmodel.AppContext) *gin.Engine {

	r := gin.Default()

	// Cấu hình CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:9000"}, // Cho phép origin từ frontend
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	//	r.POST("/login", authhandler.LoginHandler(db))
	//	r.POST("/refresh-token", authhandler.RefreshToken())

	v1 := r.Group("/api/v1")
	v1.POST("auth/register", userauthhandler.RegisterUserHandler(db, rdb, appCtx))
	//	v1.Use(authmiddleware.AuthRequire(db, rdb))
	{
	}
	return r
}
