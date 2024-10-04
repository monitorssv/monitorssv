package store

import "testing"

func TestGetLatestBlocks(t *testing.T) {
	db := initDB(t)

	blocks, err := db.GetLatestBlocks()
	if err != nil {
		t.Fatal(err)
	}
	for _, b := range blocks {
		t.Log(b)
	}
}
