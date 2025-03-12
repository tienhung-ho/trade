package bech32util

import (
	"fmt"
	"strings"
)

// NormalizeBech32Address loại bỏ khoảng trắng và chuyển thành chữ thường.
func NormalizeBech32Address(addr string) string {
	// Loại bỏ tất cả khoảng trắng
	alias := fmt.Sprintf("user-%s", strings.ToLower(addr))
	alias = strings.ReplaceAll(alias, " ", "")
	// Chuyển thành chữ thường
	return strings.ToLower(alias)
}
