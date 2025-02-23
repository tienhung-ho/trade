package authbusiness

import (
	"client/internal/common/apperrors"
	usermodel "client/internal/model/mysql/user"
	hashutil "client/internal/util/hash"
	tokenutil "client/internal/util/token"
	"context"
	"os"
	"strconv"
	"time"
)

type AuthLoginInterface interface {
	GetUser(ctx context.Context, cond map[string]interface{},
		morekeys ...string) (*usermodel.User, error)
	GetUserByWalletAddress(ctx context.Context, address string) (*usermodel.User, error)
}

type AuthLoginBusiness struct {
	store      AuthLoginInterface
	jwtService *JwtService
}

func NewAuthLoginBiz(store AuthLoginInterface, jwtService *JwtService) *AuthLoginBusiness {
	return &AuthLoginBusiness{
		store:      store,
		jwtService: jwtService,
	}
}

func (biz *AuthLoginBusiness) LoginWeb2(ctx context.Context, data *usermodel.UserLoginWeb2,
	morekeys ...string) (*usermodel.User, *tokenutil.Token, error) {

	costEnv := os.Getenv("COST")
	costInt, err := strconv.Atoi(costEnv)
	if err != nil {
		return nil, nil, err
	}

	record, err := biz.store.GetUser(ctx, map[string]interface{}{"email": data.Email})

	if err != nil {
		return nil, nil, apperrors.ErrCannotGetEntity(usermodel.EntityName, err)
	}

	hashUtil := hashutil.NewPasswordManager(costInt)
	ok := hashUtil.VerifyPassword(record.Password, data.Password)

	if !ok {
		return nil, nil, apperrors.ErrPasswordInvalid(usermodel.EntityName, err)
	}

	timeExpireAccess := time.Duration(1 * time.Hour)
	accessToken, err := biz.jwtService.GenerateToken(record.ID, record.Wallets.WalletID, record.Email, timeExpireAccess)

	if err != nil {
		return nil, nil, apperrors.ErrInternal(err)
	}

	timeExpireRefresh := time.Duration(30 * 24 * time.Hour)
	refreshToken, err := biz.jwtService.GenerateToken(record.ID, record.Wallets.WalletID, record.Email, timeExpireRefresh)
	if err != nil {
		return nil, nil, apperrors.ErrInternal(err)
	}

	token := tokenutil.NewToken(accessToken, refreshToken)

	return record, token, nil
}

func (biz *AuthLoginBusiness) LoginWeb3(ctx context.Context, data *usermodel.UserLoginWeb3,
	morekeys ...string) (*usermodel.User, *tokenutil.Token, error) {

	record, err := biz.store.GetUserByWalletAddress(ctx, data.WalletAddress)

	if err != nil {
		return nil, nil, apperrors.ErrCannotGetEntity(usermodel.EntityName, err)
	}

	hashUtilMnem := hashutil.NewMnemonicSHA()
	ok := hashUtilMnem.CompareHashSHA256(record.Wallets.EncryptedMnemonic, data.Mnemonic)

	if !ok {
		return nil, nil, apperrors.ErrMnemonicInvalid(usermodel.EntityName, err)
	}

	timeExpireAccess := time.Duration(1 * time.Hour)
	accessToken, err := biz.jwtService.GenerateToken(record.ID, record.Wallets.WalletID, record.Email, timeExpireAccess)

	if err != nil {
		return nil, nil, apperrors.ErrInternal(err)
	}

	timeExpireRefresh := time.Duration(30 * 24 * time.Hour)
	refreshToken, err := biz.jwtService.GenerateToken(record.ID, record.Wallets.WalletID, record.Email, timeExpireRefresh)
	if err != nil {
		return nil, nil, apperrors.ErrInternal(err)
	}

	token := tokenutil.NewToken(accessToken, refreshToken)

	return record, token, nil
}
