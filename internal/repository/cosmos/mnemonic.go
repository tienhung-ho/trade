package cosmosrepo

import (
	"fmt"

	"github.com/cosmos/go-bip39"
)

// Đáp ứng interface domain.MnemonicGenerator
func (g *Cosmos) GenerateEntropy(bitSize int) ([]byte, error) {
	entropy, err := bip39.NewEntropy(bitSize)
	if err != nil {
		return nil, fmt.Errorf("couldn't generate entropy: %w", err)
	}
	return entropy, nil
}

func (g *Cosmos) GenerateMnemonic(entropy []byte) (string, error) {
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("couldn't generate mnemonic: %w", err)
	}
	return mnemonic, nil
}
