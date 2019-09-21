package ldpc

import "math/big"

//Ethereum Block Header for test
type Hash [32]byte
type Address [20]byte
type Bloom [256]byte
type BlockNonce [8]byte

type ethHeader struct {
	ParentHash  Hash
	UncleHash   Hash
	Coinbase    Address
	Root        Hash
	TxHash      Hash
	ReceiptHash Hash
	Bloom       Bloom
	Difficulty  *big.Int
	Number      *big.Int
	GasLimit    uint64
	GasUsed     uint64
	Time        uint64
	Extra       []byte
	MixDigest   Hash
	Nonce       BlockNonce
}
