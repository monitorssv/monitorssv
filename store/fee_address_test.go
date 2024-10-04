package store

import (
	"testing"
)

func TestUpdateClusterFeeAddress(t *testing.T) {
	db := initDB(t)

	err := db.CreateOrUpdateClusterFeeAddress(&FeeAddressInfo{
		Owner:      "0xabf1ADf95AA7eD243672CeFC194E8411779300df",
		FeeAddress: "0x69eEd4905BC2A4a6381F2791c7644D1018AaC843",
	})
	if err != nil {
		t.Fatal(err)
	}

	fee, err := db.GetClusterFeeAddress("0xabf1ADf95AA7eD243672CeFC194E8411779300df")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fee)

	fee, err = db.GetClusterFeeAddress("0x69eEd4905BC2A4a6381F2791c7644D1018AaC843")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fee.FeeAddress == "")
}
