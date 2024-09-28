package ssv

import "testing"

func TestGetOperatorName(t *testing.T) {
	name, err := GetOperatorName("mainnet", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(name)
	name2, err := GetOperatorName("holesky", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(name2)
}

func TestGetOperatorEarnings(t *testing.T) {
	ssv := initSSV(t)
	earnings, err := ssv.GetOperatorEarnings([]uint64{1, 2, 3, 4})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(earnings)
	var operatorEarnings = make([]string, 4)
	t.Log(operatorEarnings[3])

}
