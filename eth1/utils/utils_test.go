package utils

import (
	"math/big"
	"testing"
)

func TestToSSV(t *testing.T) {
	f := ToSSV(big.NewInt(0), "%.18f")
	t.Log(f)

	f2 := ToSSV(big.NewInt(10000000000), "%.18f")
	t.Log(f2)

	f3 := ToSSV(big.NewInt(10000000000000000), "%.18f")
	t.Log(f3)
}
