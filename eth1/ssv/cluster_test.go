package ssv

import (
	"testing"
)

func TestGetContractCreator(t *testing.T) {
	ssv := initSSV(t)
	creations, err := GetContractCreator("0x87393BE8ac323F2E63520A6184e5A8A9CC9fC051", ssv.cfg.EtherScan.Endpoint, ssv.cfg.EtherScan.ApiKey)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(creations)
}
