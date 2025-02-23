package hashutil

import (
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// PasswordManager chứa các hàm để làm việc với mật khẩu
type passwordManager struct {
	cost int
}

// NewPasswordManager tạo một PasswordManager mới với cost chỉ định
func NewPasswordManager(cost int) *passwordManager {
	return &passwordManager{cost: cost}
}

// HashPassword tạo hash cho mật khẩu
func (pm *passwordManager) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), pm.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (pm *passwordManager) HashName(name string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(name), pm.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword so sánh mật khẩu với hash đã lưu
func (pm *passwordManager) VerifyPassword(hash, passwordLogin string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwordLogin))
	return err == nil
}

func (pm *passwordManager) HashMnemonic(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), pm.cost+2)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (pm *passwordManager) VerifyMnemonic(hash, mnemonicLogin string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(mnemonicLogin))
	return err == nil
}

type mnemonicSHAHash struct{}

func NewMnemonicSHA() *mnemonicSHAHash {
	return &mnemonicSHAHash{}
}

func (mnemonicSHAHash) HashSHA256(plainText string) string {
	h := sha256.New()
	h.Write([]byte(plainText))
	// h.Sum(nil) trả về mảng byte 32 bytes
	hashBytes := h.Sum(nil)

	// Chuyển mảng byte thành chuỗi hex
	return hex.EncodeToString(hashBytes)
}

func (m mnemonicSHAHash) CompareHashSHA256(hashValue, plainValue string) bool {
	// Băm plainValue
	hashedPlain := m.HashSHA256(plainValue)
	// So sánh chuỗi
	return hashedPlain == hashValue
}
