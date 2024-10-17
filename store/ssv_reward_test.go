package store

import (
	"math/big"
	"testing"
)

func TestGetSSVReward(t *testing.T) {
	db := initDB(t)
	ssvReward := db.GetSSVReward("0x064f1e75bbf6da069d7077802fe8271b7acf2c0a")
	t.Log(ssvReward)
}

func TestCreateOrUpdateSSVReward(t *testing.T) {
	db := initDB(t)
	amount := big.NewInt(1435333062)
	err := db.CreateOrUpdateSSVReward("0x4bdf5df5ae729b6797ef1649d3d5b6f63c525ed96118b75d7c38410d316f2bbe", "0x01f959146edf9e2b044e3a1a82dada5b7a3b756f", amount, "0x46362421cdc889695bba28f913d0fc00ef725ff387c02171f17bf3dd6edec07b,0x330e61f3ee2a12753a15d705032cbb622cb3ff1167ca58a24d8a1da970631315,0xcfc8c7fef0552fdfdb4a7a44da5a55c0fba68cc92dfc521c6fab16902676aa4e,0x38e6954566222b37d6e2ebaa3ce3784abd0e28fa58416ade370cbc709a201672,0x541bb22330de0c8c9007b1fa5f523c08d4980c4614ee99645127caf93d2111d6,0x781f406ec98815b13d49360eb3896e7dfa67efd51eaddafafdca21de4b890510,0x5adfb0dcbb21ab056a9955a58d50b3f6ab0d683ebb4be27a9164f6a9c90fd442,0xd04f51ef55d9436012ab43536ec7f5d6bc9c510e0d13af0a1ae0b22172a61d05,0xa7bf9292e79fb26b2669c495e674d4f23593836906e2b87bf6b72e2fe4d15aa4")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetAllAccount(t *testing.T) {
	db := initDB(t)
	accounts, err := db.GetAllSSVRewardAccount()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(accounts)
}
