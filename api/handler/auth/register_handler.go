package userauthhandler

import (
	"client/internal/common/apperrors"
	"client/internal/common/appresponses"
	cosmosmodel "client/internal/model/cosmos"
	usermodel "client/internal/model/mysql/user"
	cosmosrepo "client/internal/repository/cosmos"
	userrepo "client/internal/repository/mysql/user"
	walletrepo "client/internal/repository/mysql/wallets"
	authbusiness "client/internal/service/auth"
	validation "client/internal/util/validate"
	"net/http"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterUserHandler(db *gorm.DB, rdb *redis.Client, appCtx *cosmosmodel.AppContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var data usermodel.UserRegister

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, apperrors.ErrInvalidRequest(err))
			return
		}

		validate := validator.New()
		validate.RegisterValidation("vietnamese_phone", func(fl validator.FieldLevel) bool {
			return validation.IsValidVietnamesePhoneNumber(fl.Field().String())
		})

		// Thực hiện validate
		err := validate.Struct(data)
		if err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				//appErr := common.ErrValidation(validationErrors)
				c.JSON(http.StatusBadRequest, apperrors.ErrValidation(validationErrors))
				return
			}

			// Xử lý lỗi khác nếu có
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		store := userrepo.NewMySQLOrder(db)
		authQueryClient := authtypes.NewQueryClient(appCtx.GRPCConn)
		cosmosStore := cosmosrepo.NewCosmos(appCtx.ClientCtx, appCtx.TxFactory, appCtx.Keyring, authQueryClient, &appCtx.ProtoCodec)
		walletStore := walletrepo.NewMySQLWallet(db)
		biz := authbusiness.NewAuthBiz(store, db, cosmosStore, walletStore)

		id, err := biz.RegisterUser(c.Request.Context(), &data)

		if err != nil {

			c.JSON(http.StatusInternalServerError, err)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, appresponses.NewDataResponse(id, "register new user successfully"))

	}
}
