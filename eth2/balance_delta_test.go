package eth2

import "testing"

func TestValidatorMonitor(t *testing.T) {
	bm := initBeaconMonitor(t)
	err := bm.validatorMonitor(313123)
	if err != nil {
		t.Fatal(err)
	}
}
