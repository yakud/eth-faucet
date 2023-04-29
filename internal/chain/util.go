package chain

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func EtherToWeiWithRound(amount float64) *big.Int {
	ether := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	wei := new(big.Float).Mul(big.NewFloat(amount), new(big.Float).SetInt(ether))
	weiInt, _ := wei.Int(nil) // rounding to nearest integer
	return weiInt
}

func Has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

func IsValidAddress(address string, checksummed bool) bool {
	if !common.IsHexAddress(address) {
		return false
	}
	return !checksummed || common.HexToAddress(address).Hex() == address
}
