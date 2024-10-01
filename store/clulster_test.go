package store

import (
	"testing"
	"time"
)

func TestUpdateClusterFeeAddress(t *testing.T) {
	db := initDB(t)

	err := db.UpdateClusterFeeAddress("0xabf1ADf95AA7eD243672CeFC194E8411779300df", "0x69eEd4905BC2A4a6381F2791c7644D1018AaC843")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateClustersOwner(t *testing.T) {
	db := initDB(t)
	owners, err := db.GetNoUpdatedClustersOwner()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(owners), owners)
}

func TestGetActiveClusters(t *testing.T) {
	db := initDB(t)
	start := time.Now()
	activeClusters, err := db.GetActiveClusters()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(activeClusters))
	cost := time.Since(start)
	t.Log(cost)

	count, err := db.GetActiveClusterCount()
	t.Log(count, err)
}
