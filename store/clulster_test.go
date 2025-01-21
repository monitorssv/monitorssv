package store

import (
	"testing"
	"time"
)

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

func TestGet30DayLiquidationRankingClusters(t *testing.T) {
	db := initDB(t)
	clusters, totalCount, err := db.Get30DayLiquidationRankingClusters(1, 20, 21663093)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(totalCount)
	t.Log(len(clusters), clusters)
}
