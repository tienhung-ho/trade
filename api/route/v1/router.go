package routerv1

import (
	userauthhandler "client/api/handler/auth"
	web3handler "client/api/handler/web3"
	jwtmiddleware "client/api/middleware/jwt"
	categoryv1 "client/api/route/v1/category"
	cosmosv1 "client/api/route/v1/cosmos"
	imagev1 "client/api/route/v1/image"
	orderv1 "client/api/route/v1/order"
	productv1 "client/api/route/v1/product"
	walletv1 "client/api/route/v1/wallet"
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
	v1.POST("auth/login-web-2", userauthhandler.LoginWeb2Handler(db))
	v1.POST("auth/login-web-3", userauthhandler.LoginWeb3Handler(db))
	v1.POST("auth/refresh-token", userauthhandler.RefreshToken())
	v1.POST("auth/request-nonce", web3handler.RequestNonce(rdb))
	v1.POST("auth/verify-signature", web3handler.VerifySignature(rdb))
	v1.Use(jwtmiddleware.AuthRequire(db, rdb))
	v1.Use(jwtmiddleware.AuthMiddleware())
	{
		category := v1.Group("/category")
		{
			categoryv1.CategoryRouter(category, db, rdb)
		}
		image := v1.Group("/image")
		{
			imagev1.ImageRouter(image, db)
		}
		product := v1.Group("/product")
		{
			productv1.ProductRouter(product, db, rdb)
		}
		order := v1.Group("/order")
		{
			orderv1.OrderRouter(order, db, rdb, appCtx)
		}
		cosmos := v1.Group("/cosmos")
		{
			cosmosv1.CosmosRouter(cosmos, db, rdb, appCtx)
		}
		wallet := v1.Group("/wallet")
		{
			walletv1.WallterRouter(wallet, db, rdb)
		}

	}
	return r
}
