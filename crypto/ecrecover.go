package crypto

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

var (
	signedMessagePrefix = "\x19Ethereum Signed Message:\n%d"
)

func Ecrecover(msg []byte, sign []byte) (string, error) {
	prefixPack := encodePacked([]byte(fmt.Sprintf(signedMessagePrefix, len(msg))), msg)
	prefixedHash := crypto.Keccak256(prefixPack)

	var v = sign[64]
	V := new(big.Int).SetBytes([]byte{sign[64]})
	if V.Uint64() > 1 {
		if V.Uint64() > 29 {
			v = byte(V.Uint64() - 2 - 8 - 27)
		} else {
			v = byte(V.Uint64() - 27)
		}
	}

	sign[64] = v

	pubKey, err := crypto.SigToPub(prefixedHash, sign)
	if err != nil {
		return "", err
	}

	addr := crypto.PubkeyToAddress(*pubKey)
	return addr.String(), nil
}

func encodePacked(input ...[]byte) []byte {
	return bytes.Join(input, nil)
}
