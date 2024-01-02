package hdwallet

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"strconv"
	"strings"
)

var (
	DefaultParams       = &BTCParams
	DefaultPassword     = ""
	DefaultLanguage     = English
	DefaultPurpose      = ZeroQuote + 44
	DefaultCoinType     = BTC
	DefaultAccount      = ZeroQuote
	DefaultChange       = Zero
	DefaultAddressIndex = Zero
)

type KeyOption struct {
	Params *chaincfg.Params

	Mnemonic string
	Password string
	Language string
	Seed     []byte

	// bip44
	Purpose      uint32
	CoinType     uint32
	Account      uint32
	Change       uint32
	AddressIndex uint32
}

// GetPath https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki#path-levels
func (o *KeyOption) GetPath() []uint32 {
	return []uint32{
		o.Purpose,
		o.CoinType,
		o.Account,
		o.Change,
		o.AddressIndex,
	}
}

func (o *KeyOption) GetPathString() string {
	return fmt.Sprintf("m/%d'/%d'/%d'/%d/%d",
		o.Purpose-ZeroQuote,
		o.CoinType-ZeroQuote,
		o.Account-ZeroQuote,
		o.Change,
		o.AddressIndex)
}

func newKeyOption(opts ...KeyOptionFunc) *KeyOption {
	opt := &KeyOption{
		Params:       DefaultParams,
		Password:     DefaultPassword,
		Language:     DefaultLanguage,
		Purpose:      DefaultPurpose,
		CoinType:     DefaultCoinType,
		Account:      DefaultAccount,
		Change:       DefaultChange,
		AddressIndex: DefaultAddressIndex,
	}
	for _, f := range opts {
		f(opt)
	}
	return opt
}

type KeyOptionFunc func(option *KeyOption)

func ParamsFunc(p *chaincfg.Params) KeyOptionFunc {
	return func(option *KeyOption) {
		option.Params = p
	}
}

func MnemonicFunc(m string) KeyOptionFunc {
	return func(option *KeyOption) {
		option.Mnemonic = m
	}
}

func PasswordFunc(p string) KeyOptionFunc {
	return func(option *KeyOption) {
		option.Password = p
	}
}

func LanguageFunc(l string) KeyOptionFunc {
	return func(option *KeyOption) {
		option.Language = l
	}
}

func SeedFunc(s []byte) KeyOptionFunc {
	return func(option *KeyOption) {
		option.Seed = s
	}
}

func PurposeFunc(p uint32) KeyOptionFunc {
	return func(option *KeyOption) {
		option.Purpose = p
	}
}

func CoinTypeFunc(c uint32) KeyOptionFunc {
	return func(option *KeyOption) {
		option.CoinType = c
	}
}

func AccountFunc(a uint32) KeyOptionFunc {
	return func(option *KeyOption) {
		option.Account = a
	}
}

func ChangeFunc(c uint32) KeyOptionFunc {
	return func(option *KeyOption) {
		option.Change = c
	}
}

func AddressIndexFunc(a uint32) KeyOptionFunc {
	return func(option *KeyOption) {
		option.AddressIndex = a
	}
}

func PathFunc(path string) KeyOptionFunc {
	return func(o *KeyOption) {
		path = strings.TrimPrefix(path, "m/")
		paths := strings.Split(path, "/")
		if len(paths) != 5 {
			return
		}
		o.Purpose = PathNumber(paths[0])
		o.CoinType = PathNumber(paths[1])
		o.Account = PathNumber(paths[2])
		o.Change = PathNumber(paths[3])
		o.AddressIndex = PathNumber(paths[4])
	}
}

// PathNumber 44' => 0x80000000 + 44
func PathNumber(str string) uint32 {
	num64, _ := strconv.ParseInt(strings.TrimSuffix(str, "'"), 10, 64)
	num := uint32(num64)
	if strings.HasSuffix(str, "'") {
		num += ZeroQuote
	}
	return num
}
