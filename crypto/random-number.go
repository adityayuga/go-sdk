package crypto

import (
	"crypto/rand"
	"math/big"
)

var numberTable = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
var maxNumberTableIndex = big.NewInt(int64(len(numberTable)))

func GenerateCode(n int) string {
	if n <= 0 {
		return ""
	}

	b := make([]byte, n)
	for i := 0; i < n; i++ {
		numberIndex, err := rand.Int(rand.Reader, maxNumberTableIndex)
		if err != nil {
			return ""
		}
		b[i] = numberTable[numberIndex.Int64()]
	}

	return string(b)
}
