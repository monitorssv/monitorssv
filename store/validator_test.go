package store

import (
	"testing"
	"time"
)

func TestGetActiveValidatorCount(t *testing.T) {
	db := initDB(t)

	count, err := db.GetActiveValidatorCount()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(count)

	count2, err := db.GetActiveButExitedValidatorCount("1853c5e50b539d5c944e6db8cd54ff839f3bb756ecd39f41a4cc72f7400054dd")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(count2)
}

func TestGetActiveButExitedValidators(t *testing.T) {
	db := initDB(t)

	validators, totalCount, err := db.GetActiveButExitedValidators(1, 10)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(totalCount)
	t.Log(validators)
}

func TestGetValidatorByPublicKey(t *testing.T) {
	db := initDB(t)
	validatorInfo, err := db.GetValidatorByPublicKey("849d44839ba6dcb18d351c3a9e3c66cabef2f7132c7854785b3f14472c1f8cfe4e7337fe6cd7edf254acc7025ffd8f2a")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(validatorInfo)
}

func TestGetValidatorByClusterId(t *testing.T) {
	db := initDB(t)
	validators, totalCount, err := db.GetValidatorByClusterId(1, 10, "1853c5e50b539d5c944e6db8cd54ff839f3bb756ecd39f41a4cc72f7400054dd")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(totalCount)
	t.Log(validators)
}

func TestGetValidators(t *testing.T) {
	db := initDB(t)
	start := time.Now()
	validators, totalCount, err := db.GetValidators(1, 10)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(totalCount)
	t.Log(validators)
	t.Log(time.Since(start))
}

func TestGetValidatorByPubKeyAndBlock(t *testing.T) {
	db := initDB(t)
	validatorInfo, err := db.GetValidatorByPubKeyAndBlock("849d44839ba6dcb18d351c3a9e3c66cabef2f7132c7854785b3f14472c1f8cfe4e7337fe6cd7edf254acc7025ffd8f2a", 17508205)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(validatorInfo)
}

func TestUpdateValidatorIndex(t *testing.T) {
	db := initDB(t)
	err := db.UpdateValidatorIndex("849d44839ba6dcb18d351c3a9e3c66cabef2f7132c7854785b3f14472c1f8cfe4e7337fe6cd7edf254acc7025ffd8f2a", 249366)
	t.Log(err)
}

func TestUpdateValidatorIndex2(t *testing.T) {
	db := initDB(t)
	blocks, err := db.GetAllBlocks()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(blocks[len(blocks)-1])
	t.Log("block:", len(blocks))
	for _, block := range blocks {
		err = db.UpdateValidatorIndex(block.PublicKey, int64(block.Proposer))
		if err != nil {
			t.Fatal(err)
		}
	}
}
