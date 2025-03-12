package cosmosservice

import (
	"client/internal/common/apperrors"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
)

type FaucetBusiness struct {
	store CosmosInterface
}

func NewFaucetBiz(store CosmosInterface) *FaucetBusiness {
	return &FaucetBusiness{
		store: store,
	}
}

func (biz *FaucetBusiness) FaucetToken(ctx context.Context, addr string) (string, error) {

	addrEnv := os.Getenv("ALICE")
	denom := os.Getenv("COIN_NAME")
	adminName := os.Getenv("ADMIN_NAME")
	if adminName == "" {
		return "", errors.New("ADMIN_NAME is empty")
	}

	// 5) Mint + Send tokens
	var initCoint int32 = 50000
	msg := biz.store.MintAndSendTokens(addrEnv, denom, initCoint, addr)
	if err := biz.store.UpdateTxFactoryAccountSequence(adminName); err != nil {

		return "", fmt.Errorf("cannot update tx factory for admin: %w", err)
	}

	txBuilder, err := biz.store.BuildTx(msg, denom, "500", 300000)
	if err != nil {

		return "", apperrors.ErrInternal(err)
	}

	if err := biz.store.SignTx(ctx, adminName, txBuilder); err != nil {

		return "", apperrors.ErrInternal(err)
	}

	txBytes, err := biz.store.EncodeTxBytes(txBuilder)
	if err != nil {

		return "", fmt.Errorf("cannot encode tx: %w", err)
	}

	res, err := biz.store.BroadcastTx(txBytes)
	if err != nil {

		return "", fmt.Errorf("cannot broadcast tx: %w", err)
	}
	log.Println(res)

	return res.String(), nil
}
