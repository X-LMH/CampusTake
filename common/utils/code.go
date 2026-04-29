package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// Generate6DigitCode 生成 6 位数字验证码（如 042381）
func Generate6DigitCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000)) // [0, 999999]
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
