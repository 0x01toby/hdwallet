package hdwallet

var (
	BTCLegacy       = AddressType(0)
	BTCTaproot      = AddressType(1)
	BTCNativeSegwit = AddressType(2)
	BTCNestedSegwit = AddressType(3)
)

type AddressType uint32

func (a AddressType) EqualTo(t AddressType) bool {
	return uint32(a) == uint32(t)
}

type Wallet interface {
	Type() uint32
	Name() string
	Symbol() string
	Address() string
	key() *Key
}
