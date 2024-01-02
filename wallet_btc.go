package hdwallet

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
)

type Btc struct {
	name        string
	symbol      string
	k           *Key
	AddressType AddressType
}

func NewBTCFromKey(key *Key, addressType AddressType) Wallet {
	return &Btc{
		name:        "Bitcoin",
		symbol:      "BTC",
		k:           key,
		AddressType: addressType,
	}
}

func NewBtcLegacyFromMasterKey(master *Key) Wallet {
	key, _ := master.GetChildrenKey(CoinTypeFunc(BTC), PurposeFunc(ZeroQuote+44))
	return &Btc{
		name:        "Bitcoin",
		symbol:      "BTC",
		k:           key,
		AddressType: BTCLegacy}
}

func NewBtcTaprootFromMasterKey(master *Key) Wallet {
	key, _ := master.GetChildrenKey(CoinTypeFunc(BTC), PurposeFunc(ZeroQuote+86))
	return &Btc{
		name:        "Bitcoin",
		symbol:      "BTC",
		k:           key,
		AddressType: BTCTaproot}
}

func NewBtcNativeSegwitFromMasterKey(master *Key) Wallet {
	key, _ := master.GetChildrenKey(CoinTypeFunc(BTC), PurposeFunc(ZeroQuote+84))
	return &Btc{
		name:        "Bitcoin",
		symbol:      "BTC",
		k:           key,
		AddressType: BTCNativeSegwit}
}

func NewBtcNestedSegwitFromMasterKey(master *Key) Wallet {
	key, _ := master.GetChildrenKey(CoinTypeFunc(BTC), PurposeFunc(ZeroQuote+49))
	return &Btc{
		name:        "Bitcoin",
		symbol:      "BTC",
		k:           key,
		AddressType: BTCNestedSegwit}
}

func (b *Btc) Type() uint32 {
	return b.k.Opt.CoinType
}

func (b *Btc) Name() string {
	return b.name
}

func (b *Btc) Symbol() string {
	return b.symbol
}

func (b *Btc) key() *Key {
	return b.k
}

func (b *Btc) Address() string {
	return b.AddressWithType(b.AddressType)
}

func (b *Btc) AddressWithType(addressType AddressType) string {
	if addressType.EqualTo(BTCLegacy) {
		return b.legacy()
	} else if addressType.EqualTo(BTCTaproot) {
		address, err := b.taprootAddress()
		if err != nil {
			return ""
		}
		return address
	} else if addressType.EqualTo(BTCNativeSegwit) {
		address, err := b.nativeSegwit()
		if err != nil {
			return ""
		}
		return address
	} else if addressType.EqualTo(BTCNestedSegwit) {
		address, err := b.nestedWit()
		if err != nil {
			return ""
		}
		return address
	}
	return ""
}

func (b *Btc) legacy() string {
	address, err := b.k.Extended.Address(b.k.Opt.Params)
	if err != nil {
		return ""
	}
	return address.EncodeAddress()
}

func (b *Btc) taprootAddress() (string, error) {
	serializedPubKey := b.k.Public.SerializeCompressed()
	internalPubkey, err := btcec.ParsePubKey(serializedPubKey)
	if err != nil {
		return "", err
	}
	script := txscript.ComputeTaprootKeyNoScript(internalPubkey)
	key := schnorr.SerializePubKey(script)
	taproot, err := btcutil.NewAddressTaproot(key, b.k.Opt.Params)
	if err != nil {
		return "", err
	}
	return taproot.EncodeAddress(), nil
}

func (b *Btc) nativeSegwit() (string, error) {
	serializedPubKey := b.k.Public.SerializeCompressed()
	hash160 := btcutil.Hash160(serializedPubKey)
	hash, err := btcutil.NewAddressWitnessPubKeyHash(hash160, b.k.Opt.Params)
	if err != nil {
		return "", err
	}
	return hash.EncodeAddress(), nil
}

func (b *Btc) nestedWit() (string, error) {
	serializedPubKey := b.k.Public.SerializeCompressed()
	hash160 := btcutil.Hash160(serializedPubKey)
	hash, err := btcutil.NewAddressWitnessPubKeyHash(hash160, b.k.Opt.Params)
	if err != nil {
		return "", err
	}
	script, err := txscript.PayToAddrScript(hash)
	bytes := btcutil.Hash160(script)
	keyHash, err := btcutil.NewAddressScriptHashFromHash(bytes, b.k.Opt.Params)
	if err != nil {
		return "", err
	}
	return keyHash.EncodeAddress(), nil
}
