package store

import (
	"testing"
)

func TestUpdateScanEth1Block(t *testing.T) {
	db := initDB(t)
	err := db.UpdateScanEth1Block(20425900)
	if err != nil {
		t.Fatal(err)
	}
}
