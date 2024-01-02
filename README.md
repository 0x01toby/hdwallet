# hdwallet

from https://github.com/foxnut/go-hdwallet

## Supports several address derivation for Bitcoin
1. legacy address
2. taproot address
3. nested segwit address
4. native segwit address

## example

``` go
package hdwallet

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func GetMnemonic() string {
	return ""
}

func TestBtc_Address(t *testing.T) {
	masterKey, err := NewMasterKey(MnemonicFunc(GetMnemonic()))
	assert.NoError(t, err)
	nativeSegwitWallet := NewBtcNativeSegwitFromMasterKey(masterKey)
	t.Log("native segwit:		", nativeSegwitWallet.Address())
	nestedSegwitWallet := NewBtcNestedSegwitFromMasterKey(masterKey)
	t.Log("netsted segwit:		", nestedSegwitWallet.Address())
	taprootWallet := NewBtcTaprootFromMasterKey(masterKey)
	t.Log("taproot:				", taprootWallet.Address())
	legacyWallet := NewBtcLegacyFromMasterKey(masterKey)
	t.Log("legacy:				", legacyWallet.Address())
}
```