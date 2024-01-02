package hdwallet

import (
	"crypto/ecdsa"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/txscript"
)

type Key struct {
	Opt      *KeyOption
	Extended *hdkeychain.ExtendedKey

	// btc
	Private *btcec.PrivateKey
	Public  *btcec.PublicKey

	// eth
	PrivateECDSA *ecdsa.PrivateKey
	PublicECDSA  *ecdsa.PublicKey
}

func NewMasterKey(opts ...KeyOptionFunc) (*Key, error) {
	var (
		err error
		o   = newKeyOption(opts...)
	)
	if len(o.Seed) <= 0 {
		o.Seed, err = NewSeed(o.Mnemonic, o.Password, o.Language)
	}
	if err != nil {
		return nil, err
	}
	extended, err := hdkeychain.NewMaster(o.Seed, o.Params)
	if err != nil {
		return nil, err
	}
	key := &Key{
		Opt:      o,
		Extended: extended,
	}
	err = key.init()
	return key, err
}

func (k *Key) init() error {
	var err error
	k.Private, err = k.Extended.ECPrivKey()
	if err != nil {
		return err
	}
	k.Public, err = k.Extended.ECPubKey()
	if err != nil {
		return err
	}
	k.PrivateECDSA = k.Private.ToECDSA()
	k.PublicECDSA = &k.PrivateECDSA.PublicKey
	return nil
}

func (k *Key) GetChildrenKey(opts ...KeyOptionFunc) (*Key, error) {
	var (
		err error
		o   = newKeyOption(opts...)
	)
	extended := k.Extended
	for _, i := range o.GetPath() {
		extended, err = extended.Derive(i)
		if err != nil {
			return nil, err
		}
	}
	key := &Key{
		Opt:      o,
		Extended: extended,
	}
	err = key.init()
	return key, err
}

func (k *Key) AddressBTC() (string, error) {
	address, err := k.Extended.Address(k.Opt.Params)
	if err != nil {
		return "", err
	}
	return address.EncodeAddress(), nil
}

func (k *Key) AddressP2WPKH() (string, error) {
	pubHash, err := k.PublicHash()
	if err != nil {
		return "", err
	}
	addr, err := btcutil.NewAddressWitnessPubKeyHash(pubHash, k.Opt.Params)
	if err != nil {
		return "", err
	}
	return addr.EncodeAddress(), nil
}

func (k *Key) AddressP2WPKHInP2SH() (string, error) {
	pubHash, err := k.PublicHash()
	if err != nil {
		return "", err
	}
	addr, err := btcutil.NewAddressWitnessPubKeyHash(pubHash, k.Opt.Params)
	if err != nil {
		return "", err
	}
	script, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return "", err
	}
	addr1, err := btcutil.NewAddressScriptHash(script, k.Opt.Params)
	if err != nil {
		return "", err
	}
	return addr1.EncodeAddress(), nil
}

func (k *Key) PublicHash() ([]byte, error) {
	address, err := k.Extended.Address(k.Opt.Params)
	if err != nil {
		return nil, err
	}
	return address.ScriptAddress(), nil
}
