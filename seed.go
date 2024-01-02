package hdwallet

import (
	"github.com/tyler-smith/go-bip39"
	"github.com/tyler-smith/go-bip39/wordlists"
)

func setLanguage(language string) {
	switch language {
	case English:
		bip39.SetWordList(wordlists.English)
	case ChineseSimplified:
		bip39.SetWordList(wordlists.ChineseSimplified)
	case ChineseTraditional:
		bip39.SetWordList(wordlists.ChineseTraditional)
	}
}

func NewSeed(mnemonic, password, language string) ([]byte, error) {
	setLanguage(language)
	return bip39.NewSeedWithErrorChecking(mnemonic, password)
}
