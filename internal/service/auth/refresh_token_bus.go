package authbusiness

import (
	"client/internal/common/apperrors"
	jwtmodel "client/internal/model/jwt"
	tokenutil "client/internal/util/token"
	"context"
	"time"
)

type refreshTokenBusiness struct {
	jwtService *JwtService
}

func NewRefreshTokenBusiness(jwtService *JwtService) *refreshTokenBusiness {
	return &refreshTokenBusiness{
		jwtService: jwtService,
	}
}

func (biz *refreshTokenBusiness) RefreshToken(ctx context.Context, claim *jwtmodel.AccountJwtClaims) (*tokenutil.Token, error) {

	timeExpireAccess := time.Duration(1 * time.Hour)
	accessToken, err := biz.jwtService.GenerateToken(claim.ID, claim.WalletID, claim.Email, timeExpireAccess)
	if err != nil {
		return nil, apperrors.ErrInternal(err)
	}

	timeExpireRefresh := time.Duration(30 * 1 * time.Hour)
	refreshToken, err := biz.jwtService.GenerateToken(claim.ID, claim.WalletID, claim.Email, timeExpireRefresh)
	if err != nil {
		return nil, apperrors.ErrInternal(err)
	}

	return tokenutil.NewToken(accessToken, refreshToken), nil
}
