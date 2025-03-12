package web3handler

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/sha3"
)

var ctx = context.Background()

// Hàm tạo nonce ngẫu nhiên
func generateNonce() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%s-%d", strconv.FormatInt(r.Int63(), 16), time.Now().Unix())
}

// Tạo thông điệp chuẩn để ký
func createSignMessage(nonce string) string {
	return fmt.Sprintf("Sign this message to authenticate with our service\nNonce: %s\nTimestamp: %d",
		nonce, time.Now().Unix())
}

// Hàm trợ giúp so sánh message
func compareMessages(expected string, received []byte) {
	receivedStr := string(received)
	if expected != receivedStr {
		log.Println("Mismatch in message!")
		log.Printf("Expected (string): %q", expected)
		log.Printf("Received (string): %q", receivedStr)
		log.Printf("Expected (hex): %x", []byte(expected))
		log.Printf("Received (hex): %x", received)
	} else {
		log.Println("Message matches expected")
	}
}

// Endpoint cấp phát nonce, lưu vào Redis với TTL 5 phút
func RequestNonce(rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		walletAddress := c.Query("wallet")
		if walletAddress == "" {
			c.JSON(400, gin.H{"error": "Wallet address is required"})
			return
		}

		// Kiểm tra định dạng địa chỉ Cosmos
		if !isValidCosmosAddress(walletAddress) {
			c.JSON(400, gin.H{"error": "Invalid Cosmos wallet address format"})
			return
		}

		nonce := generateNonce()
		message := createSignMessage(nonce)

		log.Print(message)
		// Lưu cả nonce và message vào Redis với TTL 5 phút
		nonceKey := fmt.Sprintf("nonce:%s", walletAddress)
		messageKey := fmt.Sprintf("message:%s", walletAddress)

		pipe := rdb.Pipeline()
		pipe.Set(ctx, nonceKey, nonce, 5*time.Minute)
		pipe.Set(ctx, messageKey, message, 5*time.Minute)
		_, err := pipe.Exec(ctx)

		if err != nil {
			c.JSON(500, gin.H{"error": "Error setting nonce in redis: " + err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"nonce":   nonce,
			"message": message,
		})
	}
}

// Kiểm tra định dạng địa chỉ Cosmos
func isValidCosmosAddress(address string) bool {
	return strings.HasPrefix(address, "cosmos") ||
		strings.HasPrefix(address, "osmo1") ||
		strings.HasPrefix(address, "juno") ||
		strings.HasPrefix(address, "stars")
}

// Chỉnh sửa hàm verifyUserSignatureDirect để nhận expectedMessage làm tham số
func verifyUserSignatureDirect(walletAddress string, signatureB64 string, pubKeyB64 string, signDocJson string, expectedMessage string) (bool, error) {
	// Decode signature from base64
	sigBytes, err := base64.StdEncoding.DecodeString(signatureB64)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature: %v", err)
	}

	// Decode public key from base64
	pubKeyBytes, err := base64.StdEncoding.DecodeString(pubKeyB64)
	if err != nil {
		return false, fmt.Errorf("failed to decode public key: %v", err)
	}

	log.Printf("Pubkey bytes: %v", pubKeyBytes)
	log.Printf("Public key (hex): %x", pubKeyBytes)

	// Parse signDoc from JSON
	var signDoc struct {
		ChainID       string `json:"chainId"`
		AccountNumber string `json:"accountNumber"`
		AuthInfoBytes string `json:"authInfoBytes"`
		BodyBytes     string `json:"bodyBytes"`
	}

	if err := json.Unmarshal([]byte(signDocJson), &signDoc); err != nil {
		return false, fmt.Errorf("failed to unmarshal signDoc JSON: %v", err)
	}

	log.Println("Received signDoc:", signDoc)

	// Decode bodyBytes from base64
	signDocBytes, err := base64.StdEncoding.DecodeString(signDoc.BodyBytes)
	if err != nil {
		return false, fmt.Errorf("failed to decode bodyBytes: %v", err)
	}

	log.Println("Decoded bodyBytes (as string):", string(signDocBytes))
	log.Printf("Decoded bodyBytes (hex): %x", signDocBytes)

	// So sánh message được ký với message mong đợi (expectedMessage)
	compareMessages(expectedMessage, signDocBytes)

	// Tạo đối tượng public key
	pk := &secp256k1.PubKey{Key: pubKeyBytes}

	// METHOD 1: Kiểm tra chữ ký trực tiếp
	if pk.VerifySignature(signDocBytes, sigBytes) {
		log.Println("Verified with direct method")
		accAddress := sdk.AccAddress(pk.Address())
		derivedAddress := accAddress.String()
		log.Println("Derived address:", derivedAddress)
		if derivedAddress != walletAddress {
			return false, fmt.Errorf("address mismatch: derived %s, expected %s", derivedAddress, walletAddress)
		}
		return true, nil
	}

	// METHOD 2: Thử với SHA-256 hash
	hashedMsg := sha256.Sum256(signDocBytes)
	if pk.VerifySignature(hashedMsg[:], sigBytes) {
		log.Println("Verified with SHA-256 hash method")
		accAddress := sdk.AccAddress(pk.Address())
		derivedAddress := accAddress.String()
		log.Println("Derived address:", derivedAddress)
		if derivedAddress != walletAddress {
			return false, fmt.Errorf("address mismatch: derived %s, expected %s", derivedAddress, walletAddress)
		}
		return true, nil
	}

	// METHOD 3: Thử với Keccak-256 hash
	keccak := sha3.NewLegacyKeccak256()
	keccak.Write(signDocBytes)
	keccakHash := keccak.Sum(nil)
	if pk.VerifySignature(keccakHash, sigBytes) {
		log.Println("Verified with Keccak-256 hash method")
		accAddress := sdk.AccAddress(pk.Address())
		derivedAddress := accAddress.String()
		log.Println("Derived address:", derivedAddress)
		if derivedAddress != walletAddress {
			return false, fmt.Errorf("address mismatch: derived %s, expected %s", derivedAddress, walletAddress)
		}
		return true, nil
	}

	return false, fmt.Errorf("signature verification failed with all methods")
}

func VerifySignature(rdb *redis.Client) func(c *gin.Context) {
	return func(c *gin.Context) {

		var req struct {
			WalletAddress string `json:"walletAddress"`
			Signature     string `json:"signature"` // chữ ký dạng base64
			PubKey        string `json:"pubKey"`    // public key dạng base64
			SignDoc       string `json:"signDoc"`   // signDoc đã marshal và mã hóa base64
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Lấy nonce và message từ Redis (để đảm bảo tính một lần)
		messageKey := fmt.Sprintf("message:%s", req.WalletAddress)
		nonceKey := fmt.Sprintf("nonce:%s", req.WalletAddress)

		expectedMessage, err := rdb.Get(ctx, messageKey).Result()
		if err == redis.Nil {
			c.JSON(400, gin.H{"error": "Message not found or expired"})
			return
		} else if err != nil {
			c.JSON(500, gin.H{"error": "Error retrieving message from redis: " + err.Error()})
			return
		}

		// Xác minh chữ ký sử dụng signDoc và expectedMessage từ Redis
		verified, err := verifyUserSignatureDirect(req.WalletAddress, req.Signature, req.PubKey, req.SignDoc, expectedMessage)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if !verified {
			c.JSON(400, gin.H{"error": "Signature verification failed: address mismatch"})
			return
		}

		jwtSecret := os.Getenv("WEB3_SECRET_KEY")
		jwtExpiration := time.Duration(7 * 24 * time.Hour)
		token, err := CreateAuthToken(req.WalletAddress, jwtSecret, jwtExpiration)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error creating authentication token: " + err.Error()})
			return
		}
		// Sau khi xác minh thành công, xóa nonce và message khỏi Redis
		pipe := rdb.Pipeline()
		pipe.Del(ctx, nonceKey)
		pipe.Del(ctx, messageKey)
		_, err = pipe.Exec(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting nonce and message: " + err.Error()})
			return
		}

		c.SetCookie("wall_token", token, 7*24*3600, "/", "", true, true)

		c.JSON(http.StatusOK, gin.H{
			"status":  "verified",
			"address": req.WalletAddress,
			"token":   token,
		})
	}
}

// Tạo JWT token sau khi xác thực thành công (tùy chọn)
func CreateAuthToken(address string, secretKey string, expiration time.Duration) (string, error) {
	claims := struct {
		Address string `json:"address"`
		jwt.RegisteredClaims
	}{
		Address: address,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return signedToken, nil
}
