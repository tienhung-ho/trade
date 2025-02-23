package authbusiness

import (
	jwtmodel "client/internal/model/jwt"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtInterface interface {
	GenerateToken(id uint64, roleId int, email string, expireTime int) (string, error)
	ValidateToken(tokenString string) (*jwtmodel.AccountJwtClaims, error)
}

type JwtService struct {
	secretkey string
	issuer    string
	audience  string
}

func NewJwtService(secretkey, issuer, audience string) *JwtService {
	return &JwtService{secretkey, issuer, audience}
}

func (s *JwtService) GenerateToken(id uint64, walletID uint64, email string, expireTime time.Duration) (string, error) {

	claims := &jwtmodel.AccountJwtClaims{
		ID:       id,
		WalletID: walletID,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireTime)),
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Audience:  []string{s.audience},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign and get the complete encoded token as a string using the secret
	signedToken, err := token.SignedString([]byte(s.secretkey))
	if err != nil {
		return "", errors.New("error signing token")
	}

	return signedToken, nil
}

func (s *JwtService) ValidateToken(tokenString string) (*jwtmodel.AccountJwtClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &jwtmodel.AccountJwtClaims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.secretkey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token is expired")
		}

		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, errors.New("token signature is invalid")
		}

		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, errors.New("token is not valid yet")
		}

		return nil, errors.New("error parsing token")

	}

	claims, ok := token.Claims.(*jwtmodel.AccountJwtClaims)

	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Kiểm tra thêm Issuer và Audience
	if claims.Issuer != s.issuer {
		return nil, errors.New("invalid token issuer")
	}

	// Kiểm tra Audience thủ công
	if len(claims.Audience) == 0 || claims.Audience[0] != s.audience {
		return nil, errors.New("invalid token audience")
	}

	return claims, nil
}
