package eth2

import "testing"

func TestValidatorMonitor(t *testing.T) {
	bm := initBeaconMonitor(t)
	err := bm.validatorMonitor(313123)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheckBalanceChange(t *testing.T) {
	bm := initBeaconMonitor(t)
	bm.updateBalanceHistory(1, 1, 100)
	t.Log(bm.validatorBalanceHistory)
	t.Log(bm.checkBalanceChange(1, true))
	bm.updateBalanceHistory(1, 2, 101)
	t.Log(bm.validatorBalanceHistory)
	bm.updateBalanceHistory(1, 3, 102)
	t.Log(bm.validatorBalanceHistory)
}
