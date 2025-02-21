package authbusiness

import (
	"client/internal/common/apperrors"
	usermodel "client/internal/model/mysql/user"
	walletmodel "client/internal/model/mysql/wallet"
	cosmosservice "client/internal/service/cosmos"
	hashutil "client/internal/util/hash"
	"errors"
	"fmt"
	"log"

	"context"
	"os"
	"strconv"

	"gorm.io/gorm"
)

type AuthInterface interface {
	RegisterUser(db *gorm.DB, data *usermodel.UserRegister, morekeys ...string) (uint64, error)
	GetUser(ctx context.Context, cond map[string]interface{}, morekeys ...string) (*usermodel.User, error)
}

type WalletInterface interface {
	CreateWallet(db *gorm.DB, data *walletmodel.UserWallet, morekeys ...string) (uint64, uint64, error)
}

type AuthBusiness struct {
	store       AuthInterface
	db          *gorm.DB // GORM DB được inject
	cosmosStore cosmosservice.CosmosInterface
	walletStore WalletInterface
}

func NewAuthBiz(
	store AuthInterface,
	db *gorm.DB,
	cosmos cosmosservice.CosmosInterface,
	walletStore WalletInterface,
) *AuthBusiness {
	return &AuthBusiness{
		store:       store,
		db:          db,
		cosmosStore: cosmos,
		walletStore: walletStore,
	}
}

func (biz *AuthBusiness) RegisterUser(ctx context.Context,
	data *usermodel.UserRegister, morekeys ...string) (uint64, error) {

	costEnv := os.Getenv("COST")
	costInt, err := strconv.Atoi(costEnv)
	if err != nil {
		return 0, err
	}

	addrEnv := os.Getenv("ALICE")
	denom := os.Getenv("COIN_NAME")
	adminName := os.Getenv("ADMIN_NAME")
	if adminName == "" {
		return 0, errors.New("ADMIN_NAME is empty")
	}

	// 1) Mở transaction duy nhất
	tx := biz.db.Begin()
	if err := tx.Error; err != nil {
		return 0, apperrors.ErrDB(err)
	}

	// Dùng defer để rollback nếu panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 2) Hash password
	hashUtil := hashutil.NewPasswordManager(costInt)
	hashed, err := hashUtil.HashPassword(data.Password)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	data.Password = hashed

	// 3) Tạo mnemonic
	entropy, err := biz.cosmosStore.GenerateEntropy(256)
	if err != nil {
		tx.Rollback()
		return 0, apperrors.ErrInternal(fmt.Errorf("cannot generate entropy: %w", err))
	}
	mnemonic, err := biz.cosmosStore.GenerateMnemonic(entropy)
	if err != nil {
		tx.Rollback()
		return 0, apperrors.ErrInternal(fmt.Errorf("cannot generate mnemonic: %w", err))
	}
	hashedMnemonic, err := hashUtil.HashPassword(data.Password)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	// 4) Tạo key alias trên Cosmos
	alias := fmt.Sprintf("user-%s", data.Fullname)
	hdPath := "m/44'/118'/0'/0/0"
	err = biz.cosmosStore.CreateNewUser(alias, mnemonic, hdPath)
	if err != nil {
		tx.Rollback()
		return 0, apperrors.ErrInternal(fmt.Errorf("cannot create cosmos key for user: %w", err))
	}

	newAddr, err := biz.cosmosStore.GetAddressFromKeyName(alias)
	if err != nil {
		tx.Rollback()
		// Optionally xóa key rác
		_ = biz.cosmosStore.DeleteKey(alias)
		return 0, apperrors.ErrInternal(err)
	}

	// 5) Mint + Send tokens
	var initCoint int32 = 5
	msg := biz.cosmosStore.MintAndSendTokens(addrEnv, denom, initCoint, newAddr)
	if err := biz.cosmosStore.UpdateTxFactoryAccountSequence(adminName); err != nil {
		tx.Rollback()
		_ = biz.cosmosStore.DeleteKey(alias)
		return 0, fmt.Errorf("cannot update tx factory for admin: %w", err)
	}

	txBuilder, err := biz.cosmosStore.BuildTx(msg)
	if err != nil {
		tx.Rollback()
		_ = biz.cosmosStore.DeleteKey(alias)
		return 0, apperrors.ErrInternal(err)
	}

	if err := biz.cosmosStore.SignTx(ctx, adminName, txBuilder); err != nil {
		tx.Rollback()
		_ = biz.cosmosStore.DeleteKey(alias)
		return 0, apperrors.ErrInternal(err)
	}

	txBytes, err := biz.cosmosStore.EncodeTxBytes(txBuilder)
	if err != nil {
		tx.Rollback()
		_ = biz.cosmosStore.DeleteKey(alias)
		return 0, fmt.Errorf("cannot encode tx: %w", err)
	}

	res, err := biz.cosmosStore.BroadcastTx(txBytes)
	if err != nil {
		tx.Rollback()
		_ = biz.cosmosStore.DeleteKey(alias)
		return 0, fmt.Errorf("cannot broadcast tx: %w", err)
	}
	log.Println(res)

	// 6) Tạo user trong DB
	recordId, err := biz.store.RegisterUser(tx, data, morekeys...) // <--- DÙNG transaction tx
	if err != nil {
		tx.Rollback()
		_ = biz.cosmosStore.DeleteKey(alias)
		return 0, apperrors.ErrCannotCreateEntity(usermodel.EntityName, err)
	}

	// 7) Tạo wallet trong DB
	walletType := os.Getenv("WALLET_TYPE")
	wallet := walletmodel.NewUserWallet(recordId, newAddr, hashedMnemonic, walletType, fmt.Sprint(initCoint), nil)

	walletID, userID, err := biz.walletStore.CreateWallet(tx, wallet) // <--- DÙNG transaction tx
	if err != nil {
		tx.Rollback()
		_ = biz.cosmosStore.DeleteKey(alias)
		return 0, apperrors.ErrCannotCreateEntity(walletmodel.EntityName, err)
	}
	log.Printf("WalletID = %d, UserID = %d\n", walletID, userID)

	// 8) Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		_ = biz.cosmosStore.DeleteKey(alias)
		return 0, apperrors.ErrDB(err)
	}

	// 9) Thành công, return
	return recordId, nil
}
